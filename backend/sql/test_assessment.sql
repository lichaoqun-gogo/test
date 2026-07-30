DELETE FROM introduce_assessment WHERE task_id = 'TASK-001';

INSERT INTO introduce_assessment (task_id, component_name, assess_dimension, assess_status, deviation_explanation) VALUES 
('TASK-001', 'Vue.js', 'functionality', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'security', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'reliability', 'deviation', '社区长期支持版本不足3年，企业版支持有限。计划建立内部技术团队维护机制，定期跟进社区更新。'),
('TASK-001', 'Vue.js', 'performance', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'compatibility', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'maintainability', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'extensibility', 'fully_meet', NULL),
('TASK-001', 'Vue.js', 'usability', 'fully_meet', NULL);
