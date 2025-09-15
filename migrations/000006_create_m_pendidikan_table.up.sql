CREATE TABLE m_pendidikan (
    id_pendidikan       uuid NOT NULL,
    nama_pendidikan     varchar(255) NOT NULL DEFAULT ''::character varying,
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          int4 NOT NULL DEFAULT 0,
    updated_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by          int4 NOT NULL DEFAULT 0,
    deleted_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by          int4 NOT NULL DEFAULT 0,
    CONSTRAINT m_pendidikan_pkey PRIMARY KEY (id_pendidikan)
);

CREATE INDEX idx_m_pendidikan_nama_pendidikan ON m_pendidikan(nama_pendidikan);

-- Populate data
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-7196-a6e0-6b2553714335', 'SD (Sekolah Dasar) dan Sederajat', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-73fb-8922-08a7897f572f', 'SMP (Sekolah Menengah Pertama) dan Sederajat', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-7f11-95c7-a4cb89959cd7', 'SMA (Sekolah Menengah Atas) dan Sederajat', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-7119-b388-65b033c8a343', 'SMK (Sekolah Menengah Kejuruan) dan Sederajat', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-7936-a7b2-5da9de1cac12', 'D1', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-74b4-a83b-277881ad4e9c', 'D2', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-78eb-8210-282e85a289cd', 'D3 (Ahli Madya)', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-70a9-9773-5ac917ad7e2f', 'D4 (Sarjana Terapan)', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-785f-92e4-a5e6ab8ae25f', 'S1 (Sarjana)', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c75-2cf3-736e-9e23-2bbabb1d2b36', 'S2 (Magister)', 1, 1, 1);
INSERT INTO m_pendidikan (id_pendidikan, nama_pendidikan, created_by, updated_by, deleted_by) VALUES ('01994c79-6d4d-7e5e-9d30-3d28773ae539', 'S3 (Doktor)', 1, 1, 1);
