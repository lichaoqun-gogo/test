-- 009 缺陷实体/处置流程分离（增量建表部分，幂等）
-- 注意：数据回填与删列在 Go 侧 migrateDefectEntityFlow 完成（information_schema 守卫，避免连接池会话变量问题）
-- 目标结构：
--   defect（外部实体，已建）/ defect_ledger（内部实体，纯化）/ defect_flow（缺陷级处置流程，新）
--   defect_disposal / defect_status_history 改为按 defect_no（缺陷编号）关联，不再依赖实体表主键
CREATE TABLE IF NOT EXISTS `defect_flow` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `defect_no` varchar(100) NOT NULL COMMENT '缺陷编号（外部: redis#13401 / 内部: DF-xxx），一缺陷一流程',
  `flow_type` varchar(20) DEFAULT 'defect' COMMENT '流程类型: defect-缺陷, vuln-漏洞',
  `maint_dept` varchar(128) DEFAULT NULL COMMENT '运维部门（权限过滤与裂变依据）',
  `deadline` datetime(3) DEFAULT NULL COMMENT '上会/研判确定的处置截止时间',
  `status` varchar(20) NOT NULL DEFAULT 'pending_analysis' COMMENT '状态: pending_analysis-待研判, pending_review-待上会, assigned-已分配, ignored-已忽略, closed-已闭环',
  `analyze_result` varchar(500) DEFAULT NULL COMMENT '专员研判结果',
  `review_meeting_no` varchar(50) DEFAULT NULL COMMENT '评审会议编号',
  `review_conclusion` varchar(20) DEFAULT NULL COMMENT '评审结论: no_action-无需处置, need_fix-需要修复',
  `dispatch_type` varchar(20) DEFAULT NULL COMMENT '分发类型: analyze-研判分配等',
  `dispatch_target` varchar(50) DEFAULT NULL COMMENT '分发目标人员（研判人）',
  `dispatch_dept` varchar(100) DEFAULT NULL COMMENT '分发目标部门',
  `dispatch_time` datetime(3) DEFAULT NULL COMMENT '分发时间',
  `dispatcher` varchar(20) DEFAULT NULL COMMENT '分发人',
  `analyze_by` varchar(50) DEFAULT NULL COMMENT '研判人',
  `analyze_time` datetime(3) DEFAULT NULL COMMENT '研判时间',
  `review_by` varchar(50) DEFAULT NULL COMMENT '评审人',
  `review_time` datetime(3) DEFAULT NULL COMMENT '评审时间',
  `close_time` datetime(3) DEFAULT NULL COMMENT '闭环时间',
  `create_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `update_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_defect_no` (`defect_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='缺陷处置流程表（缺陷级状态机，按缺陷编号关联实体）';

ALTER TABLE `defect_disposal` ADD COLUMN `defect_no` varchar(100) NOT NULL DEFAULT '' COMMENT '缺陷编号（按编号关联缺陷实体）' AFTER `id`;
ALTER TABLE `defect_disposal` ADD KEY `idx_defect_no` (`defect_no`);
ALTER TABLE `defect_status_history` ADD COLUMN `defect_no` varchar(100) NOT NULL DEFAULT '' COMMENT '缺陷编号（按编号关联缺陷实体）' AFTER `id`;
ALTER TABLE `defect_status_history` ADD KEY `idx_defect_no` (`defect_no`);
