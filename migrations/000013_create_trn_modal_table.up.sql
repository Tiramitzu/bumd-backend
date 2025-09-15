CREATE TABLE trn_modal (
    id_modal           uuid NOT NULL,
    id_bumd            uuid NOT NULL,
    id_prov            int4 NOT NULL DEFAULT 0,
    id_kab             int4 NOT NULL DEFAULT 0,
    pemegang_modal     varchar(100) NOT NULL DEFAULT ''::character varying,
    no_ba_modal        varchar(100) NOT NULL DEFAULT ''::character varying,
    tanggal_modal      date NOT NULL,
    jumlah_modal       float8 NOT NULL DEFAULT 0,
    keterangan_modal   text NOT NULL DEFAULT ''::text,
    created_at         timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by         int4 NOT NULL DEFAULT 0,
    updated_at         timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by         int4 NOT NULL DEFAULT 0,
    deleted_at         timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by         int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_modal_pkey PRIMARY KEY (id_modal)
);

-- Add foreign key constraint to BUMD table
ALTER TABLE trn_modal
    ADD CONSTRAINT fk_trn_modal_id_bumd
    FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add indexes for better performance
CREATE INDEX idx_trn_modal_id_bumd
    ON trn_modal(id_bumd);


