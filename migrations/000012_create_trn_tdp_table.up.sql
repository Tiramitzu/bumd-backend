CREATE TABLE trn_tdp (
    id_tdp               uuid NOT NULL,
    id_bumd              uuid NOT NULL,
    nomor_tdp            varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_tdp varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_tdp          date NOT NULL,
    kualifikasi_tdp      int4 NOT NULL DEFAULT 0,
    klasifikasi_tdp      varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_tdp     date NULL,
    file_tdp             text NOT NULL DEFAULT '',
    created_at           timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by           int4 NOT NULL DEFAULT 0,
    updated_at           timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by           int4 NOT NULL DEFAULT 0,
    deleted_at           timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by           int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_tdp_pkey PRIMARY KEY (id_tdp)
);

-- Add foreign key constraint to BUMD table
ALTER TABLE trn_tdp
    ADD CONSTRAINT fk_trn_tdp_id_bumd
    FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_tdp_id_bumd
    ON trn_tdp(id_bumd);

CREATE INDEX idx_trn_tdp_nomor
    ON trn_tdp(nomor_tdp);

CREATE INDEX idx_trn_tdp_tanggal
    ON trn_tdp(tanggal_tdp);

CREATE INDEX idx_trn_tdp_masa_berlaku
    ON trn_tdp(masa_berlaku_tdp);

CREATE INDEX idx_trn_tdp_deleted_at
    ON trn_tdp(deleted_at);


