-- Drop indexes first
DROP INDEX IF EXISTS idx_bumd_id_daerah;
DROP INDEX IF EXISTS idx_bumd_id_bentuk_hukum;
DROP INDEX IF EXISTS idx_bumd_id_bentuk_usaha;
DROP INDEX IF EXISTS idx_bumd_nama;
DROP INDEX IF EXISTS idx_bumd_deleted_at;

-- Drop foreign key constraints
ALTER TABLE bumd DROP CONSTRAINT IF EXISTS fk_bumd_id_bentuk_hukum;
ALTER TABLE bumd DROP CONSTRAINT IF EXISTS fk_bumd_id_bentuk_usaha;

-- Drop table
DROP TABLE IF EXISTS bumd;
