CREATE TABLE mst_bentuk_usaha (
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
CREATE INDEX idx_mst_bentuk_usaha_nama 
    ON mst_bentuk_usaha(nama);

-- Add index on deleted_at for soft delete queries
CREATE INDEX idx_mst_bentuk_usaha_deleted_at 
    ON mst_bentuk_usaha(deleted_at);

-- Populate data
INSERT INTO mst_bentuk_usaha (nama, deskripsi, created_by) VALUES
('Air Minum', '', 1),
('Aneka Usaha', '', 1),
('BPD', '', 1),
('BPR', '', 1),
('JAMKRIDA', '', 1),
('Pasar', '', 1),
('Lainnya', '', 1),
('Migas', '', 1),
('Pariwisata', '', 1),
('MINERBA', '', 1),
('AGRO', '', 1),
('Kepelabuhan', '', 1),
('Transportasi', '', 1)