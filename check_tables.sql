-- Cek struktur tabel margin_setups
-- Jalankan query ini di database untuk melihat struktur tabel

SELECT
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'margin_setups'
ORDER BY ordinal_position;

-- Cek apakah tabel ada
SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'margin_setups'
);

-- Cek data sample
SELECT * FROM margin_setups LIMIT 5;
