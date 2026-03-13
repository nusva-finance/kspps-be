-- Add missing roles for KSPPS system
INSERT INTO roles (name, description, is_active) VALUES
('manager', 'Manager with full management access', true),
('teller', 'Teller with transaction access', true),
('cs', 'Customer Service with limited access', true)
ON CONFLICT (name) DO NOTHING;
