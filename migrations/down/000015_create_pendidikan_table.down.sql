-- Drop indexes first
DROP INDEX IF EXISTS idx_pendidikan_id_bumd;
DROP INDEX IF EXISTS idx_pendidikan_nomor;
DROP INDEX IF EXISTS idx_pendidikan_tanggal;
DROP INDEX IF EXISTS idx_pendidikan_masa_berlaku;
DROP INDEX IF EXISTS idx_pendidikan_deleted_at;

-- Drop foreign key constraint
ALTER TABLE pendidikan DROP CONSTRAINT IF EXISTS fk_pendidikan_id_bumd;

-- Drop table
DROP TABLE IF EXISTS pendidikan;
