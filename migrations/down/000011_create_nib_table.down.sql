-- Drop indexes first
DROP INDEX IF EXISTS idx_nib_id_bumd;
DROP INDEX IF EXISTS idx_nib_nomor;
DROP INDEX IF EXISTS idx_nib_tanggal;
DROP INDEX IF EXISTS idx_nib_masa_berlaku;
DROP INDEX IF EXISTS idx_nib_deleted_at;

-- Drop foreign key constraint
ALTER TABLE nib DROP CONSTRAINT IF EXISTS fk_nib_id_bumd;

-- Drop table
DROP TABLE IF EXISTS nib;
