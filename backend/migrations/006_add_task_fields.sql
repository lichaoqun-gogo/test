-- 为governance_task表添加审批流相关字段
-- 注意：MySQL不支持IF NOT EXISTS，如果字段已存在会报错，但executeSQLFile会忽略Duplicate column错误
ALTER TABLE `governance_task`
    ADD COLUMN `action_type` varchar(50) DEFAULT NULL COMMENT '业务类型：introduce(引入), upgrade(升级), withdraw(退出), usage(使用), risk_report(风险上报)' AFTER `type`,
    ADD COLUMN `applicant_id` varchar(50) DEFAULT NULL COMMENT '申请人ID' AFTER `applicant`,
    ADD COLUMN `current_node_id` bigint(20) DEFAULT NULL COMMENT '当前审批节点ID' AFTER `current_step`,
    ADD COLUMN `project_id` bigint(20) DEFAULT NULL COMMENT '关联的项目ID（开源软件）' AFTER `approval_data`,
    ADD COLUMN `project_name` varchar(200) DEFAULT NULL COMMENT '关联项目名称' AFTER `project_id`,
    ADD COLUMN `project_dept` varchar(100) DEFAULT NULL COMMENT '项目所属部门' AFTER `project_name`;

ALTER TABLE `governance_task`
    MODIFY COLUMN `current_step` varchar(100) NOT NULL DEFAULT '' COMMENT '当前审批步骤描述',
    MODIFY COLUMN `assigned_to` varchar(100) DEFAULT NULL COMMENT '当前处理人',
    MODIFY COLUMN `status` varchar(20) NOT NULL DEFAULT 'draft' COMMENT '状态: draft-草稿, approving-审批中, withdrawn-已撤回, approved-已通过, rejected-已驳回';
