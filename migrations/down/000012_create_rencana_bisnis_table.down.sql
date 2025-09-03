-- Drop indexes first
DROP INDEX IF EXISTS idx_rencana_bisnis_id_bumd;
DROP INDEX IF EXISTS idx_rencana_bisnis_nomor;
DROP INDEX IF EXISTS idx_rencana_bisnis_tanggal;
DROP INDEX IF EXISTS idx_rencana_bisnis_masa_berlaku;
DROP INDEX IF EXISTS idx_rencana_bisnis_deleted_at;

-- Drop foreign key constraint
ALTER TABLE rencana_bisnis DROP CONSTRAINT IF EXISTS fk_rencana_bisnis_id_bumd;

-- Drop table
DROP TABLE IF EXISTS rencana_bisnis;
