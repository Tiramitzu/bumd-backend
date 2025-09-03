-- Drop indexes first
DROP INDEX IF EXISTS idx_mst_bentuk_usaha_nama;
DROP INDEX IF EXISTS idx_mst_bentuk_usaha_deleted_at;

-- Drop table
DROP TABLE IF EXISTS mst_bentuk_usaha;
