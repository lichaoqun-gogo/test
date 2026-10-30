# 审批流 & 缺陷/漏洞处置 回归测试用例

> 用途：本用例供回归测试手工执行使用。统一使用 CAS 身份体系，系统内用户唯一标识为 **CAS 用户名（username）**。
> 涉及两套相互独立的流程：
> 1. **串行审批流**（治理任务）：引入 / 升级 / 退出 / 风险上报。
> 2. **缺陷/漏洞处置流**（并行处置，与审批流无关）：研判 → 上会评审 → 裂变并自动匹配处理人（专员可在「我的运维-缺陷管理」缺陷行**展开行**处置工单「处理人」列的可编辑下拉补分配/改派，仅本部门工单可改）→ 完成/忽略 → 已完成。

---

## 一、测试账号（CAS，密码见 mock CAS users.json）

| 用户名 | 姓名 | 部门 | 团队 | 角色 | 用途 |
|--------|------|------|------|------|------|
| zhangsan | 张三 | 技术部 | 前端组 | operator | 普通员工：发起申请 / 技术部维护人员 |
| lisi | 李四 | 技术部 | 前端组 | specialist | 技术部开源治理专员（node_1 / 缺陷研判） |
| zhengjie | 郑杰 | 技术部 | 前端组 | specialist | 技术部第二开源治理专员（验证一个部门多名专员） |
| wangwu | 王五 | 平台研发部 | 后端组 | specialist | 平台研发部开源治理专员 |
| zhaoliu | 赵六 | 运维部 | 后端组 | specialist | 运维部开源治理专员 |
| sunqi | 孙七 | 安全部 | 设计组 | specialist | 安全部开源治理专员 |
| liuyang | 刘洋 | 平台研发部 | 后端组 | operator | 平台研发部项目维护人员（处置处理人） |
| chenming | 陈明 | 运维部 | 后端组 | operator | 运维部项目维护人员（处置处理人） |
| qianqi | 钱七 | 安全部 | 设计组 | operator | 安全部项目维护人员（处置处理人） |
| zhouba | 周八 | 开源治理办公室 | 运营组 | specialist | 管理部门专员（node_3 管理部门评审，具备全局治理视图权限） |
| wujiu | 吴九 | 开源治理办公室 | 运营组 | specialist | 管理部门第二专员（node_3 可评审，具备全局治理视图权限） |
| admin | 系统管理员 | 管理层 | 管理层 | admin | 系统管理员（看全部缺陷台账，默认进入全局治理视图） |

> 密码：普通账号 `123456`，admin 为 `admin123`（以 `mock认证/cas-test/users.json` 为准）。

---

## 二、部门-开源治理专员映射表（dept_specialist_mapping）

> 唯一键为 `(dept, specialist_role, specialist_id)`，**一个部门可配置 1 名或多名专员**；任一同部门同角色在岗专员（is_active=1）均可处理对应审批节点。节点展示审批人取最早配置的在岗专员。

| 部门 | specialist_id | 姓名 | specialist_role |
|------|---------------|------|-----------------|
| 技术部 | lisi | 李四 | dept_specialist |
| 技术部 | zhengjie | 郑杰 | dept_specialist |
| 平台研发部 | wangwu | 王五 | dept_specialist |
| 运维部 | zhaoliu | 赵六 | dept_specialist |
| 安全部 | sunqi | 孙七 | dept_specialist |
| 技术部 | lisi | 李四 | usage_dept_specialist |
| 技术部 | zhengjie | 郑杰 | usage_dept_specialist |
| 平台研发部 | wangwu | 王五 | usage_dept_specialist |
| 运维部 | zhaoliu | 赵六 | usage_dept_specialist |
| 安全部 | sunqi | 孙七 | usage_dept_specialist |
| 开源治理办公室 | zhouba | 周八 | mgmt_dept_specialist |
| 开源治理办公室 | wujiu | 吴九 | mgmt_dept_specialist |

---

## 三、开源软件-使用部门/运维部门映射表（project）

| 软件名 | 运维部门 maint_dept | 使用部门 use_dept | 许可证 |
|--------|---------------------|-------------------|--------|
| kafka | 技术部 | 平台研发部 | Apache-2.0 |
| redis | 技术部 | 技术部 | BSD-3-Clause |
| spring-boot | 平台研发部 | 技术部 | Apache-2.0 |
| nginx | 运维部 | 安全部 | BSD-2-Clause |
| node.js | 技术部 | 技术部 | MIT |

---

## 四、项目-维护人员映射表（project_owner，软件 project_type=software）

| 软件名 | 维护人员 username | 姓名 | 所属部门 | 团队 |
|--------|-------------------|------|----------|------|
| kafka | zhangsan | 张三 | 技术部 | 前端组 |
| kafka | liuyang | 刘洋 | 平台研发部 | 后端组 |
| redis | zhangsan | 张三 | 技术部 | 前端组 |
| spring-boot | liuyang | 刘洋 | 平台研发部 | 后端组 |
| spring-boot | zhangsan | 张三 | 技术部 | 前端组 |
| nginx | chenming | 陈明 | 运维部 | 后端组 |
| nginx | qianqi | 钱七 | 安全部 | 设计组 |
| node.js | zhangsan | 张三 | 技术部 | 前端组 |

---

## 五、缺陷情报源测试记录（defect_ledger，来源后续对接外部情报接口）

| defect_no | flow_type | source | title | risk_level | affected_component | maint_dept | cve_id | status |
|-----------|-----------|--------|-------|------------|--------------------|------------|--------|--------|
| DF-TEST-001 | defect | vuln-intel | Redis 远程代码执行漏洞（Lua 脚本） | critical | redis | 技术部 | CVE-2024-31449 | pending_analysis |
| DF-TEST-002 | defect | sca-scan | Redis 数据持久化配置缺陷 | high | redis | 技术部 | — | pending_review |
| DF-TEST-003 | vuln | vuln-intel | Node.js 原型污染漏洞 | high | node.js | 技术部 | CVE-2023-30581 | pending_analysis |
| DF-TEST-004 | vuln | vuln-intel | Spring Boot Actuator 配置泄露漏洞 | medium | spring-boot | 平台研发部 | CVE-2024-22243 | pending_analysis |
| DF-TEST-005 | defect | user-report | Nginx 日志路径穿越缺陷 | medium | nginx | 运维部 | — | pending_analysis |

> 来源枚举：`vuln-intel`（漏洞情报）、`sca-scan`（SCA 扫描）、`user-report`（用户上报）、`governance-tool`（外部治理工具导入，开发中，导入即待评审/忽略）。
> 覆盖点：缺陷路线（必须上会）、漏洞路线（可直接 need_fix 裂变 / 可 ignore）、跨部门裂变（nginx：运维部+安全部）。

---

## 六、公共说明

### 6.1 串行审批流节点与跳过规则

| 节点 | 名称 | 审批人 | 跳过条件 |
|------|------|--------|----------|
| node_0 | 发起节点 | 申请人 | 提交时自动置为已通过 |
| node_1 | 本部门开源治理专员复审 | 申请人部门的 dept_specialist | 申请人 == 本部门专员 |
| node_2 | 使用部门确认 | 使用部门的 usage_dept_specialist | 仅升级/退出激活；使用部门 == 申请人部门（运维部门）时跳过 |
| node_3 | 管理部门评审 | zhouba（开源治理办公室 mgmt_dept_specialist） | 不跳过 |

### 6.2 状态机（中英文对照）

系统包含两条相互独立的流程线：

**流程一：串行审批流**（治理任务 governance_task：引入 / 升级 / 退出 / 风险上报）

| 英文状态 | 中文 | 说明 |
|---|---|---|
| draft | 草稿 | 保存未提交；被驳回/撤回后回到草稿 |
| approving | 审批中 | 已提交，流转中；可撤回 |
| approved | 已通过 | 全部节点通过 |
| rejected | 已驳回 | 任一节点驳回，回到草稿可修改重报 |
| withdrawn | 已撤回 | 审批中主动撤回，回到草稿 |

**流程二：缺陷/漏洞处置流**（与审批流无关：台账中的缺陷记录为"父"，研判/评审后跨部门裂变出多张处置单"子"**并行处置**，全部子单完成后父缺陷更新为「已完成」）

| 层级 | 英文状态 | 中文 | 说明 |
|---|---|---|---|
| 缺陷（父） | pending_analysis | 待研判 | 平台内录入/接口创建后的初始状态 |
| | pending_review | 待上会评审 | 研判结论为需上会（缺陷类必上会）；**外部治理工具导入（`source=governance-tool`，开发中）的缺陷情报直接进入该状态，导入即待评审、不裂变** |
| | assigned | 待处置 | 裂变出处置单并已自动匹配处理人，等待处置 |
| | ignored | 已忽略 | 研判认定无需处置（如不在使用范围） |
| | closed | 已完成 | 全部处置单完成/忽略 |
| 处置单（子） | pending_assign | 待处置 | 已裂变；处理人已由系统自动匹配或待专员在缺陷展开行补分配/改派 |
| | processing | 处理中 | 处理人开始处理（兼容历史数据） |
| | completed | 已完成 | 处置动作完成（升级/替换/加固等） |
| | na | 忽略 | 处置单对应项目/部门不适用（原「不适用」，UI 统一显示「忽略」） |

### 6.3 接口约定

- 所有业务接口前缀 `/api/v1`，需登录态。CAS 模式下需携带 JWT；手工测试建议直接通过前端登录后触发，或手动补请求头：
  - `X-User-ID`：CAS 用户名（如 `zhangsan`）
  - `X-User-Name`：姓名（URL 编码，如 `%E5%BC%A0%E4%B8%89`）
  - `X-User-Role`：`admin` / `specialist` / `operator`
  - `X-User-Dept`：部门（URL 编码，如 `%E6%8A%80%E6%9C%AF%E9%83%A8`）
- 审批流接口：
  - `POST /workflow/tasks/:id/submit` 提交
  - `POST /workflow/tasks/:id/approve` body `{ "nodeId": <node_id>, "opinion": "..." }`
  - `POST /workflow/tasks/:id/reject` body `{ "nodeId": <node_id>, "opinion": "..." }`
  - `POST /workflow/tasks/:id/withdraw` 撤回
  - `GET /workflow/tasks/:id/nodes` / `GET /workflow/tasks/:id/current-node` / `GET /workflow/tasks/:id/history`
- 缺陷处置接口：见「七、缺陷/漏洞处置用例」。

### 6.4 缺陷处置流程与权限（两次分配 + 自动匹配）

处置流：**第一次分配（研判人）→ 软件负责人研判 → 上会评审 → 按涉及项目裂变并自动匹配各涉及项目主岗为处理人 → 处置 → 父缺陷「已完成」**

**缺陷情报入口（两种）**：
- **平台内录入/接口创建**：`pending_analysis`（待研判），需先分配研判人 → 研判 → 待上会评审
- **外部治理工具导入（`source=governance-tool`，开发中）**：各部门专员/管理部门专员用治理工具在互联网侧研判，将"需上会"结论的情报包导入内网平台，**导入即 `pending_review`（待上会评审，不裂变）；结论为忽略的误报项导入为 `ignored` 保留研判台账**——从「待上会评审」起的后半段流程与下述完全一致，无代码差异

**裂变拆单粒度 = 涉及项目**（一张处置单 = 一个缺陷 × 一个涉及项目及所属部门）：
- 基础软件（如 redis，`project` 表命中）：软件本身一张**运维单**（`maint_dept`）+ 各使用部门各一张**使用单**（`use_dept`，项目均为该软件）
- 开源组件（如 fastjson，`project` 表未命中）：`component_snapshot` 查依赖该组件的自研项目（`project_inner`），**每个依赖项目一张使用单**（部门 = 项目所属部门）——**同一部门有多个依赖项目即拆出多张处置单**（如某部门 4 个项目都用 fastjson → 该部门 4 张使用单），不再按部门去重合并

**两种分配方式**（验证阶段为显式按钮；后期可改为系统自动执行，去掉按钮）：

| 场景 | 手动分配（按钮） | 自动分配（系统逻辑） |
|------|------------------|----------------------|
| 第一次分配（研判人） | 「分配研判人」：弹窗展示自动匹配候选（第一个为主岗），可改选或手动输入任意账号 | 「自动分配研判人」：系统直接取候选第一个（主岗） |
| 处置单处理人 | 裂变时系统已自动匹配各涉及项目主岗；专员在「我的运维 → 缺陷管理」缺陷**展开行**的处置工单表「处理人」列**直接改选处理人**（仅本部门处置单可编辑下拉，点击弹出本部门人员列表，已分配/未分配均可更改；其他部门工单只读） | 「一键自动分配」：对仍无处理人的处置单重新匹配涉及项目主岗（仅该缺陷运维部门可见） |

**自动匹配规则**（候选来源，主岗 = 列表第一个）：
- 开源基础软件（如 redis/nginx，`project` 表命中）：该软件的 `project_owner`（software）
- 开源组件（如 fastjson，`project` 表未命中）：`component_snapshot` 查依赖该组件的自研项目（`project_inner`）→ 各项目 `project_owner`（inner）

**操作入口（UI）**：
- 分配研判人 / 自动分配研判人：「我的运维 → 缺陷管理」tab → 缺陷行点「详情」→ 详情对话框内操作（列表有「研判人」列，未分配显示「未分配」标签；专员详情中**无「研判」按钮**，研判只能由被分配的研判人本人在「待我研判」操作）
- 待我研判（**仅被分配的研判人本人**，专员不再代研判）：「我的任务 → 待我研判」tab → 该行点「研判」
- **评审（上会）/ 添加会议结论**：仅全局治理视图「风险台账 → 缺陷台账」→ 缺陷行点「详情」→ 点「评审（上会）」（只有 admin 与管理部门专员 zhouba/wujiu 可见并有权操作；部门专员如李四的「我的运维-缺陷管理」详情中不显示该按钮）
- **处置工单（展开行）**：「我的运维 → 缺陷管理」tab → 缺陷行前**点箭头展开** → 展开区域展示该缺陷的**处置工单表**（工单编号/**涉及项目**/处置部门/类型/处理人/状态，涉及项目 = 软件名或自研项目编码，用于区分同部门多张单）。详情对话框内**不再展示处置工单**，仅展示缺陷信息 + 历史记录 + 分配研判人按钮
- **补分配 / 改派处理人（统一入口）**：缺陷展开行的「处理人」列为**可编辑下拉**（仅「待处置」工单且为本部门处置单可编辑）：已分配显示姓名、点开即可改选；未分配显示「请选择处理人」占位、点开同样弹出下拉；选中即生效（补分配/改派）。**仅本部门处置单可改派，其他部门工单只读显示姓名**（后端 `Assign` 亦校验本部门专员）。展开区域右上角保留「一键自动分配」按钮（仅该缺陷运维部门可见）
- 待我处置工单（完成 / 忽略）：「我的任务 → 待我处置」tab（**纯处置工单列表**：工单编号/**涉及项目**/处置部门/处置人/处置类型/状态/截止时间/操作；「处置人」列只读展示处理人姓名，无处理人显示「-」；改派请前往上方「我的运维 → 缺陷管理」展开行操作）
- 历史记录：「我的运维 → 缺陷管理」或全局治理「缺陷台账」→ 缺陷行点「详情」→ 详情底部「历史记录」区域：按时间倒序展示该缺陷全部操作（某人 何时 何物 状态如何调整：创建 / 分配研判人 / 自动分配研判人 / 研判 / 评审 / 自动分配处理人 / 分配处理人 / 改派 / 完成 / 忽略 / 闭环），含每条记录的操作人、时间、动作、状态流转（处置工单与父缺陷状态分别标注）与备注

**缺陷台账可见性（「我的运维 → 缺陷管理」列表，缺陷全景图）**：
- admin / 管理部门专员：全部缺陷（**全景**）
- 部门专员：`maint_dept=本部门` 的缺陷，**或本部门有处置单的缺陷**（跨部门裂变后使用部门专员也可看到该缺陷，以便展开改派本部门处置单）
- **展开行的处置单可见性（读全景）**：**能看见该缺陷的用户即可查看该缺陷的全部处置单**——运维部门专员 / 使用部门专员 / 管理部门专员均看**全部处置单（只读展示）**
- **改派权限（统一）**：仅**本处置单所属部门专员**可改派（后端 `Assign` 校验本部门专员）；**管理部门专员不可跨部门改派**，`DefectLedger`（缺陷台账）处置工单表仅只读展示处理人；未匹配到主岗的工单由对应部门专员在展开行补分配

**操作权限**：

| 动作 | 谁操作 |
|------|--------|
| 分配研判人 / 自动分配研判人 | 专员（admin / 管理部门专员 zhouba、wujiu / 缺陷 `maint_dept` 的部门专员） |
| 研判 | **仅被分配的研判人本人**（dispatchTarget；admin 可代研判，专员不再兜底代研判） |
| 上会评审（添加会议结论） | **仅 admin / 管理部门专员**（全局治理视图「缺陷台账」操作；后端接口同权限校验） |
| 补分配 / 改派 / 一键自动分配（处置单处理人） | 对应部门的专员（`dept_specialist` / `usage_dept_specialist`） |
| 完成 / 忽略 | **仅处理人本人**（assigneeName/assigneeUsername；admin 可代操作，专员不再兜底代处置） |

**接口对照**（UI 与接口等价，可任选验证方式）：

- `GET /risk/defects/:id/analyst-candidates` 自动匹配候选（第一个为主岗）
- `POST /risk/defects/:id/assign-analyst` body `{"analystUsername":"zhangsan"}` 手动分配研判人
- `POST /risk/defects/:id/auto-assign-analyst` 自动分配研判人（取主岗）
- `POST /risk/defects/:id/analyze` body `{"conclusion":"need_review|need_fix|ignore","analyzeResult":"..."}` 研判
- `POST /risk/defects/:id/review` body `{"reviewConclusion":"need_fix|no_action","reviewMeetingNo":"...","remark":"...","deadline":"YYYY-MM-DD HH:mm:ss"}` 上会评审
- `GET /risk/defects/:id/disposals` / `GET /risk/defects/:id/history`
- `POST /risk/defects/:id/auto-assign-disposals` 一键自动分配处置单（对无处理人的工单匹配主岗）
- `POST /risk/disposals/:id/assign` 补分配/改派处理人 body `{"assigneeUsername":"...","assigneeName":"..."}`
- `POST /risk/disposals/:id/complete` body `{"ticketNo":"...","note":"..."}`
- `POST /risk/disposals/:id/na` body `{"remark":"..."}`

### 6.5 权限档案与全局治理视图（对接外部权限系统的数据契约）

登录后前端调用 `GET /api/v1/auth/me` 获取用户信息与权限档案，据此渲染可见菜单和按钮（切换入口、全局治理菜单等）：

```json
{
  "code": 0,
  "data": {
    "user": { "username": "zhouba", "nickname": "周八", "role": "specialist", "department": "开源治理办公室", "team": "运营组" },
    "roles": ["specialist", "mgmt_specialist"],
    "permissions": ["global_governance_view"]
  }
}
```

本地授权规则（等价于外部权限系统中的配置，后续对接时仅替换 `internal/service/permission` 实现，契约不变）：

| 条件 | 授予 |
|------|------|
| `users.role = 'admin'` | `global_governance_view` |
| `dept_specialist_mapping` 中为有效 `mgmt_dept_specialist` | `mgmt_specialist` 角色 + `global_governance_view` |
| `users.can_toggle_admin_view = 1`（模拟外部权限系统单独授权，可直接改库） | `global_governance_view` |

视图行为：

- 具备 `global_governance_view` 的用户：顶部显示「个人工作台 / 全局治理」切换下拉框；admin 默认进入全局治理视图，其余（含管理部门专员）默认个人工作台；切换/刷新后保持上次视图。
- 不具备该权限的用户：无切换按钮，仅个人工作台。

### 6.6 菜单接口（角色并集模型）

`GET /api/v1/auth/menus` **无参数**，后端按用户全部角色一次返回菜单并集（每角色一棵树）：

```json
{
  "code": 0,
  "data": {
    "roles": ["specialist", "admin"],
    "menus": { "specialist": [ ...菜单树... ], "admin": [ ...菜单树... ] }
  }
}
```

模型说明：一个用户可具备多个角色，**菜单可见性 = 全部角色菜单的并集**（如角色 A 有菜单 1,2,3、角色 B 有菜单 4,5,6，则登录后 1~6 全部可见）。前端一次拉取并注册全部路由，视图切换只是按角色分组展示（个人工作台显示自身角色菜单，全局治理显示 admin 菜单）。

菜单角色本地映射规则（菜单-角色映射维护在本地 `role_menus` 表，**对接外部权限系统时仅需把"角色集合"替换为接口返回值**）：

| 用户 | 菜单角色集合 |
|------|--------------|
| 张三（operator） | `["operator"]` |
| 李四（specialist） | `["specialist"]` |
| 周八/吴九（specialist + 管理部门专员） | `["specialist", "admin"]` |
| admin | `["admin", "specialist"]`（管理员同样有个人工作台） |

---

## 七、测试用例

> 说明：用例中的「预期审批人」依据第二节映射推导。「预期结果」请以实际状态流转、节点 skip_reason、审批历史为准。
> 推荐每次用例使用新任务（新建任务），避免历史数据干扰。以下用例中的 `{id}`、`{nodeId}` 以接口实际返回为准。

---

### 【审批流】TC1 引入申请（普通员工发起，跨完整链路）

- 申请人：zhangsan（张三，技术部）
- 业务类型：introduce
- 项目（使用部门）：不参与 node_2（引入类自动跳过）

| 步骤 | 操作人 | 操作 | 接口/操作 | 预期结果 |
|------|--------|------|-----------|----------|
| 1 | zhangsan | 创建引入任务（草稿） | `POST /api/v1/tasks`，body：`{"type":"introduce","title":"引入测试组件","status":"draft","task_detail":{"component_name":"kafka","component_version":"3.5.0","introduce_reason":"测试引入","project_dept":"平台研发部"}}` | 返回 `task_no`、`id`，状态 `draft` |
| 2 | zhangsan | 提交审批 | `POST /api/v1/workflow/tasks/{id}/submit` | 状态 `approving`，当前节点=本部门开源治理专员复审（李四） |
| 3 | zhangsan | 查节点 | `GET /api/v1/workflow/tasks/{id}/nodes` | node_0 已通过；node_1 审批人 lisi；node_2 is_skipped=true（引入类跳过）；node_3 审批人 zhouba |
| 4 | lisi | 审批通过 | `POST /api/v1/workflow/tasks/{id}/approve` body `{"nodeId":<node_1的id>,"opinion":"同意"}` | 流转到管理部门评审（周八） |
| 5 | zhouba | 审批通过 | `POST /api/v1/workflow/tasks/{id}/approve` body `{"nodeId":<node_3的id>,"opinion":"同意"}` | 状态 `approved`，`assignedTo` 置空 |

**验证点**：node_2 被跳过（skip_reason 含「引入」/「不需要使用部门确认」）；最终 `approved`。

---

### 【审批流】TC2 引入申请（专员本人发起，跳过 node_1）

- 申请人：lisi（李四，技术部专员）

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | lisi | 创建并提交引入任务 | `POST /api/v1/tasks`（type=introduce, status=draft）→ `POST /workflow/tasks/{id}/submit` | 状态 `approving` |
| 2 | lisi | 查节点 | `GET /workflow/tasks/{id}/nodes` | node_1 is_skipped=true（发起人=部门专员）；node_2 跳过；当前节点=node_3（周八） |
| 3 | zhouba | 审批通过 | `POST /workflow/tasks/{id}/approve` | 状态 `approved` |

**验证点**：node_1 skip_reason 含「发起人即为部门专员」；任务提交后直接到 node_3。

---

### 【审批流】TC3 升级申请（跨部门，全节点）

- 申请人：zhangsan（张三，技术部）
- 业务类型：upgrade
- 组件：kafka（运维=技术部，使用=平台研发部）

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | zhangsan | 创建升级任务 | `POST /api/v1/tasks`，body：`{"type":"upgrade","title":"升级kafka","status":"draft","task_detail":{"component_name":"kafka","current_version":"3.4.0","target_version":"3.6.0","upgrade_reason":"安全修复","project_dept":"平台研发部"}}` | 返回 `id` |
| 2 | zhangsan | 提交 | `POST /workflow/tasks/{id}/submit` | 当前节点=node_1（李四） |
| 3 | lisi | 通过 node_1 | `POST /workflow/tasks/{id}/approve` | 流转到 node_2（使用部门确认，王五） |
| 4 | wangwu | 通过 node_2 | `POST /workflow/tasks/{id}/approve` | 流转到 node_3（周八） |
| 5 | zhouba | 通过 node_3 | `POST /workflow/tasks/{id}/approve` | 状态 `approved` |

**验证点**：node_2 不跳过，审批人为使用部门（平台研发部）专员王五；全链路 4 节点。

---

### 【审批流】TC4 退出申请（同部门，跳过 node_2）

- 申请人：zhangsan（张三，技术部）
- 业务类型：withdraw
- 组件：redis（运维=技术部，使用=技术部）

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | zhangsan | 创建退出任务 | `POST /api/v1/tasks`，body：`{"type":"withdraw","title":"退出redis","status":"draft","task_detail":{"component_name":"redis","component_version":"6.2.0","withdraw_reason":"不再使用","project_dept":"技术部"}}` | 返回 `id` |
| 2 | zhangsan | 提交 | `POST /workflow/tasks/{id}/submit` | 当前节点=node_1（李四） |
| 3 | lisi | 通过 node_1 | `POST /workflow/tasks/{id}/approve` | 跳过 node_2，直接到 node_3（周八） |
| 4 | zhouba | 通过 node_3 | `POST /workflow/tasks/{id}/approve` | 状态 `approved` |

**验证点**：node_2 is_skipped=true，skip_reason 含「使用部门与运维部门相同」或「所有使用部门与运维部门相同」。

---

### 【审批流】TC5 风险上报（defect 任务，node_2 跳过）

- 申请人：zhangsan（张三，技术部）
- 业务类型：defect（前端「风险上报」）

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | zhangsan | 创建风险上报任务 | `POST /api/v1/tasks`，body：`{"type":"defect","title":"发现高风险依赖","status":"draft","task_detail":{"risk_description":"某依赖存在已知漏洞"}}` | 返回 `id` |
| 2 | zhangsan | 提交 | `POST /workflow/tasks/{id}/submit` | 状态 `approving`，当前节点=node_1（李四） |
| 3 | lisi | 通过 node_1 | `POST /workflow/tasks/{id}/approve` | 流转到 node_3（周八） |
| 4 | zhouba | 通过 node_3 | `POST /workflow/tasks/{id}/approve` | 状态 `approved` |

**验证点**：node_2 跳过（defect 类型非升级/退出）。

---

### 【审批流】TC6 撤回 / 驳回 / 重新提交

- 申请人：zhangsan；审批人：lisi

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | zhangsan | 创建并提交引入任务 | `POST /tasks` → `POST /workflow/tasks/{id}/submit` | 状态 `approving`，当前节点 lisi |
| 2 | zhangsan | 撤回 | `POST /workflow/tasks/{id}/withdraw` | 状态 `draft` |
| 3 | zhangsan | 重新提交 | `POST /workflow/tasks/{id}/submit` | 状态 `approving`，节点复用并重置 |
| 4 | lisi | 驳回 | `POST /workflow/tasks/{id}/reject` body `{"nodeId":<node_1>,"opinion":"材料不全"}` | 状态 `rejected`，当前节点「已驳回」 |

**验证点**：撤回后历史含 `withdraw`；驳回后为终态 `rejected`。

---

**测试数据重置**（用例使用过的缺陷再次回归前执行，以 DF-TEST-001 为例，替换编号即可）：

```sql
UPDATE defect_ledger SET status='pending_analysis',
  dispatch_type='', dispatch_target='', dispatch_dept='', dispatch_time=NULL, dispatcher='',
  analyze_by='', analyze_time='', review_meeting_no='', review_by='', review_time='', deadline=NULL
WHERE defect_no='DF-TEST-001';
DELETE FROM defect_disposal WHERE defect_id=(SELECT id FROM defect_ledger WHERE defect_no='DF-TEST-001');
DELETE FROM defect_status_history WHERE defect_id=(SELECT id FROM defect_ledger WHERE defect_no='DF-TEST-001');
```

DF-TEST-002（回到「待上会评审」，用于从待评审直接测后半段流程；**预置研判人 zhangsan/张三**，保证「待上会评审必有研判人」规则自洽）：

```sql
UPDATE defect_ledger SET status='pending_review',
  dispatch_type='analyze', dispatch_target='zhangsan', dispatch_dept='技术部',
  dispatch_time=NOW(), dispatcher='lisi',
  analyze_by='张三', analyze_time=NOW(),
  review_meeting_no='', review_by='', review_time='', deadline=NULL
WHERE defect_no='DF-TEST-002';
DELETE FROM defect_disposal WHERE defect_id=(SELECT id FROM defect_ledger WHERE defect_no='DF-TEST-002');
DELETE FROM defect_status_history WHERE defect_id=(SELECT id FROM defect_ledger WHERE defect_no='DF-TEST-002');
```

---

### 【缺陷处置】TC7 缺陷全流程·手动分配路线（redis）

- 数据：DF-TEST-001（defect，redis，技术部，待研判）
- 涉及账号：lisi（技术部专员，负责分配研判人/改派处理人）、zhangsan（redis 主岗，负责研判/处置）、zhouba（管理部门专员，负责评审上会）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | lisi | 顶部菜单「我的运维」→「缺陷管理」tab → DF-TEST-001 行点「详情」 | 弹出详情，顶部按钮区显示「分配研判人」「自动分配研判人」（已分配后消失） |
| 2 | lisi | 详情中点「分配研判人」 | 弹窗列出候选：张三（带绿色「主岗」标签）；张三默认已选中 |
| 3 | lisi | 候选保持张三，点「确认分配」 | 提示"已分配研判人"；列表刷新后该缺陷研判人为 zhangsan，状态仍「待研判」 |
| 4 | zhangsan | 切换登录 →「我的任务」→「待我研判」tab → DF-TEST-001 行点「研判」→ 结论选「需上会评审」→ 填研判说明 → 点「确认研判」 | 状态变为「待上会评审」，研判人记录为张三 |
| 5 | zhouba | 切换登录 → 顶部视图切换为「全局治理」→「风险台账 → 缺陷台账」→ DF-TEST-001 行点「详情」→ 点「评审（上会）」→ 结论选「需处置」→ 填会议号（如 M-2026-001）和处置期限 → 点「确认评审」 | 父缺陷状态变「待处置」；缺陷行前出现箭头，点开展开区有「处置工单」表：1 张技术部单（涉及项目：redis）、状态「待处置」、处理人已自动匹配为张三（redis 主岗） |
| 6 | lisi | 切换登录 →「我的运维」→「缺陷管理」→ DF-TEST-001 行点箭头展开 → 处置工单该行「处理人」列为可编辑下拉（显示张三）→ 点开选其他技术部人员再改回张三 | 提示"处理人已更新"；处理人按选择变化，工单状态仍「待处置」；历史记录新增「分配处理人/改派」记录 |
| 7 | zhangsan | 切换登录 →「我的任务」→「待我处置」tab | 列表出现该工单（工单编号 DP-xxx、处置部门技术部、处置人列显示张三、状态待处置；纯处置工单列表，无缺陷编号/标题列） |
| 8 | zhangsan | 「待我处置」中该行点「完成」→ 填工单号（如 INC-001）和处置说明 → 提交 | 工单状态变「已完成」；父缺陷自动变「已完成」（列表状态列） |

**验证点**：未分配研判人时详情中显示两个分配按钮；研判只能由被分配的研判人本人操作（专员详情中无「研判」按钮、不再代研判）；**评审按钮仅出现在全局治理「缺陷台账」，李四「我的运维-缺陷管理」详情中无评审按钮**；评审「需处置」确认时系统即自动匹配该部门主岗为处理人（工单与父缺陷状态一致，均为「待处置」）；改派/补分配统一在「我的运维-缺陷管理」缺陷行**展开行**处置工单「处理人」列的可编辑下拉操作（已分配/未分配均可点开改选），「我的任务-待我处置」处置人列只读展示；全部工单完成后父缺陷变「已完成」；缺陷列表「研判人」列显示张三（未分配时显示「未分配」）；详情底部「历史记录」按时间倒序展示：创建（情报接入）→ 分配研判人（李四→张三）→ 研判（张三，待研判→待上会评审）→ 评审（周八，→待处置）→ 自动分配处理人（周八，系统自动匹配张三）→ 改派（李四）→ 完成（张三，待处置→已完成）→ 自动闭环（→已完成），每条含操作人、时间、动作、状态流转与备注。

---

### 【缺陷处置】TC8 评审结论「无需处置」（no_action）

- 数据：DF-TEST-002（defect，kafka，技术部，**待上会评审，已预置**）
- 账号：zhouba（管理部门专员，负责评审）
- 说明：**该缺陷预置为「待上会评审」，即「从待评审开始」直接测后半段流程的入口数据**——对应外部治理工具导入（结论=需上会）后的状态，与平台内研判后（TC7 步骤 1-4）的状态完全等价；本用例验证评审「无需处置」分支
- **规则：进入「待上会评审」必然已完成研判，研判人必有值**——预置数据含研判人 zhangsan（张三），台账列表/详情「研判人」列显示张三，**不会出现「待上会评审 + 研判人未分配」**（前端研判人列对已研判状态回退显示实际研判人 `analyze_by`，仅待研判未分配时显示「未分配」）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | zhouba | 全局治理视图 →「风险台账 → 缺陷台账」→ DF-TEST-002「详情」→ 点「评审（上会）」→ 结论选「无需处置」→ 确认评审 | 状态直接变「已忽略」，不生成处置工单 |

---

### 【漏洞处置】TC9 漏洞全流程·自动分配路线（node.js）

- 数据：DF-TEST-003（vuln，node.js，技术部，待研判）
- 账号：lisi（分配）、zhangsan（node.js 主岗，研判/处置）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | lisi | 「我的运维」→「缺陷管理」→ DF-TEST-003「详情」→ 点「自动分配研判人」 | 提示"已自动分配研判人 张三（主岗）"，无需选人 |
| 2 | zhangsan | 登录 →「我的任务」→「待我研判」→ DF-TEST-003 行点「研判」→ 结论选「直接处置」（漏洞可跳过上会）→ 填研判说明 → 确认研判 | 父缺陷状态变「待处置」，缺陷行前出现箭头，点开展开区出现处置工单：1 张技术部单、状态「待处置」、处理人已自动匹配为张三（node.js 主岗） |
| 3 | zhangsan | 「我的任务」→「待我处置」→ 该工单点「完成」→ 填工单号 INC-002 → 提交 | 工单「已完成」，父缺陷自动「已完成」 |

**验证点**：漏洞研判可直接处置（跳过上会）但分配研判人步骤不可跳过；研判直接处置时系统即自动匹配该部门主岗为处理人（工单/父缺陷均为「待处置」）；专员可在「我的运维-缺陷管理」缺陷行**展开行**处置工单「处理人」下拉中改选处理人。

---

### 【漏洞处置】TC10 漏洞忽略 + 多负责人主岗（spring-boot）

- 数据：DF-TEST-004（vuln，spring-boot，平台研发部，待研判）
- 账号：wangwu（平台研发部专员）、liuyang（spring-boot 主岗）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | wangwu | 「我的运维」→「缺陷管理」→ DF-TEST-004「详情」→ 点「分配研判人」 | 弹窗候选两个：刘洋（主岗）、张三——刘洋默认选中（spring_boot 第一个 owner） |
| 2 | wangwu | 保持刘洋，点「确认分配」 | 研判人 liuyang |
| 3 | liuyang | 登录 →「我的任务」→「待我研判」→ DF-TEST-004 行点「研判」→ 结论选「无需处置」→ 确认研判 | 状态变「已忽略」（漏洞允许直接忽略；缺陷则必须上会） |

---

### 【缺陷处置】TC12 跨部门裂变 + 改派 + 忽略闭环（nginx）

- 数据：DF-TEST-005（defect，nginx，运维部维护/安全部使用，待研判）
- 账号：zhaoliu（运维部专员）、chenming（nginx 主岗）、qianqi（nginx 负责人）、sunqi（安全部使用专员）、zhouba（管理部门专员，负责评审上会）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | zhaoliu | 「我的运维」→「缺陷管理」→ DF-TEST-005「详情」→ 点「分配研判人」 | 候选：陈明（主岗）、钱七 |
| 2 | zhaoliu | 选陈明 → 确认分配 | 研判人 chenming |
| 3 | chenming | 登录 →「我的任务」→「待我研判」→ DF-TEST-005 行点「研判」→「需上会评审」→ 确认 | 状态「待上会评审」 |
| 4 | zhouba | 切换登录 → 全局治理视图 →「风险台账 → 缺陷台账」→ DF-TEST-005「详情」→「评审（上会）」→「需处置」→ 确认评审 | 父缺陷「待处置」，缺陷行前出现箭头，点开展开区可见处置工单 2 张均「待处置」（涉及项目均为 nginx）：运维部单处理人已自动匹配陈明（nginx 主岗）；安全部单处理人空（未匹配到安全部主岗，「处理人」列显示「请选择处理人」占位） |
| 5 | zhaoliu | 切回 →「我的运维」→「缺陷管理」→ DF-TEST-005 行点箭头展开 → 展开区「处置工单」表中运维部单「处理人」下拉点开选钱七再改回陈明 | 提示"处理人已更新"；运维部单处理人按选择变化、状态仍「待处置」 |
| 6 | sunqi | 登录 →「我的运维」→「缺陷管理」→ DF-TEST-005 行点箭头展开 → 展开区「处置工单」表中安全部单（处理人空）「处理人」下拉点开选孙七 | 安全部单处理人变孙七，状态仍「待处置」 |
| 7 | chenming | 「我的任务」→「待我处置」→ 运维部单点「完成」→ 填 INC-003 → 提交 | 运维单「已完成」，父缺陷仍「待处置」 |
| 8 | sunqi | 「我的任务」→「待我处置」→ 安全部单点「忽略」→ 确认 | 安全部单「忽略」；父缺陷自动「已完成」（工单全部完结即完成） |

**验证点**：跨部门裂变（维护+使用两张单）；评审时自动匹配各部门主岗；改派/补分配统一在「我的运维-缺陷管理」缺陷行**展开行**处置工单「处理人」列的可编辑下拉操作，未匹配/已匹配的工单专员均可点开补分配/改派处理人；**可见性**：sunqi（安全部使用专员）虽非该缺陷运维部门，但因其部门有处置单可看到 DF-TEST-005（列表放开「本部门有处置单的缺陷」），且展开行**看全部两张单（读全景）**，但只有安全部单「处理人」列为可编辑下拉（运维部单只读）；zhaoliu（运维部专员）展开行同样看**全景**（运维部+安全部两张单），安全部单「处理人」列只读；管理部门专员在 DefectLedger 看全景但**不可跨部门改派**；某单「忽略」同样计入完结，父缺陷「已完成」。

---

### 【缺陷处置】TC17 组件缺陷全链路·按涉及项目裂变（fastjson）

- 数据：DF-TEST-006（defect，**组件 fastjson**，affected_component_type=component，技术部，待研判）
- 依赖链：fastjson 被 3 个自研项目依赖（`component_snapshot → project_inner`）：ORD-CENTER 订单中心（技术部，主岗张三）、PAY-PLATFORM 支付平台（平台研发部，主岗刘洋）、ORDER-SYS 订单系统（技术研发部，**未配置项目负责人**）
- 账号：lisi（技术部专员）、zhangsan（订单中心主岗）、liuyang（支付平台主岗）、wangwu（技术研发部专员）、zhouba（管理部门专员，全局可见）

| 步骤 | 登录账号 | 页面与操作 | 预期结果 |
|------|----------|------------|----------|
| 1 | lisi | 「我的运维」→「缺陷管理」→ DF-TEST-006「详情」→ 点「分配研判人」 | 候选：张三（主岗，订单中心）、刘洋（支付平台）——组件通过「组件依赖表 → 自研项目 → 项目负责人」匹配，非 project 表；订单系统无负责人不进入候选 |
| 2 | lisi | 保持张三（主岗）→ 确认分配 | 研判人 zhangsan |
| 3 | zhangsan | 登录 →「我的任务」→「待我研判」→ DF-TEST-006 行点「研判」→「需上会评审」→ 确认 | 状态「待上会评审」 |
| 4 | zhouba | 切换登录 → 全局治理视图 →「风险台账 → 缺陷台账」→ DF-TEST-006「详情」→「评审（上会）」→「需处置」→ 填会议号 M-2026-004 → 确认评审 | 父缺陷「待处置」；**按涉及项目拆出 3 张处置单**均「待处置」，展开行处置工单表「涉及项目」列分别为：订单中心/支付平台/订单系统：订单中心单（技术部）处理人已自动匹配张三、支付平台单（平台研发部）已自动匹配刘洋、订单系统单（技术研发部）处理人空（该项目未配置负责人，待技术研发部专员补分配） |
| 5 | wangwu | 登录 →「我的运维」→「缺陷管理」→ DF-TEST-006 行点箭头展开 → 处置工单表中技术研发部单（涉及项目：订单系统）「处理人」列点开选技术研发部人员 | 技术研发部单处理人已指定，状态仍「待处置」——技术研发部专员虽非该缺陷运维部门，但因其部门有处置单可看到缺陷并改派本部门单 |
| 6 | zhangsan / liuyang | 分别登录 →「我的任务」→「待我处置」→ 各自单「完成」→ 填 INC-006/INC-007 → 提交 | 订单中心单、支付平台单「已完成」；父缺陷仍「待处置」（订单系统单未完成） |
| 7 | wangwu | 「我的任务」→「待我处置」→ 技术研发部单「完成」→ 填 INC-008 → 提交 | 技术研发部单「已完成」，父缺陷自动「已完成」 |

**验证点**：组件缺陷按「涉及项目」裂变（每依赖项目一张使用单，跨 3 部门 3 张，不再按部门合并）；`project_owner` 缺失时该单不自动分配、由对应部门专员补分配；运维部门专员 / 使用部门专员（技术部、平台研发部、技术研发部）/ 管理部门专员均可查看该缺陷**全部处置单（读全景）**，但仅能改派**本部门**处置单（管理部门专员不可跨部门改派）。

---

### 【缺陷处置】TC11 缺陷台账权限过滤

- 使用 DF-TEST-005（defect，nginx，maint=运维部，pending_analysis）

| 步骤 | 操作人 | 操作 | 接口 | 预期结果 |
|------|--------|------|------|----------|
| 1 | lisi（技术部） | 查台账 | `GET /risk/defects` | 仅返回 maint_dept=技术部 的记录，不含 DF-TEST-005 |
| 2 | zhaoliu（运维部） | 查台账 | `GET /risk/defects` | 仅返回 maint_dept=运维部 的记录，含 DF-TEST-005 |
| 3 | admin 或 zhouba | 查台账 | `GET /risk/defects` | 返回全部记录 |

**验证点**：非 admin/管理部门专员按 `maint_dept` 过滤。

---

### 【权限/视图】TC13 管理部门专员：权限档案与全局治理视图切换

- 用户：zhouba（管理部门专员）

| 步骤 | 操作人 | 操作 | 接口/操作 | 预期结果 |
|------|--------|------|-----------|----------|
| 1 | zhouba | 登录后查权限档案 | `GET /api/v1/auth/me` | `data.user.username=zhouba`；`data.roles` 含 `specialist`、`mgmt_specialist`；`data.permissions` 含 `global_governance_view` |
| 2 | zhouba | 登录进入系统 | 前端登录 | 默认进入**个人工作台**（我的项目/我的任务/我的运维…），顶部显示「个人工作台/全局治理」切换下拉框 |
| 3 | zhouba | 切换到全局治理 | 点击顶部切换按钮 → 全局治理 | 显示全局治理菜单（台账管理/可信仓库/风险管理/治理中心/决策与洞察中心），跳转 `/admin/asset/software-library` |
| 4 | zhouba | 查全局缺陷台账 | 全局治理 → 风险管理 → 缺陷台账 | 可见**全部部门**的缺陷记录（与 TC11 步骤 3 一致） |
| 5 | zhouba | 刷新页面 | 浏览器刷新 | 仍停留在全局治理视图（记住上次视图） |
| 6 | zhouba | 切回个人工作台 | 点击顶部切换按钮 → 个人工作台 | 恢复专员菜单，跳转「我的任务」 |

**验证点**：前端可见资源完全由 `/auth/me` 返回的权限点驱动；非 admin 角色也可具备全局治理入口。

---

### 【权限/视图】TC14 模拟外部权限系统单独授权（改库开关）

- 用户：liuyang（普通 operator，平台研发部维护人员）
- 目的：验证 `users.can_toggle_admin_view` 作为「外部权限系统单独授权」的本地模拟开关

| 步骤 | 操作人 | 操作 | 接口/操作 | 预期结果 |
|------|--------|------|-----------|----------|
| 1 | liuyang | 登录 | 前端登录 | 无切换按钮，仅个人工作台 |
| 2 | DBA | 开启授权 | `UPDATE users SET can_toggle_admin_view = 1 WHERE username = 'liuyang';` | 执行成功 |
| 3 | liuyang | 重新登录 | 退出后重新登录 | `GET /auth/me` 的 `permissions` 含 `global_governance_view`；顶部出现切换按钮 |
| 4 | liuyang | 切换全局治理 | 点击切换按钮 | 正常进入全局治理视图（菜单与周八一致） |
| 5 | DBA | 收回授权 | `UPDATE users SET can_toggle_admin_view = 0 WHERE username = 'liuyang';` | 执行成功 |
| 6 | liuyang | 重新登录 | 退出后重新登录 | 切换按钮消失，恢复个人工作台 |

**验证点**：授权数据（库里的开关）→ 权限档案 → 前端渲染，整条链路可配置、可回收。

---

### 【权限/视图】TC15 菜单按角色并集下发（接口无角色参数）

- 用户：lisi（技术部专员，仅 specialist 一个角色）、zhouba（specialist + 管理部门专员，两个菜单角色）
- 背景：`GET /api/v1/auth/menus` **无参数**，后端按用户全部角色一次返回菜单并集 `{ roles: [...], menus: { <role>: [树] } }`；不具备的角色菜单不会出现（无越权拉取入口）。

| 步骤 | 操作人 | 操作 | 预期结果 |
|------|--------|------|----------|
| 1 | lisi | 查权限档案 | `GET /auth/me`：`permissions` 为空数组（无 `global_governance_view`） |
| 2 | lisi | 拉取菜单 | `GET /api/v1/auth/menus`（无参数） | `roles: ["specialist"]`，`menus` 仅含 specialist 一套菜单树，**无 admin 菜单** |
| 3 | zhouba | 拉取菜单 | `GET /api/v1/auth/menus`（无参数） | `roles: ["specialist","admin"]`，`menus` 含两套菜单树（并集，登录后全部可见） |
| 4 | zhouba | 前端登录后切换视图 | 个人工作台 ↔ 全局治理 | 两个视图的菜单分别取自并集中的对应角色菜单，路由均已注册可直达 |

**验证点**：菜单资源由后端按"用户全部角色"统一下发（并集），接口不暴露任何角色参数，前端不做权限判断只按角色分组渲染。

---

### 【权限/视图】TC16 一个部门多名开源治理专员（任一在岗专员可审批）

- 用户：zhangsan（技术部发起人）、zhengjie（技术部第二专员）、wujiu（管理部门第二专员）
- 背景：`dept_specialist_mapping` 唯一键已调整为 `(dept, specialist_role, specialist_id)`，同一部门同一角色可配置多名专员；节点展示审批人取最早配置的在岗专员，**实际任一同部门同角色在岗专员均可处理**；后端任务列表对审批中任务下发 `canApprove` 标记（当前用户是否可处理当前节点）。

| 步骤 | 操作人 | 操作 | 预期结果 |
|------|--------|------|----------|
| 1 | zhangsan | 发起引入申请（技术部）并提交 | node_1 审批人为李四（技术部最早配置的专员），郑杰**不参与节点展示** |
| 2 | zhengjie | 登录，进入「我的任务 → 待我审批」 | 能看到该任务（`canApprove=true`），显示「审批」按钮 |
| 3 | zhengjie | 审批通过 node_1 | 通过成功；任务流转到下一节点；审批历史记录实际处理人为**郑杰**，节点记录更新为郑杰 |
| 4 | wujiu | 登录，进入「待我审批」 | 能看到任务（node_3 为当前节点，`canApprove=true`） |
| 5 | wujiu | 审批通过 node_3 | 任务终审通过；`GET /auth/me` 权限档案含 `mgmt_specialist` + `global_governance_view`（吴九同周八） |
| 6 | zhangsan | 再发起一单，由李四（lisi）审批 node_1 | 原指定审批人流程不变（回归确认） |
| 7 | zhangsan | 以郑杰身份重复验证 node_1 跳过规则：郑杰本人发起申请 | node_1 自动跳过（发起人是部门任一在岗专员即视为重复） |

**验证点**：多专员并存不破坏原单人流程；审批历史如实记录实际处理人；发起人为任一在岗专员时 node_1 防重复跳过。

---

## 八、测试数据完整性自查清单

- [x] 部门-开源治理专员映射：4 个部门 dept_specialist + usage_dept_specialist + 1 个 mgmt_dept_specialist，与 CAS 账号一致。
- [x] 开源软件使用/运维部门：5 个软件覆盖「同部门（redis/node.js）」「跨部门（kafka/spring-boot/nginx）」。
- [x] 缺陷情报源记录：6 条，覆盖 vuln-intel / sca-scan / user-report，缺陷/漏洞两条路线，含 CVE 编号；含组件缺陷 DF-TEST-006（fastjson）。
- [x] 自研项目与组件依赖：订单中心 ORD-CENTER（技术部，主岗张三）、支付平台 PAY-PLATFORM（平台研发部，主岗刘洋）；fastjson 组件被两项目依赖（component_snapshot → project_inner → project_owner 匹配链路）。
- [x] 项目维护人员：5 个软件均配置维护人，覆盖技术部/平台研发部/运维部/安全部。
- [x] 新增普通员工（维护人员）账号：liuyang / chenming / qianqi，用于二次分配的实际处理人。
- [x] 一个部门多名专员：`dept_specialist_mapping` 唯一键调整为 `(dept, specialist_role, specialist_id)`（含存量库自动迁移）；任一同部门同角色在岗专员可处理对应审批节点；任务列表下发 `canApprove` 标记。测试数据：技术部 2 名（lisi/zhengjie）、管理部门 2 名（zhouba/wujiu）。
- [x] 权限档案链路：`/auth/me` 返回 `user/roles/permissions`；管理部门专员（zhouba/wujiu）具备 `global_governance_view`；`can_toggle_admin_view` 可模拟外部权限系统单独授权。
- [x] 视图切换：具备全局治理权限的用户显示切换入口；admin 默认全局治理视图，管理部门专员默认个人工作台；上次视图刷新后保持。
- [ ] 待后续：对接外部权限系统时，替换 `internal/service/permission` 的 LocalService 实现（HTTP 调用），接口契约与前端无需变更。
- [ ] 待后续：缺陷情报源接入外部采集接口后，将 `defect_ledger.source` 替换为真实情报数据。
