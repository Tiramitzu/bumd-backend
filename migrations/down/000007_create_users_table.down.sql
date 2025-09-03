DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_id_daerah;
DROP INDEX IF EXISTS idx_users_id_role;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop foreign key constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_id_role;

-- Drop table
DROP TABLE IF EXISTS users;
