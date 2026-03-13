-- Add deleted_at column to users table for GORM soft delete
-- This script adds the missing deleted_at column without dropping constraints

-- Add deleted_at column
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Create index on deleted_at for better query performance
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Verify the column was added
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'users' AND column_name = 'deleted_at';
