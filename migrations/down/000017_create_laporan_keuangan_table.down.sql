-- Drop indexes first
DROP INDEX IF EXISTS idx_laporan_keuangan_id_bumd;
DROP INDEX IF EXISTS idx_laporan_keuangan_tahun;
DROP INDEX IF EXISTS idx_laporan_keuangan_periode;
DROP INDEX IF EXISTS idx_laporan_keuangan_jenis_laporan;
DROP INDEX IF EXISTS idx_laporan_keuangan_deleted_at;

-- Drop foreign key constraint
ALTER TABLE laporan_keuangan DROP CONSTRAINT IF EXISTS fk_laporan_keuangan_id_bumd;

-- Drop table
DROP TABLE IF EXISTS laporan_keuangan;
