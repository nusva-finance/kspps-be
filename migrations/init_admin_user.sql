-- Initialize Admin User and Role
-- Run this only once to add default admin user

-- Check if admin user already exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin') THEN
        RAISE NOTICE 'Admin user already exists, skipping';
    ELSE
        -- Insert admin user (password: admin123)
        INSERT INTO users (username, email, password_hash, full_name, phone_number, is_active) VALUES
        ('admin', 'admin@nusva.id', '$2a$14$apHyMawxtvQUg3cl9evHGuaStv6LNJT/vPONRfY1C2iCEY3sdBTCe', 'Administrator', '081234567890', true);

        -- Assign admin role to admin user
        INSERT INTO user_roles (user_id, role_id)
        SELECT u.id, r.id
        FROM users u
        CROSS JOIN roles r
        WHERE u.username = 'admin' AND r.name = 'admin';
    END IF;
END $$;
