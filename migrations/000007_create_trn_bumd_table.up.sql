CREATE TABLE trn_bumd (
    id_bumd                 uuid NOT NULL,
    id_daerah               int4 NOT NULL DEFAULT 0,
    id_bentuk_hukum         uuid NOT NULL,
    id_bentuk_usaha         uuid NOT NULL,
    id_induk_perusahaan     uuid,
    penerapan_spi_bumd      boolean NOT NULL DEFAULT false,
    file_spi_bumd           text NOT NULL DEFAULT '',
    nama_bumd               varchar(250) NOT NULL DEFAULT ''::character varying,
    deskripsi_bumd          text NOT NULL DEFAULT '',
    alamat_bumd             text NOT NULL DEFAULT '',
    no_telp_bumd            varchar(30) NOT NULL DEFAULT ''::character varying,
    no_fax_bumd             varchar(30) NOT NULL DEFAULT ''::character varying,
    email_bumd              varchar(250) NOT NULL DEFAULT ''::character varying,
    website_bumd            varchar(250) NOT NULL DEFAULT ''::character varying,
    narahubung_bumd         varchar(30) NOT NULL DEFAULT ''::character varying,
    npwp_bumd               varchar(20) NOT NULL DEFAULT ''::character varying,
    npwp_pemberi_bumd       varchar(250) NOT NULL DEFAULT ''::character varying,
    npwp_file_bumd          text NOT NULL DEFAULT '',
    logo_bumd               text NOT NULL DEFAULT '',
    created_at              timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              int4 NOT NULL DEFAULT 0,
    updated_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by              int4 NOT NULL DEFAULT 0,
    deleted_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by              int4 NOT NULL DEFAULT 0,
    CONSTRAINT bumd_pkey PRIMARY KEY (id_bumd)
);

-- Add foreign key constraints
ALTER TABLE trn_bumd 
ADD CONSTRAINT fk_bumd_id_bentuk_hukum 
FOREIGN KEY (id_bentuk_hukum) REFERENCES m_bentuk_badan_hukum(id_bbh);

ALTER TABLE trn_bumd 
ADD CONSTRAINT fk_bumd_id_bentuk_usaha 
FOREIGN KEY (id_bentuk_usaha) REFERENCES m_bentuk_usaha(id_bu);

-- Add indexes for better performance
CREATE INDEX idx_bumd_id_daerah 
    ON trn_bumd(id_daerah);

CREATE INDEX idx_bumd_id_bentuk_hukum 
    ON trn_bumd(id_bentuk_hukum);

CREATE INDEX idx_bumd_id_bentuk_usaha 
    ON trn_bumd(id_bentuk_usaha);
CREATE INDEX idx_bumd_id_induk_perusahaan 
    ON trn_bumd(id_induk_perusahaan);

CREATE INDEX idx_bumd_penerapan_spi 
    ON trn_bumd(penerapan_spi_bumd);

CREATE INDEX idx_bumd_nama 
    ON trn_bumd(nama_bumd);

CREATE INDEX idx_bumd_deleted_at 
    ON trn_bumd(deleted_at);