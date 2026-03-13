-- Nusva KSPPS Database Schema
-- PostgreSQL Migration Script

-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- 1. TABLE: users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_locked BOOLEAN DEFAULT false,
    failed_login INTEGER DEFAULT 0,
    last_login TIMESTAMP,
    last_ip VARCHAR(45),
    force_change BOOLEAN DEFAULT false,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- ============================================
-- 2. TABLE: roles
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial roles
INSERT INTO roles (name, description) VALUES
('super_admin', 'Super Administrator with full access'),
('admin', 'Administrator with most access'),
('staff', 'Staff with limited access'),
('member', 'Member with self-access only')
ON CONFLICT (name) DO NOTHING;

-- ============================================
-- 3. TABLE: menus
-- ============================================
CREATE TABLE IF NOT EXISTS menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    path VARCHAR(200),
    parent_id INTEGER REFERENCES menus(id) ON DELETE SET NULL,
    "order" INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial menus
INSERT INTO menus (name, icon, path, parent_id, "order") VALUES
('Dashboard', 'layout-dashboard', '/dashboard', NULL, 1),
('Manajemen User', 'users', '/users', NULL, 2),
('Data Anggota', 'user-circle', '/members', NULL, 3),
('Simpanan', 'wallet', '/savings', NULL, 4),
('Pinjaman', 'banknote', '/loans', NULL, 5),
('Keamanan', 'shield', '/security', NULL, 6),
('Roles', 'shield-check', '/security/roles', 6, 1),
('Permissions', 'key', '/security/permissions', 6, 2),
('Audit Log', 'scroll-text', '/security/audit-logs', 6, 3)
ON CONFLICT DO NOTHING;

-- ============================================
-- 4. TABLE: permissions
-- ============================================
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    menu_id INTEGER NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL, -- view, create, edit, delete, approve, export
    code VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial permissions
INSERT INTO permissions (menu_id, action, code, name) VALUES
-- Dashboard
(1, 'view', 'dashboard.view', 'View Dashboard'),
-- User Management
(2, 'view', 'user_mgmt.view', 'View Users'),
(2, 'create', 'user_mgmt.create', 'Create User'),
(2, 'edit', 'user_mgmt.edit', 'Edit User'),
(2, 'delete', 'user_mgmt.delete', 'Delete User'),
-- Members
(3, 'view', 'member_mgmt.view', 'View Members'),
(3, 'create', 'member_mgmt.create', 'Create Member'),
(3, 'edit', 'member_mgmt.edit', 'Edit Member'),
(3, 'delete', 'member_mgmt.delete', 'Delete Member'),
-- Savings
(4, 'view', 'saving_mgmt.view', 'View Savings'),
(4, 'create', 'saving_mgmt.create', 'Create Saving Account'),
(4, 'edit', 'saving_mgmt.edit', 'Edit Saving Account'),
(4, 'delete', 'saving_mgmt.delete', 'Delete Saving Account'),
(4, 'export', 'saving_mgmt.export', 'Export Savings'),
-- Loans
(5, 'view', 'loan_mgmt.view', 'View Loans'),
(5, 'create', 'loan_mgmt.create', 'Create Loan Application'),
(5, 'edit', 'loan_mgmt.edit', 'Edit Loan Application'),
(5, 'delete', 'loan_mgmt.delete', 'Delete Loan Application'),
(5, 'approve', 'loan_mgmt.approve', 'Approve Loan'),
(5, 'export', 'loan_mgmt.export', 'Export Loans'),
-- Security/Roles
(7, 'view', 'role_mgmt.view', 'View Roles'),
(7, 'create', 'role_mgmt.create', 'Create Role'),
(7, 'edit', 'role_mgmt.edit', 'Edit Role'),
(7, 'delete', 'role_mgmt.delete', 'Delete Role'),
-- Security/Audit Logs
(9, 'view', 'audit.view', 'View Audit Logs'),
(9, 'export', 'audit.export', 'Export Audit Logs')
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 5. TABLE: role_permissions
-- ============================================
CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    is_allowed BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- ============================================
-- 6. TABLE: user_roles
-- ============================================
CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- ============================================
-- 7. TABLE: member_sequence_counter
-- ============================================
CREATE TABLE IF NOT EXISTS member_sequence_counter (
    id SERIAL PRIMARY KEY,
    year_code VARCHAR(2) NOT NULL,
    month_code VARCHAR(2) NOT NULL,
    gender_code VARCHAR(2) NOT NULL,
    last_seq INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(year_code, month_code, gender_code)
);

-- ============================================
-- 8. TABLE: members
-- ============================================
CREATE TABLE IF NOT EXISTS members (
    id SERIAL PRIMARY KEY,
    member_no VARCHAR(12) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    gender VARCHAR(20) NOT NULL CHECK (gender IN ('Laki-laki', 'Perempuan')),
    join_date DATE NOT NULL,
    birth_date DATE,
    birth_place VARCHAR(100),
    nik VARCHAR(16) UNIQUE,
    npwp VARCHAR(16),
    address TEXT NOT NULL,
    city VARCHAR(50),
    province VARCHAR(50),
    postal_code VARCHAR(10),
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(100),
    ktp_photo VARCHAR(255),
    npwp_photo VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_members_member_no ON members(member_no);
CREATE INDEX idx_members_nik ON members(nik);
CREATE INDEX idx_members_is_active ON members(is_active);

-- ============================================
-- 9. TABLE: saving_accounts
-- ============================================
CREATE TABLE IF NOT EXISTS saving_accounts (
    id SERIAL PRIMARY KEY,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    account_type VARCHAR(20) NOT NULL CHECK (account_type IN ('pokok', 'wajib', 'manasuka')),
    account_number VARCHAR(20) UNIQUE NOT NULL,
    balance DECIMAL(15, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(member_id, account_type)
);

CREATE INDEX idx_saving_accounts_member ON saving_accounts(member_id);
CREATE INDEX idx_saving_accounts_type ON saving_accounts(account_type);

-- ============================================
-- 10. TABLE: saving_transactions
-- ============================================
CREATE TABLE IF NOT EXISTS saving_transactions (
    id SERIAL PRIMARY KEY,
    saving_account_id INTEGER NOT NULL REFERENCES saving_accounts(id) ON DELETE CASCADE,
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('credit', 'debit')),
    amount DECIMAL(15, 2) NOT NULL,
    description TEXT,
    balance_before DECIMAL(15, 2),
    balance_after DECIMAL(15, 2),
    transaction_date DATE NOT NULL,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_saving_txns_account ON saving_transactions(saving_account_id);
CREATE INDEX idx_saving_txns_date ON saving_transactions(transaction_date);

-- ============================================
-- 11. TABLE: loan_applications
-- ============================================
CREATE TABLE IF NOT EXISTS loan_applications (
    id SERIAL PRIMARY KEY,
    application_no VARCHAR(20) UNIQUE NOT NULL,
    member_id INTEGER NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    principal_amount DECIMAL(15, 2) NOT NULL,
    margin_rate DECIMAL(5, 2) NOT NULL,
    margin_amount DECIMAL(15, 2) NOT NULL,
    term_months INTEGER NOT NULL,
    monthly_installment DECIMAL(15, 2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    purpose TEXT,
    contract_type VARCHAR(20) NOT NULL CHECK (contract_type IN ('murabahah', 'musyarakah')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'active', 'paid_off')),
    approved_by VARCHAR(100),
    approved_date TIMESTAMP,
    notes TEXT,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_applications_member ON loan_applications(member_id);
CREATE INDEX idx_loan_applications_status ON loan_applications(status);

-- ============================================
-- 12. TABLE: loan_schedules
-- ============================================
CREATE TABLE IF NOT EXISTS loan_schedules (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL REFERENCES loan_applications(id) ON DELETE CASCADE,
    sequence_number INTEGER NOT NULL,
    due_date DATE NOT NULL,
    principal DECIMAL(15, 2) NOT NULL,
    margin DECIMAL(15, 2) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT false,
    paid_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_schedules_application ON loan_schedules(application_id);
CREATE INDEX idx_loan_schedules_due_date ON loan_schedules(due_date);

-- ============================================
-- 13. TABLE: loan_transactions
-- ============================================
CREATE TABLE IF NOT EXISTS loan_transactions (
    id SERIAL PRIMARY KEY,
    application_id INTEGER NOT NULL REFERENCES loan_applications(id) ON DELETE CASCADE,
    schedule_id INTEGER REFERENCES loan_schedules(id) ON DELETE SET NULL,
    amount DECIMAL(15, 2) NOT NULL,
    transaction_date DATE NOT NULL,
    description TEXT,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_txns_application ON loan_transactions(application_id);
CREATE INDEX idx_loan_txns_date ON loan_transactions(transaction_date);

-- ============================================
-- 14. TABLE: audit_logs
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    username VARCHAR(100),
    action VARCHAR(50) NOT NULL,
    module VARCHAR(50) NOT NULL,
    record_id INTEGER,
    old_data TEXT,
    new_data TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    status VARCHAR(20) DEFAULT 'success',
    error_msg TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_module ON audit_logs(module);
CREATE INDEX idx_audit_created ON audit_logs(created_at);
CREATE INDEX idx_audit_action ON audit_logs(action);

-- ============================================
-- Functions for auto-updating updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menus_updated_at BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_members_updated_at BEFORE UPDATE ON members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_saving_accounts_updated_at BEFORE UPDATE ON saving_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_loan_applications_updated_at BEFORE UPDATE ON loan_applications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_loan_schedules_updated_at BEFORE UPDATE ON loan_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_member_sequence_counter_updated_at BEFORE UPDATE ON member_sequence_counter
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
