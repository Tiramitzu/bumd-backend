CREATE TABLE trn_akta_notaris (
    id_akta_notaris             uuid NOT NULL,
    id_bumd                     uuid NOT NULL,
    nomor_akta_notaris          varchar(100) NOT NULL DEFAULT ''::character varying,
    notaris_akta_notaris        varchar(100) NOT NULL DEFAULT ''::character varying,
    tanggal_akta_notaris        date NOT NULL,
    keterangan_akta_notaris     text NOT NULL DEFAULT '',
    file_akta_notaris           text NOT NULL DEFAULT '',
    created_at                  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                  int4 NOT NULL DEFAULT 0,
    updated_at                  timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                  int4 NOT NULL DEFAULT 0,
    deleted_at                  timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                  int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_akta_notaris_pkey PRIMARY KEY (id_akta_notaris)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_akta_notaris 
ADD CONSTRAINT fk_trn_akta_notaris_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_akta_notaris_id_bumd 
    ON trn_akta_notaris(id_bumd);

CREATE INDEX idx_trn_akta_notaris_nomor
    ON trn_akta_notaris(nomor_akta_notaris);

CREATE INDEX idx_dkmn_akta_notaris_tanggal 
    ON trn_akta_notaris(tanggal_akta_notaris);

CREATE INDEX idx_trn_akta_notaris_deleted_at 
    ON trn_akta_notaris(deleted_at);