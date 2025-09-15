CREATE TABLE trn_rencana_bisnis (
    id_rencana_bisnis   uuid NOT NULL,
    id_bumd             uuid NOT NULL,
    nomor_rencana_bisnis varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_rencana_bisnis varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_rencana_bisnis date NOT NULL,
    kualifikasi_rencana_bisnis int4 NOT NULL DEFAULT 0,
    klasifikasi_rencana_bisnis varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_rencana_bisnis date NULL,
    file_rencana_bisnis  text NOT NULL DEFAULT '',
    created_at           timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by           int4 NOT NULL DEFAULT 0,
    updated_at           timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by           int4 NOT NULL DEFAULT 0,
    deleted_at           timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by           int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_rencana_bisnis_pkey PRIMARY KEY (id_rencana_bisnis)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_rencana_bisnis 
ADD CONSTRAINT fk_trn_rencana_bisnis_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_rencana_bisnis_id_bumd 
    ON trn_rencana_bisnis(id_bumd);

CREATE INDEX idx_trn_rencana_bisnis_nomor 
    ON trn_rencana_bisnis(nomor_rencana_bisnis);

CREATE INDEX idx_trn_rencana_bisnis_tanggal 
    ON trn_rencana_bisnis(tanggal_rencana_bisnis);

CREATE INDEX idx_trn_rencana_bisnis_masa_berlaku 
    ON trn_rencana_bisnis(masa_berlaku_rencana_bisnis);