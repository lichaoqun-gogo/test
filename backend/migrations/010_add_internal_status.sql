-- 010 缺陷实体表增加内部处置状态 internal_status（幂等）
-- 设计调整：内部处置状态以实体表为准（台账直接查两张实体表即可展示状态），defect_flow 退化为流程明细表
-- 注意：状态回填与 defect_flow.status 删列在 Go 侧 migrateDefectInternalStatus 完成（information_schema 守卫）
ALTER TABLE `defect` ADD COLUMN `internal_status` varchar(20) NOT NULL DEFAULT 'pending_analysis' COMMENT '内部处置状态: pending_analysis-待研判, pending_review-待上会, assigned-已分配, ignored-已忽略, closed-已闭环' AFTER `issue_state`;
ALTER TABLE `defect` ADD KEY `idx_internal_status` (`internal_status`);
ALTER TABLE `defect_ledger` ADD COLUMN `internal_status` varchar(20) NOT NULL DEFAULT 'pending_analysis' COMMENT '内部处置状态: pending_analysis-待研判, pending_review-待上会, assigned-已分配, ignored-已忽略, closed-已闭环' AFTER `discovered_date`;
ALTER TABLE `defect_ledger` ADD KEY `idx_internal_status` (`internal_status`);
