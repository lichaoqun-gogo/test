-- CICD检查测试数据
INSERT INTO component_black (blacklist_id, component_name, component_version, component_group, project_id, project_name, alert_type, status, handler, handle_plan, ci_pipeline_id, create_time, create_by, update_time, update_by) VALUES
(1, 'log4j-core', '2.14.1', 'org.apache.logging.log4j', 1, 'ORDER-SYS', 1, '待处理', '', '', 'pipeline-001', NOW(), 'admin', NOW(), 'admin'),
(2, 'commons-collections', '3.2.1', 'org.apache.commons', 1, 'ORDER-SYS', 1, '待处理', '', '', 'pipeline-001', NOW(), 'admin', NOW(), 'admin'),
(3, 'jackson-databind', '2.9.8', 'com.fasterxml.jackson.core', 2, 'PAY-SYS', 1, '待处理', '', '', 'pipeline-002', NOW(), 'admin', NOW(), 'admin'),
(4, 'spring-boot-starter-web', '2.7.0', 'org.springframework.boot', 3, 'USER-CENTER', 0, '待处理', '', '', 'pipeline-003', NOW(), 'admin', NOW(), 'admin');
