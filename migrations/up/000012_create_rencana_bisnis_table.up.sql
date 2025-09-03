CREATE TABLE rencana_bisnis (
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
ALTER TABLE rencana_bisnis 
ADD CONSTRAINT fk_rencana_bisnis_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_rencana_bisnis_id_bumd 
    ON rencana_bisnis(id_bumd);

CREATE INDEX idx_rencana_bisnis_nomor 
    ON rencana_bisnis(nomor);

CREATE INDEX idx_rencana_bisnis_tanggal 
    ON rencana_bisnis(tanggal);

CREATE INDEX idx_rencana_bisnis_masa_berlaku 
    ON rencana_bisnis(masa_berlaku);

CREATE INDEX idx_rencana_bisnis_deleted_at 
    ON rencana_bisnis(deleted_at);
