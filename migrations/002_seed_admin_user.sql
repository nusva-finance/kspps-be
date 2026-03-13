-- Nusva KSPPS - Seed Admin User
-- Default credentials: username: admin, password: admin123

-- Insert default admin user
-- Password: admin123 (hashed with bcrypt)
INSERT INTO users (username, email, password_hash, full_name, phone_number, is_active) VALUES
('admin', 'admin@nusva.id', '$2a$14$apHyMawxtvQUg3cl9evHGuaStv6LNJT/vPONRfY1C2iCEY3sdBTCe', 'Administrator', '081234567890', true)
ON CONFLICT (username) DO NOTHING;

-- Assign admin role to admin user
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.username = 'admin' AND r.name = 'admin'
ON CONFLICT (user_id, role_id) DO NOTHING;
