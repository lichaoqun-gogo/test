-- 设置 lisi 为 admin 角色（管理部门开源治理专员，可切换全局治理/个人工作台）
-- CAS 登录时会更新 nickname/department 等信息，但不会覆盖 role
INSERT INTO users (username, nickname, email, phone, department, role, status, create_at, update_at)
VALUES ('lisi', '李四', 'lisi@qim.local', '13800000002', '开源治理办公室', 'admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    role = 'admin',
    status = 1;

-- 设置 lisi 为开源治理办公室的管理部门专员
-- 直接 UPDATE 已有记录（唯一键 uk_dept_role 约束：每个部门每个角色只能有一人）
UPDATE dept_specialist_mapping
SET specialist_id = 'lisi', specialist_name = '李四', is_active = 1
WHERE dept = '开源治理办公室' AND specialist_role = 'mgmt_dept_specialist';
