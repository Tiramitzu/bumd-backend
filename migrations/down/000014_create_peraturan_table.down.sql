-- Drop indexes first
DROP INDEX IF EXISTS idx_peraturan_id_bumd;
DROP INDEX IF EXISTS idx_peraturan_nomor;
DROP INDEX IF EXISTS idx_peraturan_jenis;
DROP INDEX IF EXISTS idx_peraturan_tanggal_berlaku;
DROP INDEX IF EXISTS idx_peraturan_deleted_at;

-- Drop foreign key constraint
ALTER TABLE peraturan DROP CONSTRAINT IF EXISTS fk_peraturan_id_bumd;

-- Drop table
DROP TABLE IF EXISTS peraturan;
