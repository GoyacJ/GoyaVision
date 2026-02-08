-- Create system_configs table
CREATE TABLE system_configs (
    key VARCHAR(255) PRIMARY KEY,
    value JSONB NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default configuration
INSERT INTO system_configs (key, value, description)
VALUES ('system.home_path', '"/assets"', '系统默认首页路径');

INSERT INTO system_configs (key, value, description)
VALUES ('system.public_menus', '[]', '未登录用户可见的菜单ID列表');

-- Add permission for system config
INSERT INTO permissions (id, code, name, description, created_at, updated_at)
VALUES (gen_random_uuid(), 'system:config:view', '查看系统配置', '允许查看系统配置', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO permissions (id, code, name, description, created_at, updated_at)
VALUES (gen_random_uuid(), 'system:config:update', '修改系统配置', '允许修改系统配置', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Assign permissions to super_admin role (assuming role exists and we find it by code 'super_admin')
-- Note: In this system, super_admin has '*' permission usually, but for completeness or other roles:
-- We'll skip explicit role assignment here as super_admin usually has all permissions.

-- Add System Config menu
INSERT INTO menus (id, parent_id, code, name, type, path, icon, component, permission, sort, visible, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    NULL,
    'system',
    '系统管理',
    1,
    '/system',
    'Setting',
    'Layout',
    'system:config:view',
    90,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

INSERT INTO menus (id, parent_id, code, name, type, path, icon, component, permission, sort, visible, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    m.id,
    'system_config',
    '系统配置',
    2,
    '/system/config',
    'Tools',
    'system/config/index',
    'system:config:view',
    1,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM menus m WHERE m.code = 'system';
