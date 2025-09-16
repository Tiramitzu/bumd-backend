CREATE TABLE trn_peraturan (
    id_peraturan                    uuid NOT NULL,
    id_bumd                         uuid NOT NULL,
    nomor_peraturan                 varchar(100) NOT NULL DEFAULT ''::character varying,
    jenis_peraturan                 uuid NOT NULL,
    tanggal_berlaku_peraturan       date NOT NULL DEFAULT '0001-01-01',
    keterangan_peraturan            text NOT NULL DEFAULT '',
    file_peraturan                  text NOT NULL DEFAULT '',
    created_at                      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                      int4 NOT NULL DEFAULT 0,
    updated_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                      int4 NOT NULL DEFAULT 0,
    deleted_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                      int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_peraturan_pkey PRIMARY KEY (id_peraturan)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_peraturan 
ADD CONSTRAINT fk_trn_peraturan_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_peraturan_id_bumd 
    ON trn_peraturan(id_bumd);

CREATE INDEX idx_trn_peraturan_nomor 
    ON trn_peraturan(nomor_peraturan);

CREATE INDEX idx_trn_peraturan_masa_berlaku 
    ON trn_peraturan(tanggal_berlaku_peraturan);