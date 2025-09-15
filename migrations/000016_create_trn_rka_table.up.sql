CREATE TABLE trn_rka (
    id_rka              uuid NOT NULL,
    id_bumd             uuid NOT NULL,
    nomor_rka           varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_rka varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_rka         date NOT NULL,
    kualifikasi_rka     int4 NOT NULL DEFAULT 0,
    klasifikasi_rka     varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_rka    date NULL,
    file_rka            text NOT NULL DEFAULT '',
    created_at          timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by          int4 NOT NULL DEFAULT 0,
    updated_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by          int4 NOT NULL DEFAULT 0,
    deleted_at          timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by          int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_rka_pkey PRIMARY KEY (id_rka)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_rka 
ADD CONSTRAINT fk_trn_rka_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_rka_id_bumd 
    ON trn_rka(id_bumd);

CREATE INDEX idx_trn_rka_nomor 
    ON trn_rka(nomor_rka);

CREATE INDEX idx_trn_rka_tanggal 
    ON trn_rka(tanggal_rka);

CREATE INDEX idx_trn_rka_masa_berlaku 
    ON trn_rka(masa_berlaku_rka);