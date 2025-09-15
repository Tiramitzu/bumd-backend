CREATE TABLE trn_pegawai (
    id_pegawai              uuid NOT NULL,
    id_bumd                 uuid NOT NULL,
    tahun_pegawai           int4 NOT NULL DEFAULT 0,
    status_pegawai          int4 NOT NULL DEFAULT 0,
    pendidikan_pegawai      uuid NOT NULL,
    jumlah_pegawai          int4 NOT NULL DEFAULT 0,
    created_at              timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              int4 NOT NULL DEFAULT 0,
    updated_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by              int4 NOT NULL DEFAULT 0,
    deleted_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by              int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_pegawai_pkey PRIMARY KEY (id_pegawai)
);

ALTER TABLE trn_pegawai ADD CONSTRAINT trn_pegawai_id_bumd_fk FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);
ALTER TABLE trn_pegawai ADD CONSTRAINT trn_pegawai_pendidikan_fk FOREIGN KEY (pendidikan_pegawai) REFERENCES m_pendidikan(id_pendidikan);

CREATE INDEX idx_trn_pegawai_id_bumd ON trn_pegawai(id_bumd);
CREATE INDEX idx_trn_pegawai_tahun ON trn_pegawai(tahun_pegawai);
CREATE INDEX idx_trn_pegawai_status_pegawai ON trn_pegawai(status_pegawai);
CREATE INDEX idx_trn_pegawai_pendidikan ON trn_pegawai(pendidikan_pegawai);
CREATE INDEX idx_trn_pegawai_jumlah_pegawai ON trn_pegawai(jumlah_pegawai);