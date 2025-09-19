CREATE TABLE trn_pengurus (
    id_pengurus                     uuid NOT NULL,
    id_bumd                         uuid NOT NULL,
    jabatan_struktur_pengurus       int4 NOT NULL DEFAULT 0,
    nama_pengurus                   varchar(255) NOT NULL DEFAULT '',
    nik_pengurus                    varchar(255) NOT NULL DEFAULT '',
    alamat_pengurus                 varchar(255) NOT NULL DEFAULT '',
    deskripsi_jabatan_pengurus      text NOT NULL DEFAULT '',
    pendidikan_akhir_pengurus       uuid NOT NULL,
    tanggal_mulai_jabatan_pengurus  date NOT NULL DEFAULT '0001-01-01',
    tanggal_akhir_jabatan_pengurus  date NOT NULL DEFAULT '0001-01-01',
    file_pengurus                   text NOT NULL DEFAULT '',
    is_active                       int4 NOT NULL DEFAULT 1,
    created_at                      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                      int4 NOT NULL DEFAULT 0,
    updated_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                      int4 NOT NULL DEFAULT 0,
    deleted_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                      int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_pengurus_pkey PRIMARY KEY (id_pengurus)
);

-- Add foreign key constraint to bumd table
ALTER TABLE trn_pengurus ADD CONSTRAINT trn_pengurus_id_bumd_fk FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);
ALTER TABLE trn_pengurus ADD CONSTRAINT trn_pengurus_pendidikan_akhir_fk FOREIGN KEY (pendidikan_akhir_pengurus) REFERENCES m_pendidikan(id_pendidikan);

-- Create index
CREATE INDEX idx_trn_pengurus_id_bumd ON trn_pengurus(id_bumd);
CREATE INDEX idx_trn_pengurus_is_active ON trn_pengurus(is_active);
CREATE INDEX idx_trn_pengurus_jabatan_struktur ON trn_pengurus(jabatan_struktur_pengurus);
CREATE INDEX idx_trn_pengurus_nama_pengurus ON trn_pengurus(nama_pengurus);
CREATE INDEX idx_trn_pengurus_nik ON trn_pengurus(nik_pengurus);