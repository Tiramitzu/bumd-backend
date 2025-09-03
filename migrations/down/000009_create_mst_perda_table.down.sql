-- Drop indexes first
DROP INDEX IF EXISTS idx_mst_perda_id_bumd;
DROP INDEX IF EXISTS idx_mst_perda_nomor_perda;
DROP INDEX IF EXISTS idx_mst_perda_tanggal_perda;
DROP INDEX IF EXISTS idx_mst_perda_deleted_at;

-- Drop foreign key constraint
ALTER TABLE mst_perda DROP CONSTRAINT IF EXISTS fk_mst_perda_id_bumd;

-- Drop table
DROP TABLE IF EXISTS mst_perda;
