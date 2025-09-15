CREATE TABLE m_bentuk_badan_hukum (
    id_bbh                  uuid NOT NULL,
    nama_bbh                varchar(250) NOT NULL DEFAULT ''::character varying,
    deskripsi_bbh           text NOT NULL DEFAULT '',
    created_at              timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              int4 NOT NULL DEFAULT 0,
    updated_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by              int4 NOT NULL DEFAULT 0,
    deleted_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by              int4 NOT NULL DEFAULT 0,
    CONSTRAINT m_bentuk_badan_hukum_pkey PRIMARY KEY (id_bbh)
);

-- Add index on nama for better search performance
CREATE INDEX idx_m_bentuk_badan_hukum_nama 
    ON m_bentuk_badan_hukum(nama_bbh);

-- Add index on deleted_at for soft delete queries
CREATE INDEX idx_m_bentuk_badan_hukum_deleted_at 
    ON m_bentuk_badan_hukum(deleted_at);

-- Populate data
INSERT INTO m_bentuk_badan_hukum (id_bbh, nama_bbh, deskripsi_bbh, created_by) VALUES
('01994c00-6768-7fde-ab34-f315eea6510f', 'Perumda', '', 1),
('01994c00-6768-74b2-9895-9ad47078f424', 'Perseroda', '', 1),
('01994c00-6768-7c83-8ca6-57c2b2a1f2cc', 'Lainnya', '', 1),
('01994c00-6768-78cf-b405-261ff89e11bb', 'Perusda', '', 1);