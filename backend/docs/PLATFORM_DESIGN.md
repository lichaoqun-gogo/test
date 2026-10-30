# 平台设计规则备查（终版）

> 用途：沉淀平台关键设计决策，供后续开发/审查/测试对照。
> 维护：凡涉及平台设计的内容讨论定稿后，同步更新本文件。

---

## 1. 生命周期治理流程 · 审批流设计

### 1.1 万能审批流模板（universal，四节点）

| 节点 | 名称 | 类型 | 审批人 | 说明 |
|---|---|---|---|---|
| node_0 | 发起 | 固定 | 申请人 | 提交时自动通过（发起动作） |
| node_1 | 本部门/运维部门专员复审 | 条件 | **引入=申请人部门**；**升级/退出=项目运维部门** `project.maint_dept`（无则退回申请人部门）`dept_specialist` | 防重复：发起人若本身就是该部门在岗专员则跳过 |
| node_2 | 使用部门确认 | 条件 | 各使用部门 `usage_dept_specialist` | **仅升级/退出申请激活**；每个使用部门一个独立节点（随机顺序），自动跳过与运维部门相同的部门 |
| node_3 | 管理部门评审 | 固定 | 开源治理办公室 `mgmt_dept_specialist` | 终审 |

**部门口径（重要）**：使用部门来自 `project.use_dept`（可逗号分隔多个）；运维部门来自 `project.maint_dept`。node_1 复审部门与 node_2 的"与运维部门相同则跳过"参照，均以 `project.maint_dept` 为准（不再用申请人部门）。

### 1.2 node_2（使用部门确认）自动跳过规则

1. 业务类型非 upgrade/withdraw（如引入申请）→ 跳过
2. 升级/退出但未指定使用部门（project_dept 为空）→ 跳过
3. 使用部门与运维部门（applicant_dept）相同 → 过滤该部门（防重复）
4. 所有使用部门都与运维部门相同 → 整个 node_2 跳过
5. 使用部门未配置 usage 专员 → 跳过该部门节点
6. 所有使用部门均未配置专员 → 整个 node_2 跳过

**为什么引入申请不需要使用部门确认**：引入是"新增"，使用方就是申请人所在部门（本部门专员已在 node_1 复审），无需再确认；升级/退出是"存量变更"，项目可能被多个使用部门使用，需各使用部门专员确认（是否同意升级/退出、确认影响范围）。

### 1.3 任务状态机（五态 + 驳回可重提）

```
draft(草稿) --提交--> approving(审批中) --全部通过--> approved(已通过/终态)
approving --撤回(仅首节点未处理)--> withdrawn(已撤回/可重提) --提交--> approving
approving --驳回--> rejected(已驳回/可重提) --编辑+提交--> approving
```

- **撤回**：仅审批中、且第一个实际审批节点（node_order>0）状态为 pending 时可撤回（重提后节点重置 pending，可再次撤回）
- **驳回**：可重提态（非终态）。发起人可编辑修改后重新提交，重新走审批流（节点重置 pending）；驳回意见留存于审批历史
- **重提节点重置规则**：重提只重置**非跳过节点**为 pending（重新审批）；**跳过节点保持跳过**——重提不改变业务类型/使用部门，跳过条件不变（曾因把跳过节点重置成 pending 空审批人导致流程卡死，已修复）
- **不可删除**：发起人不能删除自己发起的任务，平台留存记录
- **审批历史完整可追溯**：撤回/驳回重提是同一任务的多次流转，审批历史按时间正序**完整返回**（含首次提交、驳回意见、重提等），不做截断
- **更新限制**：仅 draft / withdrawn / rejected 状态任务可编辑更新（approving 数据锁定）

### 1.4 专员分类与审批权限

- **部门专员**（dept_specialist / usage_dept_specialist）：一个部门可多人，任一在岗专员可处理本部门节点
- **管理专员**（mgmt_dept_specialist，开源治理办公室）：终审 + 全局治理视图
- 节点审批人由 `dept_specialist_mapping` 表解析（后期由权限系统同步角色到表）

---

## 2. 缺陷台账权限分层

| 角色 | 个人工作台·缺陷管理（scope=dept） | 全局治理·缺陷台账（scope≠dept） |
|---|---|---|
| 管理专员 | 本部门（等同任意部门专员） | **全公司全部** |
| 部门专员 | 本部门 + 本部门有处置单的缺陷 | （按部门） |
| 普通用户 | 本部门 + 本部门有处置单的缺陷 | 无入口 |

设计依据：缺陷是台账信息，外部研判可能不准，**全局治理看全部兜底**；个人工作台按部门分工。管理专员可切换两种视图。

---

## 3. 缺陷管理设计

### 3.1 实体/流程分离（两实体表 + 流程明细 + 处置 + 历史）

- 外部缺陷实体：`defect`（defect_id = 软件#编号，如 redis#13401，幂等导入，专表专用）
- 内部缺陷实体：`defect_ledger`（defect_no = DF-xxx，风险上报字段）
- 流程明细：`defect_flow`（无 status 列）
- 处置工单：`defect_disposal`（按缺陷编号关联）
- 状态历史：`defect_status_history`
- **跨表关联一律按缺陷编号（defect_id / defect_no），不依赖实体表主键**（实体可反复导入 upsert）
- 内部处置状态 `internal_status` 放实体表（pending_analysis → pending_review / ignored / assigned / closed）

### 3.2 外部缺陷导入（幂等 + 就高不就低）

- 按 defect_id 主键 upsert，同一编号只有一条记录，无论谁导入
- 更新列不含 internal_status（状态不被普通导入覆盖）
- 研判预判升级规则（**就高不就低**）：
  ```
  等级：待研判(1) < 忽略(2) < 待上会(3) < 已分配(4) < 已闭环(5)
  ```
  新预判等级 > 当前等级才升级（忽略可被待上会覆盖）；绝不降级/回退
- 首次导入记录部门（maint_dept），重复导入不覆盖

### 3.3 研判来源与「待研判」兜底（定案）

**研判有两个来源，职责分工明确：**
- **外部采集端预判（主流程）**：采集端已研判，导入时带 `disposal_status`（待上会 pending_review / 忽略 ignored）。外部缺陷导入即研判，导入者即为研判归属，平台**不重复研判**。
- **平台内研判（兜底）**：仅当采集端**未给研判结论**（disposal_status 缺省）时，缺陷落「待研判 pending_analysis」，由管理专员分配软件负责人研判（`AssignDefectAnalyst` / `AnalyzeDefect`：pending_analysis → pending_review / assigned / ignored）。

**定案依据**：采集端偶尔可能漏传研判结论，故**保留「待研判 + 分配研判」作为兜底通道**，不删除。外部已研判缺陷不受影响，不产生重复研判。

**dispatch_type 语义（研判人 vs 跟踪人）**：`defect_flow.dispatch_type` 区分分配类型——`analyze`=真正分配了研判人（前端据此显示/隐藏「分配研判人」按钮）；`import`=导入跟踪人（导入者归属，**不算已分配研判人**）。待研判缺陷（import 或无）可分配研判人；分配后置为 analyze。**跟踪人/研判人/操作人统一存 username**（导入者取 X-User-ID），显示时由 `fillDispatchTargetNames` 映射昵称；「待我研判」只列 `dispatch_type='analyze'` 且 `dispatch_target=当前用户` 的缺陷，import 跟踪人不进待办。

**就高不就低**：外部多次导入不同预判结论时，按等级升级（忽略可被待上会覆盖），绝不降级；平台内研判同理。

### 3.4 处置流程

- 导入 → 研判（分配/自动分配研判人）→ 审核上会（需要修复/忽略）→ 处置分配 → 全部完成 → 缺陷闭环
- **研判结论口径**：缺陷可「需要上会评审」或「忽略」（与外部研判口径一致：待上会/忽略）；「直接分配（need_fix）」仅漏洞支持。缺陷需要修复分配仍须先上会评审（评审结论：需要修复/无需处置）。

**处置分配两通道（定案，与研判分配入口对齐）**：
- 入口：评审弹框内「结论=需要修复」时出现分配区块；漏洞在研判弹框「结论=需要修复」时同样出现。两通道：**自动分配**（默认）/ **手动指定**。
- **手动指定**：与内部缺陷审核分配一致——输入负责人账号，**可填多个，用逗号/顿号/分号/空格分隔**（`DisposalAssignPicker` 公共组件统一解析）；每人各建 1 张处置单（处置部门=各自所属部门），流程 `maint_dept` 同步为第一个负责人部门。不校验部门归属，任意部门人员均可指定。
- **自动分配**：按涉及项目裂变 + 自动匹配处理人（见下方裂变拆单规则）。
- **内部缺陷（风险上报）审核**：固定手动指定（无自动入口，showAuto=false）。
- **漏洞（flow_type=vuln）研判 need_fix**：同样两通道；自动通道与缺陷一致，按组件类型裂变（基础软件如 redis / 开源组件如 fastjson）。**漏洞自动分配依赖影响组件**：`defect_ledger` 保留 `affected_component` 列（漏洞创建/导入时写入，种子漏洞补全），组件为空时自动通道报错提示改用手动。

- **裂变拆单（定案）**：
  - 基础软件（命中 project）：公共中间件统一运维，**只拆 1 张运维单**（处置部门=**软件在 project 表配置的运维部门**，如 redis→技术部；`flow.maint_dept` 仅是导入时的归属部门，仅当软件未配置运维部门时兜底），不再按使用部门拆单；使用部门专员不可见该缺陷
  - 开源组件（命中 component_snapshot）：按依赖该组件的自研项目拆单（每项目 1 张，部门=项目所属部门）
- **自动匹配处理人**：处理人 = 涉及项目负责人主岗（project_owner 第一个）；**不再硬性校验「负责人部门==处置单部门」**（处理人来自涉及项目、处置单部门也来自同一项目，天然对应）；匹配不到（负责人为空）→ 单留空，由该部门专员/管理专员改派本部门单
- 内部缺陷（风险上报）：上报 → DF-xxx → 待审核 → 管理专员审核 → **手动指定一个或多个负责人**（每人一张处置单，不走裂变）

---

## 4. 漏洞情报（临时安全情报）

- 归**漏洞类**（flow_type=vuln），与缺陷（flow_type=defect）并列；**生成处置工单，不进入修复效能统计**
- 权限：开源治理专员/管理专员
- 入口：手动录入 + Excel 导入（必填列：组织/名称/漏洞编号/研判人/研判部门）
- **手动录入处置分配两通道（与缺陷评审一致）**：创建后不固定裂变——「自动分配」按组件类型裂变+匹配（软件→project / 组件→component_snapshot→project_inner）；「手动指定」输入一个或多个负责人账号（逗号/顿号/分号/空格分隔，每人一张处置单，部门=各自所属部门）。Excel 批量导入固定自动裂变（无手动通道）
- **自动分配预览可编辑（创建前）**：点「创建并裂变」先调 `/risk/vuln-intel/match` 预览待生成处置单清单——基础软件显示 软件名/运维部门/团队/运维负责人，组件显示 每个依赖项目/团队/部门/负责人；负责人可编辑（候选下拉+自定义输入），确认后调创建接口并回传 `assignOverride`（按 projectId+projectType 覆盖各处置单负责人），**以创建表单人编辑为准批量裂变**
- 幂等：漏洞编号唯一，重复录入/导入不重复生成
- 组织：开源软件/组件的组织命名空间（如 com.alibaba），存 `defect.org`

---

## 5. 缺陷修复效能看板

- **只统计外部缺陷（defect 表）及其处置工单**（defect_no ∈ defect），内部缺陷不纳入
- 缺陷级：新增缺陷按上游创建时间 created_at；等级=高危/中危/低危/提示/未知；部门取实体表 maint_dept
- 工单级：三态（进行中/逾期未闭环/已完成处置）、MTTR 从派发到完成、部门明细
- 后续邮件督办可复用本接口

---

## 6. 统一文件管理（attachment / bom）

- attachment 统一文件表（source/filename/md5/store_url，唯一键 (source, filename, store_url) 防重），生命周期流程（BOM/升级/退出）附件统一走它，source 区分业务（AttachmentSource 枚举：BOM_UPLOAD/VULDB_PLAN/VULDB_RESULT/FIRST_VERIFY/FINAL_VERIFY）
- bom 表只存元数据，store_url(bigint) 引用 attachment.id
- MinIO 目录规则：`{source}/{业务对象目录}/{文件名}`，BOM = `bom/{项目组织}/{项目名称}/{项目版本}/{文件名}`
- BOM 上传流程：uuid 查项目 → 组装目录 → 附件记录(新建/更新) → 传 MinIO → bom 记录(新建/更新) → 调外部解析服务 `/task/bomUpload(uuid, bomId)`（服务地址 config `services.bom-parser`，暂定 url1）

---

## 7. 组件快照定时清理（每小时）

- 目的：每小时归档旧数据，防 component_snapshot 膨胀
- 逻辑：
  1. 待清理 = `is_delete=1`（未删除）的 current_id
  2. 排除不清理：`project_inner.current_id`（当前最新）+ `project_inner_pro_info.is_deleted=0`（线上依赖）
  3. 逐个 commit：记录 `commit_id_history`（唯一键防重）→ 分批（可配 batchSize）显式列 `INSERT INTO ... SELECT` 归档到 `component_snapshot_history`（保留原主键）→ 置 `is_delete=0`
- cron 表达式 + 批次大小走 config（`scheduler.componentCleanupCron` / `cleanupBatchSize`）

---

## 8. 通用 HTTP 调用（httpclient）与定时任务框架

- `internal/httpclient`：通用传输层（BaseURL/超时/忽略证书校验/JSON/错误包装），`Auth` 认证策略（none/header/自定义闭包）+ `Registry` 服务注册表（config `services` 段按服务名配置）
- 三种访问方式统一：无认证（none）/ 固定 token 头（header）/ 签名认证（openapi 复用 httpclient 传输，签名作请求头注入）
- 定时任务：`robfig/cron/v3`，cron 表达式走 config（对应 Java 版配置定时表达式）；结构 = scheduler（注册触发）+ service（业务逻辑）
- 下游服务地址：配置文件静态配置（重启生效）

---

## 9. 映射表现状

| 映射表 | 维护界面 |
|---|---|
| 项目负责人（project_owner） | ✅ 有（我的运维 → 人员管理 → 编辑项目负责人） |
| 专员部门映射（dept_specialist_mapping） | ❌ 无（启动种子数据，后期由权限系统同步角色到表） |

---

## 10. 其他约定

- 缺陷状态就高不就低等级：待研判(1) < 忽略(2) < 待上会(3) < 已分配(4) < 已闭环(5)
- 组件快照历史表 = 原表去掉 is_delete，含 create_time
- 前端 operator / specialist 角色各有一份 MyOps / MyTask 页面（改动需同步两份）
