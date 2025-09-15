CREATE TABLE trn_kinerja (
    id_kinerja                      uuid NOT NULL,
    id_bumd                         uuid NOT NULL,
    tahun_kinerja                   int4 NOT NULL DEFAULT 0,
    ebit_kinerja                    float8 NOT NULL DEFAULT 0,
    ebitda_kinerja                  float8 NOT NULL DEFAULT 0,
    modal_sendiri_kinerja           float8 NOT NULL DEFAULT 0,
    penyusutan_kinerja              float8 NOT NULL DEFAULT 0,
    capital_employed_kinerja        float8 NOT NULL DEFAULT 0,
    total_aset_awal_kinerja         float8 NOT NULL DEFAULT 0,
    total_aset_akhir_kinerja        float8 NOT NULL DEFAULT 0,
    kas_kinerja                     float8 NOT NULL DEFAULT 0,
    setara_kas_kinerja              float8 NOT NULL DEFAULT 0,
    kewajiban_lancar_kinerja        float8 NOT NULL DEFAULT 0,
    harta_lancar_kinerja            float8 NOT NULL DEFAULT 0,
    penjualan_bersih_kinerja        float8 NOT NULL DEFAULT 0,
    piutang_dagang_kinerja          float8 NOT NULL DEFAULT 0,
    harga_pokok_penjualan_kinerja   float8 NOT NULL DEFAULT 0,
    persediaan_kinerja              float8 NOT NULL DEFAULT 0,
    aktiva_tetap_kinerja            float8 NOT NULL DEFAULT 0,
    akumulasi_depresiasi_kinerja    float8 NOT NULL DEFAULT 0,
    kredit_bermasalah_kinerja       float8 NOT NULL DEFAULT 0,
    total_kredit_kinerja            float8 NOT NULL DEFAULT 0,
    created_at                      timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                      int4 NOT NULL DEFAULT 0,
    updated_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by                      int4 NOT NULL DEFAULT 0,
    deleted_at                      timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by                      int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_kinerja_pkey PRIMARY KEY (id_kinerja)
);

ALTER TABLE trn_kinerja ADD CONSTRAINT trn_kinerja_id_bumd_fk FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

CREATE INDEX idx_trn_kinerja_id_bumd ON trn_kinerja(id_bumd);
CREATE INDEX idx_trn_kinerja_tahun ON trn_kinerja(tahun_kinerja);