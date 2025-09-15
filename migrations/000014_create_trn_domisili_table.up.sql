CREATE TABLE trn_domisili (
    id_domisili         uuid NOT NULL,
    id_bumd             uuid NOT NULL,
    nomor_domisili      varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_domisili varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_domisili    date NOT NULL,
    kualifikasi_domisili int4 NOT NULL DEFAULT 0,
    klasifikasi_domisili varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_domisili date NULL,
    file_domisili       text NOT NULL DEFAULT '',
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          int4 NOT NULL DEFAULT 0,
    updated_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by          int4 NOT NULL DEFAULT 0,
    deleted_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by          int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_domisili_pkey PRIMARY KEY (id_domisili)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_domisili 
ADD CONSTRAINT fk_trn_domisili_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_domisili_id_bumd 
    ON trn_domisili(id_bumd);

CREATE INDEX idx_trn_domisili_nomor 
    ON trn_domisili(nomor_domisili);

CREATE INDEX idx_trn_domisili_tanggal 
    ON trn_domisili(tanggal_domisili);

CREATE INDEX idx_trn_domisili_masa_berlaku 
    ON trn_domisili(masa_berlaku_domisili);