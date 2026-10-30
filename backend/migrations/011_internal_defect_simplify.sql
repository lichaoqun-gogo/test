-- 011 内部缺陷精简为风险上报字段 / 外部缺陷增加运维部门列（幂等）
-- 内部缺陷去掉结构字段（risk_type/affected_component/...）在 Go 侧 migrateInternalDefectSimplify 完成（information_schema 守卫）
-- 外部缺陷增加 maint_dept：首次导入记录导入者部门，重复导入不覆盖（导入逻辑控制）
ALTER TABLE `defect` ADD COLUMN `maint_dept` varchar(128) NULL COMMENT '运维部门（首次导入者部门，重复导入不覆盖）' AFTER `repair_advice`;
