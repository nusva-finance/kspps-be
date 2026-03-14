-- Fix margin_setups table structure
-- Hapus dan buat ulang tabel dengan struktur yang benar

-- 1. Hapus tabel jika ada
DROP TABLE IF EXISTS margin_setups CASCADE;

-- 2. Buat tabel dengan struktur yang benar
CREATE TABLE margin_setups (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    tenor INTEGER NOT NULL,
    margin DECIMAL(10, 6) NOT NULL DEFAULT 0,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Buat index
CREATE INDEX idx_margin_setups_category ON margin_setups(category);
CREATE INDEX idx_margin_setups_tenor ON margin_setups(tenor);
CREATE INDEX idx_margin_setups_category_tenor ON margin_setups(category, tenor);

-- 4. Buat unique constraint
ALTER TABLE margin_setups
ADD CONSTRAINT unique_category_tenor UNIQUE (category, tenor);

-- 5. Insert sample data
INSERT INTO margin_setups (category, tenor, margin) VALUES
('Elektronik', 3, 0.05),
('Elektronik', 6, 0.07),
('Elektronik', 12, 0.10),
('Pakaian', 3, 0.04),
('Pakaian', 6, 0.06),
('Pakaian', 12, 0.08),
('Kendaraan', 12, 0.12),
('Kendaraan', 24, 0.15),
('Kendaraan', 36, 0.18);

-- 6. Verifikasi
SELECT 'Table created successfully' AS status;
SELECT COUNT(*) AS total_rows FROM margin_setups;
