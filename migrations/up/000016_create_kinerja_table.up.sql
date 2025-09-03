CREATE TABLE kinerja (
    id                              BIGSERIAL       PRIMARY KEY,
    id_bumd                         INT             NOT NULL,
    tahun                           INT             NOT NULL,
    laba_bersih_sebelum_pajak       DECIMAL(15,2)   NOT NULL DEFAULT 0,
    laba_bersih_setelah_pajak       DECIMAL(15,2)   NOT NULL DEFAULT 0,
    modal_sendiri                   DECIMAL(15,2)   NOT NULL DEFAULT 0,
    penyusutan                      DECIMAL(15,2)   NOT NULL DEFAULT 0,
    total_aset_kewajiban_lancar     DECIMAL(15,2)   NOT NULL DEFAULT 0,
    total_aset_awal                 DECIMAL(15,2)   NOT NULL DEFAULT 0,
    total_aset_akhir                DECIMAL(15,2)   NOT NULL DEFAULT 0,
    kas                             DECIMAL(15,2)   NOT NULL DEFAULT 0,
    setara_kas                      DECIMAL(15,2)   NOT NULL DEFAULT 0,
    kewajiban_lancar                DECIMAL(15,2)   NOT NULL DEFAULT 0,
    harta_lancar                    DECIMAL(15,2)   NOT NULL DEFAULT 0,
    penjualan_bersih                DECIMAL(15,2)   NOT NULL DEFAULT 0,
    rata_rata_piutang_dagang        DECIMAL(15,2)   NOT NULL DEFAULT 0,
    harga_pokok_penjualan           DECIMAL(15,2)   NOT NULL DEFAULT 0,
    rata_rata_persediaan            DECIMAL(15,2)   NOT NULL DEFAULT 0,
    aktiva_tetap                    DECIMAL(15,2)   NOT NULL DEFAULT 0,
    akumulasi_depresiasi            DECIMAL(15,2)   NOT NULL DEFAULT 0,
    kredit_bermasalah               DECIMAL(15,2)   NOT NULL DEFAULT 0,
    total_kredit                    DECIMAL(15,2)   NOT NULL DEFAULT 0,
    created_at                      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by                      INT             NOT NULL DEFAULT 0,
    updated_at                      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by                      INT             NOT NULL DEFAULT 0,
    deleted_at                      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by                      INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraint to bumd table
ALTER TABLE kinerja 
ADD CONSTRAINT fk_kinerja_id_bumd 
FOREIGN KEY (id_bumd) REFERENCES bumd(id);

-- Add indexes for better performance
CREATE INDEX idx_kinerja_id_bumd ON kinerja(id_bumd);
CREATE INDEX idx_kinerja_tahun ON kinerja(tahun);
CREATE INDEX idx_kinerja_deleted_at ON kinerja(deleted_at);
