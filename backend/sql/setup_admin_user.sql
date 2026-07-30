-- 插入或更新 admin 用户，role=admin
INSERT INTO users (username, nickname, email, phone, department, role, can_toggle_admin_view, status, create_at, update_at)
VALUES ('admin', '系统管理员', 'admin@qim.local', '13800000000', '管理层', 'admin', true, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    role = 'admin',
    can_toggle_admin_view = true,
    nickname = '系统管理员',
    department = '管理层',
    email = 'admin@qim.local',
    phone = '13800000000',
    status = 1;

-- 添加为部门开源治理专员
INSERT INTO dept_specialist_mapping (dept, specialist_id, specialist_name, specialist_role, is_active)
VALUES ('管理层', 'admin', '系统管理员', 'dept_specialist', 1)
ON DUPLICATE KEY UPDATE
    specialist_name = '系统管理员',
    is_active = 1;
