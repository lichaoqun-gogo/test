-- =====================================================
-- 开源治理平台（OSM）初始化数据
-- 说明：包含菜单资源、角色菜单关联、用户、部门专员映射、审批流模板
-- 使用方式：建库并执行 schema.sql 后，执行本文件初始化基础数据
-- 注意：CAS 模式下普通用户由首次登录自动创建，此处仅预设管理员
-- =====================================================

-- =====================================================
-- 1. 菜单资源（menus 表）
-- 菜单分三套：admin（全局治理）、specialist（专员工作台）、operator（运维者）
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `sort`, `menu_type`, `hidden`, `group`) VALUES
-- ===== admin 角色：全局治理菜单（id 1-23）=====
(1,  0, '台账管理',       '',                                   '',                                       'Coin',         1, 1, 0, 'asset'),
(2,  1, '开源基础软件',   '/admin/asset/software-library',      'views/admin/asset/SoftwareLibrary',      '',             1, 2, 0, ''),
(3,  1, '开源依赖包',     '/admin/asset/dependency-library',    'views/admin/asset/DependencyLibrary',    '',             2, 2, 0, ''),
(4,  1, '开源工具',       '/admin/asset/tool-library',           'views/admin/asset/ToolLibrary',          '',             3, 2, 0, ''),
(5,  1, '自研项目',       '/admin/asset/self-developed',         'views/admin/asset/SelfDevelopedProjects','',             4, 2, 0, ''),
(6,  1, '源码备份',       '/admin/asset/source-backup',          'views/admin/asset/SourceBackup',         '',             5, 2, 0, ''),
(7,  0, '可信仓库',       '/admin/trusted-repo',                 'views/admin/asset/TrustedRepository',    'Files',        2, 2, 0, ''),
(8,  0, '风险管理',       '',                                    '',                                       'Warning',      3, 1, 0, 'risk'),
(9,  8, '名单管理',       '/admin/risk/blacklist',               'views/admin/risk/Blacklist',             '',             1, 2, 0, ''),
(10, 8, 'CICD检查',       '/admin/risk/cicd-check',              'views/admin/risk/CicdCheck',             '',             2, 2, 0, ''),
(11, 8, '许可证信息库',   '/admin/risk/license-library',         'views/admin/risk/LicenseLibrary',        '',             3, 2, 0, ''),
(12, 8, '缺陷台账',       '/admin/risk/defect-ledger',           'views/admin/risk/DefectLedger',          '',             4, 2, 0, ''),
(13, 8, '活跃度健康度监测','/admin/risk/activity-monitor',       'views/admin/risk/ActivityMonitor',       '',             5, 2, 0, ''),
(14, 8, 'SCA自动化解析入库','/admin/risk/sca-parser',            'views/admin/risk/ScaParser',             '',             6, 2, 0, ''),
(15, 0, '治理中心',       '',                                    '',                                       'Setting',      4, 1, 0, 'governance'),
(16,15, '任务列表',       '/admin/governance/task-list',         'views/admin/governance/TaskList',        '',             1, 2, 0, ''),
(17, 0, '决策与洞察中心', '',                                    '',                                       'DataAnalysis', 5, 1, 0, 'decision'),
(18,17, '治理效能看板',   '',                                    '',                                       '',             1, 1, 0, ''),
(19,18, '违规拦截统计',   '/admin/decision/interception',        'views/admin/decision/interception',      '',             1, 2, 0, ''),
(20,18, '缺陷修复效能',   '/admin/decision/repair-efficiency',   'views/admin/decision/repair-efficiency', '',             2, 2, 0, ''),
(21,18, '资产健康度分布', '/admin/decision/health-distribution',  'views/admin/decision/health-distribution','',           3, 2, 0, ''),
(22,18, '治理工作量统计', '/admin/decision/workload',             'views/admin/decision/workload',          '',             4, 2, 0, ''),
(23,17, '供应链安全洞察', '/admin/decision/core-dependencies',    'views/admin/decision/core-dependencies', '',             2, 2, 0, ''),

-- ===== specialist 角色：专员工作台菜单（id 24-30）=====
(24, 0, '我的项目', '/specialist/my-projects',  'views/operator/my-assets/MyProjects', 'Monitor', 1, 2, 0, ''),
(25, 0, '我的任务', '/specialist/my-task',      'views/specialist/MyTask',             'List',    2, 2, 0, ''),
(26, 0, '我的运维', '',                         '',                                    'Setting', 3, 1, 0, ''),
(27,26, '运维台账', '/specialist/my-ops',       'views/specialist/MyOps',              '',        1, 2, 0, ''),
(28,26, '人员管理', '/specialist/personnel-mgmt','views/operator/PersonnelMgmt',       '',        2, 2, 0, ''),
(29, 0, '我的使用', '/specialist/my-usage',     'views/operator/MyUsage',              'Files',   4, 2, 0, ''),
(30, 0, '可信仓库', '/specialist/trusted-repo', 'views/operator/TrustedRepo',          'Files',   5, 2, 0, ''),

-- ===== operator 角色：运维者菜单（id 31-35）=====
(31, 0, '我的项目', '/operator/my-projects',  'views/operator/my-assets/MyProjects', 'Monitor', 1, 2, 0, ''),
(32, 0, '我的任务', '/operator/my-task',      'views/operator/MyTask',               'List',    2, 2, 0, ''),
(33, 0, '我的运维', '/operator/my-ops',       'views/operator/MyOps',                'Setting', 3, 2, 0, ''),
(34, 0, '我的使用', '/operator/my-usage',     'views/operator/MyUsage',              'Files',   4, 2, 0, ''),
(35, 0, '可信仓库', '/operator/trusted-repo', 'views/operator/TrustedRepo',          'Files',   5, 2, 0, '');

-- =====================================================
-- 2. 角色-菜单关联（role_menus 表）
-- admin 角色对应菜单 id 1-23
-- specialist 角色对应菜单 id 24-30
-- operator 角色对应菜单 id 31-35
-- =====================================================
INSERT INTO `role_menus` (`role`, `menu_id`) VALUES
-- admin
('admin', 1),  ('admin', 2),  ('admin', 3),  ('admin', 4),  ('admin', 5),
('admin', 6),  ('admin', 7),  ('admin', 8),  ('admin', 9),  ('admin', 10),
('admin', 11), ('admin', 12), ('admin', 13), ('admin', 14), ('admin', 15),
('admin', 16), ('admin', 17), ('admin', 18), ('admin', 19), ('admin', 20),
('admin', 21), ('admin', 22), ('admin', 23),
-- specialist
('specialist', 24), ('specialist', 25), ('specialist', 26),
('specialist', 27), ('specialist', 28), ('specialist', 29), ('specialist', 30),
-- operator
('operator', 31), ('operator', 32), ('operator', 33),
('operator', 34), ('operator', 35);

-- =====================================================
-- 3. 用户数据（users 表）
-- CAS 模式下普通用户由首次登录自动创建（默认 operator 角色）
-- 此处仅预设管理员账户，上线后请修改 role 和 department 为实际值
-- =====================================================
INSERT INTO `users` (`id`, `username`, `nickname`, `email`, `phone`, `department`, `role`, `can_toggle_admin_view`, `status`, `create_at`, `update_at`) VALUES
(2, 'admin', '系统管理员', 'admin@qim.local', '13800000000', '管理层', 'admin', 0, 1, NOW(), NOW());

-- =====================================================
-- 4. 部门开源治理专员映射（dept_specialist_mapping 表）
-- 每个部门可配置 3 种专员角色：
--   dept_specialist        - 部门开源治理专员
--   usage_dept_specialist  - 使用部门专员
--   mgmt_dept_specialist   - 管理部门专员（开源治理办公室）
-- 上线后请替换为实际部门和人员
-- =====================================================
INSERT INTO `dept_specialist_mapping` (`dept`, `specialist_id`, `specialist_name`, `specialist_role`, `is_active`) VALUES
('开源治理办公室', 'admin', '系统管理员', 'mgmt_dept_specialist', 1);

-- =====================================================
-- 5. 审批流模板（approval_template 表）
-- 万能审批流模板：4 节点动态路由（发起→部门专员复审→使用部门确认→管理部门评审）
-- =====================================================
INSERT INTO `approval_template` (`id`, `template_key`, `template_name`, `description`, `is_active`) VALUES
(1, 'universal', '万能审批流模板', '包含4个逻辑节点的万能审批流，支持动态路由', 1);

-- =====================================================
-- 6. 审批流模板节点定义（approval_template_node 表）
-- node_0: 发起节点（固定）
-- node_1: 本部门专员复审（条件节点，可跳过-同一人时跳过）
-- node_2: 使用部门确认（条件节点，可跳过-同部门时跳过）
-- node_3: 管理部门评审（固定，不可跳过）
-- =====================================================
INSERT INTO `approval_template_node` (`template_id`, `node_key`, `node_name`, `node_order`, `node_type`, `condition_type`, `is_skipable`, `skip_rule`, `approver_role`, `description`) VALUES
(1, 'node_0', '发起节点',             0, 'fixed',       NULL,            0, NULL,                                  'initiator',             '任务发起人提交申请'),
(1, 'node_1', '本部门开源治理专员复审', 1, 'conditional', 'dept_specialist', 1, '{"type":"same_person","field":"approver_id"}', 'dept_specialist',       '发起人所在部门的开源治理专员进行复审'),
(1, 'node_2', '使用部门确认',           2, 'conditional', 'usage_dept',     1, '{"type":"same_dept","field":"project_dept"}', 'usage_dept_specialist', '使用部门的开源治理专员确认（仅升级/退出申请）'),
(1, 'node_3', '管理部门评审',           3, 'fixed',       'mgmt_dept',      0, NULL,                                  'mgmt_dept_specialist',  '开源治理办公室进行最终评审');
