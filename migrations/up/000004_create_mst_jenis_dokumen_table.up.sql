CREATE TABLE mst_jenis_dokumen (
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

-- Populate data
INSERT INTO mst_jenis_dokumen (nama, deskripsi, created_by) VALUES
('SOP', '', 1),
('Dokumen Tata Kelola', '', 1),
('Pengadaan Barang Jasa', '', 1),
('Kerjasama / SPK', '', 1),
('Pinjaman', '', 1),
('Rencana Bisnis', '', 1),
('Rencana Kerja Anggaran', '', 1);