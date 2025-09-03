-- Drop indexes first
DROP INDEX IF EXISTS idx_kinerja_id_bumd;
DROP INDEX IF EXISTS idx_kinerja_tahun;
DROP INDEX IF EXISTS idx_kinerja_deleted_at;

-- Drop foreign key constraint
ALTER TABLE kinerja DROP CONSTRAINT IF EXISTS fk_kinerja_id_bumd;

-- Drop table
DROP TABLE IF EXISTS kinerja;
