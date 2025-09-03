CREATE TABLE rka (
    id                 BIGSERIAL       PRIMARY KEY,
    id_bumd            INT             NOT NULL,
    nomor              VARCHAR(100)    NOT NULL DEFAULT '',
    instansi_pemberi   VARCHAR(250)    NOT NULL DEFAULT '',
    tanggal            DATE            NOT NULL,
    kualifikasi        VARCHAR(100)    NOT NULL DEFAULT '',
    klasifikasi        VARCHAR(100)    NOT NULL DEFAULT '',
    masa_berlaku       DATE            NOT NULL,
    file               TEXT            NOT NULL DEFAULT '',
    created_at         TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by         INT             NOT NULL DEFAULT 0,
    updated_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by         INT             NOT NULL DEFAULT 0,
    deleted_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by         INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE rka 
ADD CONSTRAINT fk_rka_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_rka_id_bumd 
    ON rka(id_bumd);

CREATE INDEX idx_rka_nomor 
    ON rka(nomor);

CREATE INDEX idx_rka_tanggal 
    ON rka(tanggal);

CREATE INDEX idx_rka_masa_berlaku 
    ON rka(masa_berlaku);

CREATE INDEX idx_rka_deleted_at 
    ON rka(deleted_at);
