-- Migration: Make Savings Types Dynamic
-- Create saving_types table and update saving_accounts to use foreign keys

-- ============================================
-- 1. TABLE: saving_types
-- ============================================
CREATE TABLE IF NOT EXISTS saving_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    is_required BOOLEAN DEFAULT false,
    min_balance DECIMAL(15, 2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster lookups
CREATE INDEX idx_saving_types_code ON saving_types(code);
CREATE INDEX idx_saving_types_active ON saving_types(is_active);

-- ============================================
-- 2. Insert Initial Saving Types Data
-- ============================================
INSERT INTO saving_types (code, name, description, is_required, min_balance, display_order) VALUES
('pokok', 'Simpanan Pokok', 'Simpanan wajib pertama kali menjadi anggota koperasi', true, 0, 1),
('wajib', 'Simpanan Wajib', 'Simpanan rutin yang wajib dibayar oleh anggota', true, 0, 2),
('modal', 'Simpanan Modal', 'Simpanan modal untuk keperluan investasi dan pengembangan usaha', false, 0, 3)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- 3. Alter Table saving_accounts - Add Foreign Key
-- ============================================

-- First, rename the old account_type column to account_type_old for migration
ALTER TABLE saving_accounts RENAME COLUMN account_type TO account_type_old;

-- Add new column as foreign key to saving_types
ALTER TABLE saving_accounts ADD COLUMN account_type_id INTEGER REFERENCES saving_types(id);

-- ============================================
-- 4. Migrate Existing Data
-- ============================================

-- Update existing records with mapping from account_type_old to account_type_id
UPDATE saving_accounts sa
SET account_type_id = st.id
FROM saving_types st
WHERE st.code = sa.account_type_old;

-- ============================================
-- 5. Apply Constraints
-- ============================================

-- Set NOT NULL constraint on account_type_id
ALTER TABLE saving_accounts ALTER COLUMN account_type_id SET NOT NULL;

-- Drop the old column
ALTER TABLE saving_accounts DROP COLUMN account_type_old;

-- ============================================
-- 6. Create Functions for Helper Queries
-- ============================================

-- Function to get saving type name by ID
CREATE OR REPLACE FUNCTION get_saving_type_name(type_id INTEGER)
RETURNS VARCHAR AS $$
BEGIN
    RETURN COALESCE((SELECT name FROM saving_types WHERE id = type_id), 'Unknown');
END;
$$ LANGUAGE plpgsql;

-- Function to get saving type code by ID
CREATE OR REPLACE FUNCTION get_saving_type_code(type_id INTEGER)
RETURNS VARCHAR AS $$
BEGIN
    RETURN COALESCE((SELECT code FROM saving_types WHERE id = type_id), 'unknown');
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 7. Update Trigger for saving_types
-- ============================================
CREATE TRIGGER update_saving_types_updated_at BEFORE UPDATE ON saving_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 8. Create View for Easy Account Queries
-- ============================================
CREATE OR REPLACE VIEW v_saving_accounts_with_types AS
SELECT
    sa.id,
    sa.member_id,
    sa.account_type_id,
    sa.account_number,
    sa.balance,
    sa.is_active,
    sa.created_at,
    sa.updated_at,
    st.code AS account_type_code,
    st.name AS account_type_name,
    st.description AS account_type_description,
    st.is_required,
    st.min_balance
FROM saving_accounts sa
LEFT JOIN saving_types st ON sa.account_type_id = st.id;

-- ============================================
-- 9. Create Index for Better Performance
-- ============================================
CREATE INDEX idx_saving_accounts_type_id ON saving_accounts(account_type_id);

-- ============================================
-- Migration Complete
-- ============================================
-- Notes:
-- - saving_types table now controls all savings types
-- - saving_accounts.account_type_id references saving_types.id
-- - Old hardcoded 'pokok', 'wajib', 'manasuka' strings migrated to IDs
-- - View v_saving_accounts_with_types provides easy access to account types
-- - Helper functions get_saving_type_name() and get_saving_type_code() available
