CREATE TABLE m_bentuk_usaha (
    id_bu               uuid NOT NULL,
    nama_bu             varchar(250) NOT NULL DEFAULT ''::character varying,
    deskripsi_bu        text NOT NULL DEFAULT '',
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          int4 NOT NULL DEFAULT 0,
    updated_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by          int4 NOT NULL DEFAULT 0,
    deleted_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by          int4 NOT NULL DEFAULT 0,
    CONSTRAINT m_bentuk_usaha_pkey PRIMARY KEY (id_bu)
);

-- Add index on nama for better search performance
CREATE INDEX idx_m_bentuk_usaha_nama 
    ON m_bentuk_usaha(nama_bu);

-- Add index on deleted_at for soft delete queries
CREATE INDEX idx_m_bentuk_usaha_deleted_at 
    ON m_bentuk_usaha(deleted_at);

-- Populate data
INSERT INTO m_bentuk_usaha (id_bu, nama_bu, deskripsi_bu, created_by) VALUES
('01994c01-c285-7e57-a486-fd9978083917', 'Air Minum', '', 1),
('01994c01-c285-7744-ba6a-5782f39fe366', 'Aneka Usaha', '', 1),
('01994c01-c285-7eea-aead-221a0d1f4cac', 'BPD', '', 1),
('01994c01-c285-7e6f-865a-b286738dff03', 'BPR', '', 1),
('01994c01-c285-73ae-8885-65ab6b20deb3', 'JAMKRIDA', '', 1),
('01994c01-c285-7e51-9140-e75269630f28', 'Pasar', '', 1),
('01994c01-c285-73b3-ada9-10d5180a4a2f', 'Lainnya', '', 1),
('01994c01-c285-75df-906f-5b5ea775d519', 'Migas', 'Migas', 1),
('01994c01-c285-764f-ab66-d9e0c728ba59', 'Pariwisata', 'Pariwisata', 1),
('01994c01-c285-7880-a971-70231530ba47', 'MINERBA', 'Minerba', 1),
('01994c01-c285-7e19-a547-c055c71b681b', 'AGRO', 'Agro', 1),
('01994c01-c285-72e8-a39a-4b065fa81980', 'Kepelabuhan', 'Kepelabuhan', 1),
('01994c01-c285-72c6-9bec-5ebe9fbd222d', 'Transportasi', 'Transportasi', 1);