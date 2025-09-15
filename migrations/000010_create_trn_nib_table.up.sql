CREATE TABLE trn_nib (
    id_nib                  uuid NOT NULL,
    id_bumd                 uuid NOT NULL,
    nomor_nib               varchar(100) NOT NULL DEFAULT ''::character varying,
    instansi_pemberi_nib    varchar(250) NOT NULL DEFAULT ''::character varying,
    tanggal_nib             date NOT NULL,
    kualifikasi_nib         int4 NOT NULL DEFAULT 0,
    klasifikasi_nib         varchar(100) NOT NULL DEFAULT ''::character varying,
    masa_berlaku_nib        date NULL,
    file_nib                text NOT NULL DEFAULT '',
    created_at              timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              int4 NOT NULL DEFAULT 0,
    updated_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by              int4 NOT NULL DEFAULT 0,
    deleted_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by              int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_nib_pkey PRIMARY KEY (id_nib)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_nib 
ADD CONSTRAINT fk_trn_nib_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_nib_id_bumd 
    ON trn_nib(id_bumd);

CREATE INDEX idx_trn_nib_nomor 
    ON trn_nib(nomor_nib);

CREATE INDEX idx_trn_nib_tanggal 
    ON trn_nib(tanggal_nib);

CREATE INDEX idx_trn_nib_masa_berlaku 
    ON trn_nib(masa_berlaku_nib);

CREATE INDEX idx_trn_nib_deleted_at 
    ON trn_nib(deleted_at);