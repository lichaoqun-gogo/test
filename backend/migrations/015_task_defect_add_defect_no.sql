-- 015 风险上报任务扩展表增加内部缺陷编号（系统自动生成 DF-xxx，不允许用户手填）
-- 说明：与缺陷台账内部缺陷编号风格一致（DF-纳秒时间戳），后续可打通缺陷台账
ALTER TABLE `task_defect` ADD COLUMN `internal_defect_no` varchar(64) DEFAULT NULL COMMENT '内部缺陷编号（系统自动生成 DF-xxx）' AFTER `risk_description`;
