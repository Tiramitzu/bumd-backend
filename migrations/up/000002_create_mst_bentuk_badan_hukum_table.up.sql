CREATE TABLE mst_bentuk_badan_hukum (
    id              BIGSERIAL       PRIMARY KEY,
    nama            VARCHAR(250)    NOT NULL DEFAULT '',
    deskripsi       TEXT            NOT NULL DEFAULT '',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT             NOT NULL DEFAULT 0,
    updated_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by      INT             NOT NULL DEFAULT 0,
    deleted_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by      INT             NOT NULL DEFAULT 0
);

-- Add index on nama for better search performance
CREATE INDEX idx_mst_bentuk_badan_hukum_nama 
    ON mst_bentuk_badan_hukum(nama);

-- Add index on deleted_at for soft delete queries
CREATE INDEX idx_mst_bentuk_badan_hukum_deleted_at 
    ON mst_bentuk_badan_hukum(deleted_at);

-- Populate data
INSERT INTO mst_bentuk_badan_hukum (nama, deskripsi, created_by) VALUES
('Perumda', '', 1),
('Perseroda', '', 1),
('Lainnya', '', 1),
('Perusda', '', 1)