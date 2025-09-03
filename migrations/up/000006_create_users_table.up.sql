CREATE TABLE users (
    id          BIGSERIAL       PRIMARY KEY,
    username    VARCHAR(200)    NOT NULL DEFAULT '',
    password    VARCHAR(150)    NOT NULL DEFAULT '',
    id_daerah   INT             NOT NULL DEFAULT 0,
    id_role     INT             NOT NULL DEFAULT 0,
    nama        VARCHAR(250)    NOT NULL DEFAULT '',
    logo        TEXT            NOT NULL DEFAULT '',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  INT             NOT NULL DEFAULT 0,
    updated_at  TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by  INT             NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by  INT             NOT NULL DEFAULT 0
);

-- Populate data
INSERT INTO users (username, password, id_daerah, id_role, nama, logo, created_by) VALUES
('admin', '$2a$12$c1oTnbMoQM/VuwAW2xXNhehGlTecToWCyfF1gKiv0wFg47BAc0abu', 0, 1, 'Admin Pusat', '', 1);

-- Add foreign key constraint to roles table
ALTER TABLE users 
ADD CONSTRAINT fk_users_id_role 
FOREIGN KEY (id_role) REFERENCES roles(id);

-- Add indexes for better performance
CREATE INDEX idx_users_username 
    ON users(username);

CREATE INDEX idx_users_id_daerah 
    ON users(id_daerah);

CREATE INDEX idx_users_id_role 
    ON users(id_role);

CREATE INDEX idx_users_deleted_at 
    ON users(deleted_at);
