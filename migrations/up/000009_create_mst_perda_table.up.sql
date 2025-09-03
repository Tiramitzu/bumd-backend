CREATE TABLE mst_perda (
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
ALTER TABLE mst_perda 
ADD CONSTRAINT fk_mst_perda_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_mst_perda_id_bumd 
    ON mst_perda(id_bumd);

CREATE INDEX idx_mst_perda_nomor_perda 
    ON mst_perda(nomor_perda);

CREATE INDEX idx_mst_perda_tanggal_perda 
    ON mst_perda(tanggal_perda);

CREATE INDEX idx_mst_perda_deleted_at 
    ON mst_perda(deleted_at);
