CREATE TABLE mst_akta_notaris (
    id              BIGSERIAL       PRIMARY KEY,
    id_bumd         INT             NOT NULL,
    nomor_perda     VARCHAR(100)    NOT NULL DEFAULT '',
    tanggal_perda   DATE            NOT NULL,
    modal_dasar     DECIMAL(15,2)   NOT NULL DEFAULT 0,
    keterangan      TEXT            NOT NULL DEFAULT '',
    file            TEXT            NOT NULL DEFAULT '',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT             NOT NULL DEFAULT 0,
    updated_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by      INT             NOT NULL DEFAULT 0,
    deleted_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by      INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE mst_akta_notaris 
ADD CONSTRAINT fk_mst_akta_notaris_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_mst_akta_notaris_id_bumd 
    ON mst_akta_notaris(id_bumd);

CREATE INDEX idx_mst_akta_notaris_nomor_perda 
    ON mst_akta_notaris(nomor_perda);

CREATE INDEX idx_mst_akta_notaris_tanggal_perda 
    ON mst_akta_notaris(tanggal_perda);

CREATE INDEX idx_mst_akta_notaris_deleted_at 
    ON mst_akta_notaris(deleted_at);
