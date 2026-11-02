-- =====================================================
-- 审批流系统数据库表结构
-- 创建时间：2026-07-05
-- 说明：实现万能审批流模板、动态路由、状态机等功能
-- =====================================================

-- 1. 审批流模板表
-- 定义不同类型的审批流程模板（引入、升级、退出、使用、风险上报）
CREATE TABLE IF NOT EXISTS `approval_template` (
    `id` BIGINT /*T![clustered_index] CLUSTERED */,
    `template_key` VARCHAR(50) NOT NULL COMMENT '模板唯一标识：introduce, upgrade, withdraw, usage, risk_report',
    `template_name` VARCHAR(100) NOT NULL COMMENT '模板名称，如：引入申请审批流',
    `description` VARCHAR(500) COMMENT '模板描述',
    `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用：1-启用，0-禁用',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_template_key` (`template_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批流模板表';

-- 2. 审批流模板节点定义表
-- 定义每个模板包含的审批节点及其条件规则
CREATE TABLE IF NOT EXISTS `approval_template_node` (
    `id` BIGINT /*T![clustered_index] CLUSTERED */,
    `template_id` BIGINT NOT NULL COMMENT '关联的模板ID',
    `node_key` VARCHAR(50) NOT NULL COMMENT '节点唯一标识：node_0, node_1, node_2, node_3',
    `node_name` VARCHAR(100) NOT NULL COMMENT '节点名称',
    `node_order` INT NOT NULL COMMENT '节点顺序（从0开始）',
    `node_type` VARCHAR(20) NOT NULL DEFAULT 'fixed' COMMENT '节点类型：fixed-固定节点，conditional-条件节点',
    `condition_type` VARCHAR(50) COMMENT '条件类型：dept_specialist-本部门专员，usage_dept-使用部门，mgmt_dept-管理部门',
    `is_skipable` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否可跳过：1-可跳过，0-不可跳过',
    `skip_rule` VARCHAR(200) COMMENT '跳过规则JSON，如：{"type":"same_person","field":"approver_id"}',
    `approver_role` VARCHAR(50) NOT NULL COMMENT '审批人角色：initiator-发起人，dept_specialist-部门专员，usage_dept_specialist-使用部门专员，mgmt_dept_specialist-管理部门专员',
    `description` VARCHAR(500) COMMENT '节点描述',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='审批流模板节点定义表';

-- 3. 扩展 approval_node 表（增加动态路由相关字段）
-- 注意：MySQL不支持IF NOT EXISTS，如果字段已存在会报错，但executeSQLFile会忽略Duplicate column错误
ALTER TABLE `approval_node`
    ADD COLUMN `node_key` VARCHAR(50) COMMENT '节点标识（来自模板）' AFTER `task_id`,
    ADD COLUMN `node_type` VARCHAR(20) COMMENT '节点类型：fixed, conditional' AFTER `node_order`,
    ADD COLUMN `approver_id` VARCHAR(50) COMMENT '审批人ID' AFTER `node_type`,
    ADD COLUMN `approver_dept` VARCHAR(100) COMMENT '审批人部门' AFTER `approver`,
    ADD COLUMN `is_skipped` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否被跳过' AFTER `approved_at`,
    ADD COLUMN `skip_reason` VARCHAR(200) COMMENT '跳过原因' AFTER `is_skipped`,
    ADD COLUMN `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER `create_time`;

-- 4. 扩展 approval_history 表（增加更多审计字段）
ALTER TABLE `approval_history`
    ADD COLUMN `node_key` VARCHAR(50) COMMENT '节点标识' AFTER `node_id`,
    ADD COLUMN `approver_id` VARCHAR(50) COMMENT '操作人ID' AFTER `approver`,
    ADD COLUMN `approver_dept` VARCHAR(100) COMMENT '操作人部门' AFTER `approver`,
    ADD COLUMN `from_status` VARCHAR(20) COMMENT '操作前状态' AFTER `opinion`,
    ADD COLUMN `to_status` VARCHAR(20) COMMENT '操作后状态' AFTER `from_status`,
    MODIFY COLUMN `opinion` VARCHAR(1000) COMMENT '审批意见';

-- 5. 扩展 governance_task 表（增加业务类型和审批流相关字段）
-- 注意：这些字段已在006迁移中添加，此处保留用于兼容旧系统
-- ALTER TABLE `governance_task`
--     ADD COLUMN `action_type` VARCHAR(50) COMMENT '业务类型：introduce-引入，upgrade-升级，withdraw-退出，usage-使用，risk_report-风险上报' AFTER `type`,
--     ADD COLUMN `applicant_id` VARCHAR(50) COMMENT '申请人ID' AFTER `applicant`,
--     ADD COLUMN `current_node_id` BIGINT COMMENT '当前审批节点ID' AFTER `current_step`,
--     ADD COLUMN `project_id` BIGINT COMMENT '关联的项目ID（开源软件）' AFTER `approval_data`,
--     ADD COLUMN `project_name` VARCHAR(200) COMMENT '关联项目名称' AFTER `project_id`,
--     ADD COLUMN `project_dept` VARCHAR(100) COMMENT '项目所属部门' AFTER `project_name`,
--     MODIFY COLUMN `task_no` VARCHAR(50) NOT NULL COMMENT '任务编号',
--     MODIFY COLUMN `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '任务状态：draft-草稿，approving-审批中，withdrawn-已撤回，approved-已通过，rejected-已驳回';

-- =====================================================
-- 初始化审批流模板数据
-- =====================================================

-- 插入万能审批流模板
INSERT INTO `approval_template` (`id`, `template_key`, `template_name`, `description`, `is_active`)
VALUES 
    (1, 'universal', '万能审批流模板', '包含4个逻辑节点的万能审批流，支持动态路由', 1)
ON DUPLICATE KEY UPDATE `template_name` = VALUES(`template_name`), `description` = VALUES(`description`);

-- 插入万能模板的节点定义
INSERT INTO `approval_template_node` (`template_id`, `node_key`, `node_name`, `node_order`, `node_type`, `condition_type`, `is_skipable`, `skip_rule`, `approver_role`, `description`)
VALUES 
    (1, 'node_0', '发起节点', 0, 'fixed', NULL, 0, NULL, 'initiator', '任务发起人提交申请'),
    (1, 'node_1', '本部门开源治理专员复审', 1, 'conditional', 'dept_specialist', 1, '{"type":"same_person","field":"approver_id"}', 'dept_specialist', '发起人所在部门的开源治理专员进行复审'),
    (1, 'node_2', '使用部门确认', 2, 'conditional', 'usage_dept', 1, '{"type":"same_dept","field":"project_dept"}', 'usage_dept_specialist', '使用部门的开源治理专员确认（仅升级/退出申请）'),
    (1, 'node_3', '管理部门评审', 3, 'fixed', 'mgmt_dept', 0, NULL, 'mgmt_dept_specialist', '开源治理办公室进行最终评审')
ON DUPLICATE KEY UPDATE `node_name` = VALUES(`node_name`);
