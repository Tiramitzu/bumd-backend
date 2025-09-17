CREATE TABLE trn_laporan_keuangan (
    id_laporan_keuangan         uuid NOT NULL,
    id_bumd                     uuid NOT NULL,
    id_jenis_laporan            int8 NOT NULL,
    id_jenis_laporan_item       int8 NOT NULL,
    tahun_laporan_keuangan      int4 NOT NULL DEFAULT 0,
    jumlah_laporan_keuangan     float8 NOT NULL DEFAULT 0,
    file_laporan_keuangan       text NOT NULL DEFAULT ''::text,
    created_at                  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                  int4 NOT NULL DEFAULT 0,
    updated_at                  timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                  int4 NOT NULL DEFAULT 0,
    deleted_at                  timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                  int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_laporan_keuangan_pkey PRIMARY KEY (id_laporan_keuangan)
);

-- Add foreign key constraint to trn_bumd table
ALTER TABLE trn_laporan_keuangan 
ADD CONSTRAINT fk_trn_laporan_keuangan_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

-- Add foreign key constraint to m_jenis_laporan table
ALTER TABLE trn_laporan_keuangan 
ADD CONSTRAINT fk_trn_laporan_keuangan_id_jenis_laporan
FOREIGN KEY (id_jenis_laporan) REFERENCES m_jenis_laporan(id_jenis_laporan);
