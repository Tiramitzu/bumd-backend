CREATE TABLE mst_produk (
    id              BIGSERIAL       PRIMARY KEY,
    id_bumd         INT             NOT NULL DEFAULT 0,
    id_daerah       INT             NOT NULL DEFAULT 0,
    nama            VARCHAR(250)    NOT NULL DEFAULT '',
    deskripsi       TEXT            NOT NULL DEFAULT '',
    gambar          TEXT            NOT NULL DEFAULT '',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT             NOT NULL DEFAULT 0,
    updated_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by      INT             NOT NULL DEFAULT 0,
    deleted_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by      INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE mst_produk 
ADD CONSTRAINT fk_mst_produk_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_mst_produk_id_bumd 
    ON mst_produk(id_bumd);

CREATE INDEX idx_mst_produk_id_daerah 
    ON mst_produk(id_daerah);

CREATE INDEX idx_mst_produk_deleted_at 
    ON mst_produk(deleted_at);