-- Drop indexes first
DROP INDEX IF EXISTS idx_mst_akta_notaris_id_bumd;
DROP INDEX IF EXISTS idx_mst_akta_notaris_nomor_perda;
DROP INDEX IF EXISTS idx_mst_akta_notaris_tanggal_perda;
DROP INDEX IF EXISTS idx_mst_akta_notaris_deleted_at;

-- Drop foreign key constraint
ALTER TABLE mst_akta_notaris DROP CONSTRAINT IF EXISTS fk_mst_akta_notaris_id_bumd;

-- Drop table
DROP TABLE IF EXISTS mst_akta_notaris;
