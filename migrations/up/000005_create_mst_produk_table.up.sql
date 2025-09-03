CREATE TABLE mst_produk (
    id              BIGSERIAL       PRIMARY KEY,
    id_bumd         INT             NOT NULL DEFAULT 0,
    nama            VARCHAR(250)    NOT NULL DEFAULT '',
    deskripsi       TEXT            NOT NULL DEFAULT '',
    gambar          TEXT            NOT NULL DEFAULT '',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT             NOT NULL DEFAULT 0,
    updated_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by      INT             NOT NULL DEFAULT 0,
    deleted_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by      INT             NOT NULL DEFAULT 0
);
