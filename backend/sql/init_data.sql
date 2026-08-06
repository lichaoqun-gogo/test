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

-- 初始化菜单
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `sort`, `menu_type`, `hidden`, `create_at`)
VALUES
  (1, 0, '我的项目', '/operator/my-projects', 'views/operator/my-assets/MyProjects', 'Monitor', 0, 2, 0, NOW()),
  (2, 0, '我的任务', '/operator/my-tasks', 'views/operator/MyTasks', 'List', 1, 2, 0, NOW()),
  (3, 0, '我的运维', '', '', 'Setting', 2, 1, 0, NOW()),
  (4, 3, '运维台账', '/operator/my-ops', 'views/operator/MyOps', '', 0, 2, 0, NOW()),
  (5, 3, '人员管理', '/operator/personnel-mgmt', 'views/operator/PersonnelMgmt', '', 1, 2, 0, NOW()),
  (6, 0, '我的使用', '/operator/my-usage', 'views/operator/MyUsage', 'Files', 3, 2, 0, NOW()),
  (7, 0, '可信仓库', '/operator/trusted-repo', 'views/operator/TrustedRepo', 'Files', 4, 2, 0, NOW()),
  (8, 0, '软件选型服务', '/operator/software-service', 'views/operator/software-service/SoftwareService', '', 5, 2, 0, NOW()),
  (9, 0, '我的申请', '/operator/my-applications', 'views/operator/my-applications/MyApplications', '', 6, 2, 0, NOW()),
  (10, 0, '资产中心', '', '', 'Setting', 0, 1, 0, NOW()),
  (11, 10, '开源基础软件', '/admin/asset/software-library', 'views/admin/asset/SoftwareLibrary', '', 0, 2, 0, NOW()),
  (12, 10, '开源工具库', '/admin/asset/tool-library', 'views/admin/asset/ToolLibrary', '', 1, 2, 0, NOW()),
  (13, 10, '依赖库', '/admin/asset/dependency-library', 'views/admin/asset/DependencyLibrary', '', 2, 2, 0, NOW()),
  (14, 10, '许可证管理', '/admin/asset/license-management', 'views/admin/asset/LicenseManagement', '', 3, 2, 0, NOW()),
  (15, 10, '活跃度监测', '/admin/asset/activity-monitor', 'views/admin/asset/ActivityMonitor', '', 4, 2, 0, NOW()),
  (16, 10, '源码备份', '/admin/asset/source-backup', 'views/admin/asset/SourceBackup', '', 5, 2, 0, NOW()),
  (17, 10, '可信仓库', '/admin/asset/trusted-repo', 'views/admin/asset/TrustedRepo', '', 6, 2, 0, NOW()),
  (18, 0, '风险中心', '', '', 'Warning', 1, 1, 0, NOW()),
  (19, 18, '统一风险台账', '/admin/risk/risk-ledger', 'views/admin/risk/RiskLedger', '', 0, 2, 0, NOW()),
  (20, 18, '缺陷台账', '/admin/risk/defect-ledger', 'views/admin/risk/DefectLedger', '', 1, 2, 0, NOW()),
  (21, 18, '黑白名单管理', '/admin/risk/blacklist-management', 'views/admin/risk/BlacklistManagement', '', 2, 2, 0, NOW()),
  (22, 18, '风险任务', '/admin/risk/task-list', 'views/admin/risk/TaskList', '', 3, 2, 0, NOW()),
  (23, 18, 'SCA解析记录', '/admin/risk/sca-parser', 'views/admin/risk/ScaParser', '', 4, 2, 0, NOW()),
  (24, 0, '治理中心', '', '', 'Coin', 2, 1, 0, NOW()),
  (25, 24, '引入审批', '/admin/governance/introduce-approval', 'views/admin/governance/IntroduceApproval', '', 0, 2, 0, NOW()),
  (26, 24, '升级评审', '/admin/governance/upgrade-review', 'views/admin/governance/UpgradeReview', '', 1, 2, 0, NOW()),
  (27, 24, '退出审批', '/admin/governance/withdraw-approval', 'views/admin/governance/WithdrawApproval', '', 2, 2, 0, NOW()),
  (28, 24, '存量治理', '/admin/governance/inventory-management', 'views/admin/governance/InventoryManagement', '', 3, 2, 0, NOW()),
  (29, 24, '违规审计', '/admin/governance/violation-audit', 'views/admin/governance/ViolationAudit', '', 4, 2, 0, NOW()),
  (30, 0, '决策中心', '', '', 'DataAnalysis', 3, 1, 0, NOW()),
  (31, 30, '治理看板', '/admin/decision/governance-dashboard', 'views/admin/decision/GovernanceDashboard', '', 0, 2, 0, NOW()),
  (32, 30, '安全洞察', '/admin/decision/security-insights', 'views/admin/decision/SecurityInsights', '', 1, 2, 0, NOW()),
  (33, 30, '核心依赖看板', '/admin/decision/core-dependencies', 'views/admin/decision/CoreDependencies', '', 2, 2, 0, NOW()),
  (34, 0, '个人工作台', '', '', 'SwitchButton', 0, 1, 0, NOW()),
  (35, 34, '我的项目', '/specialist/my-projects', 'views/operator/my-assets/MyProjects', 'Monitor', 0, 2, 0, NOW()),
  (36, 34, '我的任务', '/specialist/my-tasks', 'views/operator/MyTasks', 'List', 1, 2, 0, NOW()),
  (37, 34, '我的运维', '', '', 'Setting', 2, 1, 0, NOW()),
  (38, 37, '运维台账', '/specialist/my-ops', 'views/specialist/MyOps', '', 0, 2, 0, NOW()),
  (39, 37, '人员管理', '/specialist/personnel-mgmt', 'views/operator/PersonnelMgmt', '', 1, 2, 0, NOW()),
  (40, 34, '我的使用', '/specialist/my-usage', 'views/operator/MyUsage', 'Files', 3, 2, 0, NOW()),
  (41, 34, '可信仓库', '/specialist/trusted-repo', 'views/operator/TrustedRepo', 'Files', 4, 2, 0, NOW());

-- 初始化角色菜单关联
INSERT INTO `role_menus` (`id`, `role`, `menu_id`)
VALUES
  (1, 'operator', 1), (2, 'operator', 2), (3, 'operator', 3), (4, 'operator', 4), (5, 'operator', 5), (6, 'operator', 6), (7, 'operator', 7), (8, 'operator', 8), (9, 'operator', 9),
  (10, 'specialist', 35), (11, 'specialist', 36), (12, 'specialist', 37), (13, 'specialist', 38), (14, 'specialist', 39), (15, 'specialist', 40), (16, 'specialist', 41),
  (17, 'admin', 10), (18, 'admin', 11), (19, 'admin', 12), (20, 'admin', 13), (21, 'admin', 14), (22, 'admin', 15), (23, 'admin', 16), (24, 'admin', 17),
  (25, 'admin', 18), (26, 'admin', 19), (27, 'admin', 20), (28, 'admin', 21), (29, 'admin', 22), (30, 'admin', 23),
  (31, 'admin', 24), (32, 'admin', 25), (33, 'admin', 26), (34, 'admin', 27), (35, 'admin', 28), (36, 'admin', 29),
  (37, 'admin', 30), (38, 'admin', 31), (39, 'admin', 32), (40, 'admin', 33),
  (41, 'admin', 34), (42, 'admin', 35), (43, 'admin', 36), (44, 'admin', 37), (45, 'admin', 38), (46, 'admin', 39), (47, 'admin', 40), (48, 'admin', 41);

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
