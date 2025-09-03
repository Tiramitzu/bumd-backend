CREATE TABLE peraturan (
    id                 BIGSERIAL       PRIMARY KEY,
    id_bumd            INT             NOT NULL,
    nomor              VARCHAR(100)    NOT NULL DEFAULT '',
    jenis              VARCHAR(100)    NOT NULL DEFAULT '',
    tanggal_berlaku    DATE            NOT NULL,
    keterangan         TEXT            NOT NULL DEFAULT '',
    file               TEXT            NOT NULL DEFAULT '',
    created_at         TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by         INT             NOT NULL DEFAULT 0,
    updated_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by         INT             NOT NULL DEFAULT 0,
    deleted_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by         INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE peraturan 
ADD CONSTRAINT fk_peraturan_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_peraturan_id_bumd 
    ON peraturan(id_bumd);

CREATE INDEX idx_peraturan_nomor 
    ON peraturan(nomor);

CREATE INDEX idx_peraturan_jenis 
    ON peraturan(jenis);

CREATE INDEX idx_peraturan_tanggal_berlaku 
    ON peraturan(tanggal_berlaku);

CREATE INDEX idx_peraturan_deleted_at 
    ON peraturan(deleted_at);
