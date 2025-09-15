CREATE TABLE trn_perda_pendirian (
    id_perda_pendirian                uuid NOT NULL,
    id_bumd                           uuid NOT NULL,
    nomor_perda_pendirian             varchar(100) NOT NULL DEFAULT ''::character varying,
    tanggal_perda_pendirian           date NOT NULL,
    modal_dasar_perda_pendirian       float8 NOT NULL DEFAULT 0,
    keterangan_perda_pendirian        text NOT NULL DEFAULT '',
    file_perda_pendirian              text NOT NULL DEFAULT '',
    created_at                        timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                        int4 NOT NULL DEFAULT 0,
    updated_at                        timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                        int4 NOT NULL DEFAULT 0,
    deleted_at                        timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                        int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_perda_pendirian_pkey PRIMARY KEY (id_perda_pendirian)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_perda_pendirian 
ADD CONSTRAINT fk_trn_perda_pendirian_id_bumd
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_perda_pendirian_id_bumd_pendirian 
    ON trn_perda_pendirian(id_bumd);

CREATE INDEX idx_trn_perda_pendirian_nomor_perda_pendirian
    ON trn_perda_pendirian(nomor_perda_pendirian);

CREATE INDEX idx_trn_perda_pendirian_tanggal_perda_pendirian 
    ON trn_perda_pendirian(tanggal_perda_pendirian);

CREATE INDEX idx_trn_perda_pendirian_deleted_at_pendirian 
    ON trn_perda_pendirian(deleted_at);