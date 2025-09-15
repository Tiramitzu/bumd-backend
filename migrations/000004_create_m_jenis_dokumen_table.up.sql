CREATE TABLE m_jenis_dokumen (
    id_jd               uuid NOT NULL,
    nama_jd             varchar(250) NOT NULL DEFAULT ''::character varying,
    deskripsi_jd        text NOT NULL DEFAULT '',
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          int4 NOT NULL DEFAULT 0,
    updated_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by          int4 NOT NULL DEFAULT 0,
    deleted_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by          int4 NOT NULL DEFAULT 0,
    CONSTRAINT m_jenis_dokumen_pkey PRIMARY KEY (id_jd)
);

-- Populate data
INSERT INTO m_jenis_dokumen (id_jd, nama_jd, deskripsi_jd, created_by) VALUES
('01994c04-699d-75e0-a288-13980f8c854d', 'SOP', '', 1),
('01994c04-699d-7b41-a4bd-b553f36dc916', 'Dokumen Tata Kelola', '', 1),
('01994c04-699d-7912-9c24-9e088476e909', 'Pengadaan Barang Jasa', '', 1),
('01994c04-699d-76ce-b71a-1871df178927', 'Kerjasama / SPK', '', 1),
('01994c04-699d-7aa0-81e4-4ccbf1356654', 'Pinjaman', '', 1),
('01994c04-699d-791a-b108-80190c082bfe', 'Rencana Bisnis', '', 1),
('01994c04-699d-7fc1-a7b7-232c84796b9f', 'Rencana Kerja Anggaran', '', 1);