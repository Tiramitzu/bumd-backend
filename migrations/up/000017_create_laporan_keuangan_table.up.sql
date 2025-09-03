CREATE TABLE laporan_keuangan (
    id                 BIGSERIAL       PRIMARY KEY,
    id_bumd            INT             NOT NULL,
    tahun              INT             NOT NULL,
    periode            VARCHAR(50)     NOT NULL DEFAULT '',
    jenis_laporan      VARCHAR(100)    NOT NULL DEFAULT '',
    file               TEXT            NOT NULL DEFAULT '',
    keterangan         TEXT            NOT NULL DEFAULT '',
    created_at         TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by         INT             NOT NULL DEFAULT 0,
    updated_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by         INT             NOT NULL DEFAULT 0,
    deleted_at         TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by         INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE laporan_keuangan 
ADD CONSTRAINT fk_laporan_keuangan_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_laporan_keuangan_id_bumd 
    ON laporan_keuangan(id_bumd);

CREATE INDEX idx_laporan_keuangan_tahun 
    ON laporan_keuangan(tahun);

CREATE INDEX idx_laporan_keuangan_periode 
    ON laporan_keuangan(periode);

CREATE INDEX idx_laporan_keuangan_jenis_laporan 
    ON laporan_keuangan(jenis_laporan);

CREATE INDEX idx_laporan_keuangan_deleted_at 
    ON laporan_keuangan(deleted_at);
