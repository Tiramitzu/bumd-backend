-- Drop indexes first
DROP INDEX IF EXISTS idx_rka_id_bumd;
DROP INDEX IF EXISTS idx_rka_nomor;
DROP INDEX IF EXISTS idx_rka_tanggal;
DROP INDEX IF EXISTS idx_rka_masa_berlaku;
DROP INDEX IF EXISTS idx_rka_deleted_at;

-- Drop foreign key constraint
ALTER TABLE rka DROP CONSTRAINT IF EXISTS fk_rka_id_bumd;

-- Drop table
DROP TABLE IF EXISTS rka;
