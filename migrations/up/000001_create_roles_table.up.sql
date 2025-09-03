CREATE TABLE roles (
    id          BIGSERIAL       PRIMARY KEY,
    nama        VARCHAR(250)    NOT NULL DEFAULT '',
    created_at  TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  INT             NOT NULL DEFAULT 0,
    updated_at  TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by  INT             NOT NULL DEFAULT 0,
    deleted_at  TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by  INT             NOT NULL DEFAULT 0
);

-- Populate data
INSERT INTO roles (nama, created_by) VALUES
('Admin Pusat', 0),
('Admin Daerah', 1),
('BUMD', 2);
