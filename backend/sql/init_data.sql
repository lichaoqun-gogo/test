-- 开源治理平台初始化数据
-- 执行前请确保已运行 schema.sql 创建表结构

-- 初始化审批模板
INSERT INTO `approval_template` (`id`, `template_key`, `template_name`, `description`, `is_active`, `create_time`, `update_time`)
VALUES
  (1, 'introduce', '引入申请审批流程', NULL, 1, NOW(), NOW()),
  (2, 'upgrade', '升级申请审批流程', NULL, 1, NOW(), NOW()),
  (3, 'withdraw', '退出申请审批流程', NULL, 1, NOW(), NOW()),
  (4, 'usage', '使用申请审批流程', NULL, 1, NOW(), NOW()),
  (5, 'risk_report', '风险上报审批流程', NULL, 1, NOW(), NOW());

-- 初始化模板节点
INSERT INTO `approval_template_node` (`id`, `template_id`, `node_key`, `node_name`, `node_order`, `node_type`, `approver_role`, `is_skipable`, `create_time`, `update_time`)
VALUES
  (1, 1, 'node_0', '使用部门专员审核', 0, 'fixed', 'usage_dept_specialist', 0, NOW(), NOW()),
  (2, 1, 'node_1', '管理部门专员审核', 1, 'fixed', 'mgmt_dept_specialist', 0, NOW(), NOW()),
  (3, 2, 'node_0', '使用部门专员审核', 0, 'fixed', 'usage_dept_specialist', 0, NOW(), NOW()),
  (4, 2, 'node_1', '管理部门专员审核', 1, 'fixed', 'mgmt_dept_specialist', 0, NOW(), NOW()),
  (5, 3, 'node_0', '使用部门专员审核', 0, 'fixed', 'usage_dept_specialist', 0, NOW(), NOW()),
  (6, 3, 'node_1', '管理部门专员审核', 1, 'fixed', 'mgmt_dept_specialist', 0, NOW(), NOW()),
  (7, 4, 'node_0', '使用部门专员审核', 0, 'fixed', 'usage_dept_specialist', 0, NOW(), NOW()),
  (8, 5, 'node_0', '使用部门专员审核', 0, 'fixed', 'usage_dept_specialist', 0, NOW(), NOW()),
  (9, 5, 'node_1', '管理部门专员审核', 1, 'fixed', 'mgmt_dept_specialist', 0, NOW(), NOW());

-- menus / role_menus 种子已废弃：菜单统一由「统一权限系统」按角色下发或前端本地映射（见 PLATFORM_DESIGN.md 第 11 章）

-- 初始化管理员用户（CAS登录后会自动创建/更新用户信息）
-- 注：实际部署时建议删除此记录，改为CAS认证自动同步
INSERT INTO `users` (`id`, `username`, `nickname`, `email`, `phone`, `department`, `team`, `role`, `can_toggle_admin_view`, `status`, `create_at`, `update_at`)
VALUES
  (1, 'admin', '系统管理员', 'admin@company.com', '13800000000', '技术部-管理组', '管理组', 'admin', 1, 1, NOW(), NOW()),
  (2, 'lisi', '李四', 'lisi@company.com', '13800000001', '技术部', '前端组', 'specialist', 0, 1, NOW(), NOW()),
  (3, 'zhangsan', '张三', 'zhangsan@company.com', '13800000002', '技术部-前端组', '前端组', 'operator', 0, 1, NOW(), NOW());

-- 初始化许可证分组
INSERT INTO `licensegroup` (`id`, `_name`, `RISKWEIGHT`, `UUID`, `create_time`, `update_time`)
VALUES
  (1, '宽松许可证', 0, UUID(), NOW(), NOW()),
  (2, '弱保护许可证', 1, UUID(), NOW(), NOW()),
  (3, '强保护许可证', 3, UUID(), NOW(), NOW());

-- 初始化许可证数据（示例）
INSERT INTO `license` (`id`, `LICENSE`, `_name`, `IS_OSIAPPROVED`, `IS_DEPRECATED`, `UUID`, `create_time`, `update_time`)
VALUES
  (1, 'MIT', 'MIT License', 1, 0, UUID(), NOW(), NOW()),
  (2, 'Apache-2.0', 'Apache License 2.0', 1, 0, UUID(), NOW(), NOW()),
  (3, 'BSD-3-Clause', 'BSD 3-Clause', 1, 0, UUID(), NOW(), NOW()),
  (4, 'GPL-3.0', 'GNU General Public License v3.0', 1, 0, UUID(), NOW(), NOW()),
  (5, 'LGPL-2.1', 'GNU Lesser General Public License v2.1', 1, 0, UUID(), NOW(), NOW());

-- 许可证分组关联
INSERT INTO `licensegroup_license` (`licensegroup_id`, `license_id`, `create_time`, `update_time`)
VALUES
  (1, 1, NOW(), NOW()), (1, 2, NOW(), NOW()), (1, 3, NOW(), NOW()),
  (2, 5, NOW(), NOW()),
  (3, 4, NOW(), NOW());
