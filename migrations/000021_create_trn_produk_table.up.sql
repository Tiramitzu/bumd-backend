CREATE TABLE trn_produk (
    id_produk               uuid NOT NULL,
    id_bumd                 uuid NOT NULL,
    nama_produk             varchar(255) NOT NULL DEFAULT '',
    deskripsi_produk        text NOT NULL DEFAULT '',
    foto_produk             text NOT NULL DEFAULT '',
    is_show                 int4 NOT NULL DEFAULT 0,
    created_at              timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by              int4 NOT NULL DEFAULT 0,
    updated_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    updated_by              int4 NOT NULL DEFAULT 0,
    deleted_at              timestamp NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    deleted_by              int4 NOT NULL DEFAULT 0,
    CONSTRAINT trn_produk_pkey PRIMARY KEY (id_produk)
);

ALTER TABLE trn_produk ADD CONSTRAINT trn_produk_id_bumd_fk FOREIGN KEY (id_bumd) REFERENCES trn_bumd(id_bumd);

CREATE INDEX idx_trn_produk_id_bumd ON trn_produk(id_bumd);
CREATE INDEX idx_trn_produk_nama_produk ON trn_produk(nama_produk);
CREATE INDEX idx_trn_produk_deskripsi ON trn_produk(deskripsi_produk);
CREATE INDEX idx_trn_produk_foto_produk ON trn_produk(foto_produk);
CREATE INDEX idx_trn_produk_is_show ON trn_produk(is_show);