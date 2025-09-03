CREATE TABLE bumd (
    id              BIGSERIAL       PRIMARY KEY,
    id_daerah       INT             NOT NULL DEFAULT 0,
    id_bentuk_hukum INT             NOT NULL DEFAULT 0,
    id_bentuk_usaha INT             NOT NULL DEFAULT 0,
    nama            VARCHAR(250)    NOT NULL DEFAULT '',
    deskripsi       TEXT            NOT NULL DEFAULT '',
    alamat          TEXT            NOT NULL DEFAULT '',
    no_telp         VARCHAR(30)     NOT NULL DEFAULT '',
    email           VARCHAR(250)    NOT NULL DEFAULT '',
    website         VARCHAR(250)    NOT NULL DEFAULT '',
    created_at      TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      INT             NOT NULL DEFAULT 0,
    updated_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    updated_by      INT             NOT NULL DEFAULT 0,
    deleted_at      TIMESTAMP       NOT NULL DEFAULT '0001-01-01 00:00:00',
    deleted_by      INT             NOT NULL DEFAULT 0
);

-- Add foreign key constraints
ALTER TABLE bumd 
ADD CONSTRAINT fk_bumd_id_bentuk_hukum 
FOREIGN KEY (id_bentuk_hukum) REFERENCES mst_bentuk_badan_hukum(id);

ALTER TABLE bumd 
ADD CONSTRAINT fk_bumd_id_bentuk_usaha 
FOREIGN KEY (id_bentuk_usaha) REFERENCES mst_bentuk_usaha(id);

-- Add indexes for better performance
CREATE INDEX idx_bumd_id_daerah 
    ON bumd(id_daerah);

CREATE INDEX idx_bumd_id_bentuk_hukum 
    ON bumd(id_bentuk_hukum);

CREATE INDEX idx_bumd_id_bentuk_usaha 
    ON bumd(id_bentuk_usaha);

CREATE INDEX idx_bumd_nama 
    ON bumd(nama);

CREATE INDEX idx_bumd_deleted_at 
    ON bumd(deleted_at);
