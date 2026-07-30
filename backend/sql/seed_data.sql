-- =====================================================
-- 开源治理平台 V3 - 测试数据脚本
-- =====================================================

-- 1. 开源基础软件数据
INSERT INTO `project` (`ACTIVE`, `AUTHOR`, `PUBLISHER`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `license`, `UUID`, `type`, `applicable_level`, `maint_dept`, `use_dept`, `country`, `create_by`) VALUES
(1, 'Redis contributors', 'Redis Ltd', 'redis', 'Redis', '7.0.11', '开源内存数据存储', 'FRAMEWORK', 'BSD-3-Clause', UUID(), 'java', 'A', '运维部', '技术研发部,平台研发部', '美国', 'admin'),
(1, 'Pivotal Software', 'VMware', 'org.springframework.boot', 'Spring Boot', '3.2.0', 'Spring应用框架', 'LIBRARY', 'Apache-2.0', UUID(), 'java', 'A', '平台研发部', '技术研发部', '美国', 'admin'),
(1, 'Oracle', 'Oracle', 'mysql', 'MySQL', '8.0.35', '关系型数据库', 'DATABASE', 'GPL-2.0', UUID(), 'java', 'B', '运维部', '技术研发部,质量管理部', '美国', 'admin'),
(1, 'Igor Sysoev', 'Nginx Inc', 'nginx', 'Nginx', '1.25.3', '高性能HTTP服务器', 'CONTAINER', 'BSD-2-Clause', UUID(), 'go', 'A', '运维部', '平台研发部', '俄罗斯', 'admin'),
(1, 'Apache', 'Apache', 'org.apache.kafka', 'Kafka', '3.6.0', '分布式消息队列', 'LIBRARY', 'Apache-2.0', UUID(), 'java', 'A', '平台研发部', '技术研发部', '美国', 'admin'),
(1, 'Apache', 'Apache', 'org.apache.zookeeper', 'Zookeeper', '3.9.1', '分布式协调服务', 'LIBRARY', 'Apache-2.0', UUID(), 'java', 'B', '平台研发部', '技术研发部', '美国', 'admin'),
(1, 'Vue.js contributors', 'Vue.js', 'vue', 'Vue.js', '3.4.0', '渐进式JavaScript框架', 'LIBRARY', 'MIT', UUID(), 'js', 'A', '交易组', '技术研发部', '中国', 'admin'),
(1, 'Node.js contributors', 'OpenJS Foundation', 'nodejs', 'Node.js', '20.10.0', 'JavaScript运行时环境', 'RUNTIME', 'MIT', UUID(), 'js', 'A', '交易组', '平台研发部', '美国', 'admin'),
(1, 'MyBatis contributors', 'MyBatis', 'org.mybatis', 'MyBatis', '3.5.15', '持久层框架', 'LIBRARY', 'Apache-2.0', UUID(), 'java', 'B', '支付组', '技术研发部', '中国', 'admin'),
(1, 'Apache', 'Apache', 'org.apache.commons', 'Commons Lang', '3.14.0', 'Java工具类库', 'LIBRARY', 'Apache-2.0', UUID(), 'java', 'A', '支付组', '技术研发部,平台研发部', '美国', 'admin');

-- 2. 自研项目数据
INSERT INTO `project_inner` (`pro_code`, `active`, `type`, `is_boundary`, `group_name`, `dept_name`, `module_code`, `branch`, `current_ID`, `create_by`) VALUES
('ORDER-SYS', 1, 'java', 1, '交易组', '技术研发部', 'order-service', 'main', 'ORD-2024-001', 'admin'),
('PAY-SYS', 1, 'java', 1, '支付组', '平台研发部', 'pay-service', 'main', 'PAY-2024-001', 'admin'),
('USER-CENTER', 1, 'java', 1, '用户组', '技术研发部', 'user-service', 'main', 'USR-2024-001', 'admin'),
('MSG-PLATFORM', 1, 'java', 0, '消息组', '平台研发部', 'msg-service', 'develop', 'MSG-2024-001', 'admin'),
('RISK-CONTROL', 1, 'java', 1, '风控组', '安全部', 'risk-service', 'main', 'RSK-2024-001', 'admin'),
('DATA-PLATFORM', 1, 'go', 0, '数据组', '技术研发部', 'data-service', 'main', 'DAT-2024-001', 'admin');

-- 3. 开源依赖包数据
INSERT INTO `component_snapshot` (`current_id`, `_group`, `_name`, `_version`, `DESCRIPTION`, `CLASSIFIER`, `INTERNAL`, `LICENSE`, `LAST_RISKSCORE`, `PROJECT_ID`) VALUES
('SRC-001', 'com.google.guava', 'guava', '32.1.3-jre', 'Google核心库', 'LIBRARY', 0, 'Apache-2.0', 0.120, 1),
('SRC-002', 'org.apache.commons', 'commons-lang3', '3.14.0', 'Apache通用工具', 'LIBRARY', 0, 'Apache-2.0', 0.080, 1),
('SRC-003', 'com.fasterxml.jackson.core', 'jackson-databind', '2.16.0', 'JSON处理库', 'LIBRARY', 0, 'Apache-2.0', 0.150, 2),
('SRC-004', 'org.springframework', 'spring-core', '6.1.1', 'Spring核心模块', 'LIBRARY', 0, 'Apache-2.0', 0.090, 2),
('SRC-005', 'io.netty', 'netty-all', '4.1.100.Final', '网络通信框架', 'LIBRARY', 0, 'Apache-2.0', 0.110, 3),
('SRC-006', 'com.alibaba', 'fastjson', '2.0.42', 'JSON处理库', 'LIBRARY', 0, 'Apache-2.0', 0.200, 3),
('SRC-007', 'org.mybatis', 'mybatis', '3.5.14', 'ORM框架', 'LIBRARY', 0, 'Apache-2.0', 0.070, 4),
('SRC-008', 'com.github.pagehelper', 'pagehelper', '6.1.0', '分页插件', 'LIBRARY', 0, 'MIT', 0.050, 4);

-- 4. 黑名单配置数据
INSERT INTO `component_core` (`_format`, `type`, `_name`, `_group`, `version_exp`, `applicable_level`, `support_team`, `reason`, `advise`, `_enable`, `alert_type`, `create_by`, `is_boundary`) VALUES
(1, 'maven', 'log4j', 'org.apache.log4j', '[0, 2.17.0)', 'A', '安全部', '存在远程代码执行漏洞(CVE-2021-44228)', '升级到2.17.0或以上版本', 1, 1, 'admin', 2),
(1, 'maven', 'fastjson', 'com.alibaba', '[0, 1.2.83)', 'A', '安全部', '存在反序列化漏洞', '升级到1.2.83或以上版本', 1, 1, 'admin', 2),
(2, 'pypi', 'requests', 'requests', '[0, 2.31.0)', 'B', '安全部', '存在SSL验证绕过漏洞', '升级到2.31.0或以上版本', 1, 2, 'admin', 2),
(3, 'npm', 'lodash', 'lodash', '[0, 4.17.21)', 'A', '安全部', '存在原型链污染漏洞', '升级到4.17.21或以上版本', 1, 1, 'admin', 2);

-- 5. 许可证数据
INSERT INTO `license` (`_COMMENT`, `IS_DEPRECATED`, `FSFLIBRE`, `LICENSE`, `_name`, `IS_OSIAPPROVED`, `UUID`, `create_by`) VALUES
('MIT许可证，宽松型', 0, 1, 'MIT', 'MIT License', 1, UUID(), 'admin'),
('Apache许可证2.0版本', 0, 1, 'Apache-2.0', 'Apache License 2.0', 1, UUID(), 'admin'),
('GNU通用公共许可证2.0', 0, 1, 'GPL-2.0', 'GNU General Public License v2.0', 1, UUID(), 'admin'),
('GNU通用公共许可证3.0', 0, 1, 'GPL-3.0', 'GNU General Public License v3.0', 1, UUID(), 'admin'),
('BSD 3条款许可证', 0, 1, 'BSD-3-Clause', 'BSD 3-Clause License', 1, UUID(), 'admin'),
('BSD 2条款许可证', 0, 1, 'BSD-2-Clause', 'BSD 2-Clause License', 1, UUID(), 'admin'),
('LGPL许可证', 0, 1, 'LGPL-3.0', 'GNU Lesser General Public License v3.0', 1, UUID(), 'admin');

-- 6. 许可证分组数据
INSERT INTO `licensegroup` (`_name`, `RISKWEIGHT`, `UUID`, `create_by`) VALUES
('宽松型许可证', 10, UUID(), 'admin'),
('强Copyleft型', 80, UUID(), 'admin'),
('弱Copyleft型', 40, UUID(), 'admin'),
('非标准许可证', 60, UUID(), 'admin');

-- 7. 许可证分组关联
INSERT INTO `licensegroup_license` (`licensegroup_id`, `license_id`, `create_by`) VALUES
(1, 1, 'admin'), (1, 2, 'admin'), (1, 5, 'admin'), (1, 6, 'admin'),
(2, 3, 'admin'), (2, 4, 'admin'),
(3, 7, 'admin');

-- 8. 缺陷台账数据
INSERT INTO `defect_ledger` (`defect_no`, `source`, `title`, `risk_level`, `risk_type`, `affected_component`, `affected_component_type`, `discovered_date`, `status`, `description`, `create_by`) VALUES
('DEF-2024-001', 'sca', 'Log4j远程代码执行漏洞', 'critical', 'vulnerability', 'org.apache.log4j:log4j', 'dependency', '2024-01-15', 'pending_analysis', 'Log4j存在严重远程代码执行漏洞(CVE-2021-44228)', 'admin'),
('DEF-2024-002', 'external', 'Fastjson反序列化漏洞', 'high', 'vulnerability', 'com.alibaba:fastjson', 'dependency', '2024-02-20', 'pending_review', 'Fastjson存在反序列化漏洞风险', 'admin'),
('DEF-2024-003', 'user_report', 'Spring核心模块安全风险', 'medium', 'security', 'org.springframework:spring-core', 'project', '2024-03-10', 'dispatched', 'Spring核心模块存在潜在安全问题', 'admin'),
('DEF-2024-004', 'sca', 'Nginx配置错误导致信息泄露', 'low', 'misconfiguration', 'nginx', 'project', '2024-04-05', 'ignored', 'Nginx配置不当可能泄露服务器信息', 'admin'),
('DEF-2024-005', 'external', 'MySQL权限提升漏洞', 'high', 'vulnerability', 'mysql', 'project', '2024-05-12', 'closed', 'MySQL存在权限提升漏洞', 'admin');

-- 9. 任务主表数据
INSERT INTO `governance_task` (`task_no`, `type`, `title`, `applicant`, `applicant_dept`, `status`, `priority`, `current_step`, `description`, `assigned_to`, `create_by`) VALUES
('TASK-2024-001', 'introduce', '引入Redis缓存组件', '张三', '技术研发部', 'completed', 'normal', '管理员审批通过', '项目需要使用Redis作为缓存层', '管理员', 'admin'),
('TASK-2024-002', 'upgrade', '升级Spring Boot到3.2版本', '李四', '平台研发部', 'processing', 'urgent', '待管理员审批', 'Spring Boot 3.2包含重要安全修复', '管理员', 'admin'),
('TASK-2024-003', 'introduce', '引入Kafka消息中间件', '王五', '技术研发部', 'pending', 'normal', '使用部门审批中', '项目需要消息队列支持', '使用部门专员', 'admin'),
('TASK-2024-004', 'defect', '修复Log4j漏洞', '赵六', '安全部', 'completed', 'urgent', '已修复', '处理DEF-2024-001缺陷', '张三', 'admin'),
('TASK-2024-005', 'withdraw', '退出Fastjson组件', '张三', '技术研发部', 'draft', 'normal', '草稿', '由于安全风险退出Fastjson', '', 'admin');

-- 10. 引入申请扩展表数据
INSERT INTO `task_introduce` (`task_id`, `component_name`, `component_version`, `license_type`, `introduce_reason`, `use_scenario`, `pre_eval_conclusion`) VALUES
(1, 'Redis', '7.0.11', 'BSD-3-Clause', '项目需要高性能缓存方案', '用于用户会话缓存和热点数据缓存', 'pass'),
(3, 'Kafka', '3.6.0', 'Apache-2.0', '系统间需要异步消息通信', '订单系统和消息平台之间的事件驱动', 'pass');

-- 11. 升级申请扩展表数据
INSERT INTO `task_upgrade` (`task_id`, `component_name`, `current_version`, `target_version`, `upgrade_reason`, `impact_analysis`, `pre_eval_conclusion`) VALUES
(2, 'Spring Boot', '2.7.18', '3.2.0', '包含重要安全修复和性能提升', '需要验证所有依赖兼容性，预计影响5个微服务', 'warning');

-- 12. 缺陷处置扩展表数据
INSERT INTO `task_defect` (`task_id`, `defect_id`, `defect_no`, `defect_type`, `affected_component`, `fix_plan`, `need_fix`, `production_change_no`) VALUES
(4, 1, 'DEF-2024-001', 'vulnerability', 'org.apache.log4j:log4j:2.14.1', '升级到log4j 2.17.1版本', 1, 'CHG-2024-001');

-- 13. 审批节点数据
INSERT INTO `approval_node` (`task_id`, `node_name`, `node_order`, `approver`, `status`, `opinion`, `approved_at`) VALUES
(1, '使用部门意见', 1, '使用部门专员', 'approved', '同意引入', '2024-01-16 10:00:00'),
(1, '管理员审批', 2, '管理员', 'approved', '审批通过', '2024-01-17 14:30:00'),
(2, '使用部门意见', 1, '使用部门专员', 'approved', '同意升级', '2024-02-21 09:00:00'),
(2, '管理员审批', 2, '管理员', 'pending', '', NULL),
(3, '使用部门意见', 1, '使用部门专员', 'pending', '', NULL),
(3, '管理员审批', 2, '管理员', 'not_required', '', NULL);

-- 14. 审批历史数据
INSERT INTO `approval_history` (`task_id`, `task_no`, `node_id`, `node_name`, `approver`, `action`, `opinion`, `create_time`) VALUES
(1, 'TASK-2024-001', 1, '使用部门意见', '使用部门专员', 'approve', '同意引入，Redis性能优秀', '2024-01-16 10:00:00'),
(1, 'TASK-2024-001', 2, '管理员审批', '管理员', 'approve', '审批通过，可以引入', '2024-01-17 14:30:00'),
(2, 'TASK-2024-002', 3, '使用部门意见', '使用部门专员', 'approve', '同意升级，但需要充分测试', '2024-02-21 09:00:00');

-- 15. 人员-资源映射数据
INSERT INTO `specialist_resource_mapping` (`resource_type`, `resource_id`, `resource_name`, `owner`, `owner_dept`, `maintainer`, `maintainer_dept`) VALUES
('project', 1, 'Redis', '张三', '技术研发部', '张三', '技术研发部'),
('project', 2, 'Spring Boot', '李四', '平台研发部', '李四', '平台研发部'),
('project', 3, 'MySQL', '王五', '运维部', '王五', '运维部'),
('project', 4, 'Nginx', '赵六', '安全部', '赵六', '安全部'),
('inner_project', 1, '订单系统', '张三', '技术研发部', '张三', '技术研发部'),
('inner_project', 2, '支付系统', '李四', '平台研发部', '李四', '平台研发部'),
('inner_project', 3, '用户中心', '王五', '技术研发部', '王五', '技术研发部');
