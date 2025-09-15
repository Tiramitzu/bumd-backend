CREATE TABLE trn_siup (
    id_siup               uuid NOT NULL,
    id_bumd               uuid NOT NULL,
    nomor_siup            varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_siup varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_siup          date NOT NULL,
    kualifikasi_siup      int4 NOT NULL DEFAULT 0,
    klasifikasi_siup      varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_siup     date NULL,
    file_siup             text NOT NULL DEFAULT '',
    created_at            timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by            int4 NOT NULL DEFAULT 0,
    updated_at            timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by            int4 NOT NULL DEFAULT 0,
    deleted_at            timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by            int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_siup_pkey PRIMARY KEY (id_siup)
);

-- Add foreign key constraint to BUMD table
ALTER TABLE trn_siup
    ADD CONSTRAINT fk_trn_siup_id_bumd
    FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_siup_id_bumd
    ON trn_siup(id_bumd);

CREATE INDEX idx_trn_siup_nomor
    ON trn_siup(nomor_siup);

CREATE INDEX idx_trn_siup_tanggal
    ON trn_siup(tanggal_siup);

CREATE INDEX idx_trn_siup_masa_berlaku
    ON trn_siup(masa_berlaku_siup);

CREATE INDEX idx_trn_siup_deleted_at
    ON trn_siup(deleted_at);


