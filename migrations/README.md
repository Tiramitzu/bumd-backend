# Database Migrations

This directory contains database migrations organized in a clean folder structure.

## 📁 Folder Structure

```
migrations/
├── up/                    # Up migrations (CREATE tables)
│   ├── 000001_create_roles_table.up.sql
│   ├── 000002_create_mst_bentuk_badan_hukum_table.up.sql
│   ├── 000003_create_mst_bentuk_usaha_table.up.sql
│   ├── 000004_create_mst_jenis_dokumen_table.up.sql
│   ├── 000005_create_mst_produk_table.up.sql
│   ├── 000006_create_sys_config_table.up.sql
│   ├── 000007_create_users_table.up.sql
│   ├── 000008_create_bumd_table.up.sql
│   ├── 000009_create_mst_perda_table.up.sql
│   ├── 000010_create_mst_akta_notaris_table.up.sql
│   ├── 000011_create_nib_table.up.sql
│   ├── 000012_create_rencana_bisnis_table.up.sql
│   ├── 000013_create_rka_table.up.sql
│   ├── 000014_create_peraturan_table.up.sql
│   ├── 000015_create_pendidikan_table.up.sql
│   ├── 000016_create_kinerja_table.up.sql
│   └── 000017_create_laporan_keuangan_table.up.sql
├── down/                  # Down migrations (DROP tables)
│   ├── 000001_create_roles_table.down.sql
│   ├── 000002_create_mst_bentuk_badan_hukum_table.down.sql
│   ├── 000003_create_mst_bentuk_usaha_table.down.sql
│   ├── 000004_create_mst_jenis_dokumen_table.down.sql
│   ├── 000005_create_mst_produk_table.down.sql
│   ├── 000006_create_sys_config_table.down.sql
│   ├── 000007_create_users_table.down.sql
│   ├── 000008_create_bumd_table.down.sql
│   ├── 000009_create_mst_perda_table.down.sql
│   ├── 000010_create_mst_akta_notaris_table.down.sql
│   ├── 000011_create_nib_table.down.sql
│   ├── 000012_create_rencana_bisnis_table.down.sql
│   ├── 000013_create_rka_table.up.sql
│   ├── 000014_create_peraturan_table.down.sql
│   ├── 000015_create_pendidikan_table.down.sql
│   ├── 000016_create_kinerja_table.down.sql
│   └── 000017_create_laporan_keuangan_table.down.sql
└── README.md             # This file
```

## 🔄 Migration Dependencies

The migrations are ordered to respect foreign key dependencies:

1. **000001-000006**: Master/lookup tables (no dependencies)
2. **000007**: Users table (depends on roles)
3. **000008**: BUMD table (depends on master tables)
4. **000009-000017**: Business tables (depend on BUMD)

## 🚀 Usage

### Running Migrations
```bash
# Run all pending migrations
migrate -path migrations -database "your_connection_string" up

# Run specific number of migrations
migrate -path migrations -database "your_connection_string" step 5

# Check current version
migrate -path migrations -database "your_connection_string" version
```

### Rolling Back Migrations
```bash
# Rollback specific number of migrations
migrate -path migrations -database "your_connection_string" down 3

# Rollback to specific version
migrate -path migrations -database "your_connection_string" goto 5

# Rollback all migrations
migrate -path migrations -database "your_connection_string" down
```

## 📝 File Naming Convention

- **Up migrations**: `000XXX_description.up.sql`
- **Down migrations**: `000XXX_description.down.sql`
- **Version numbers**: Sequential 6-digit numbers (000001, 000002, etc.)

## ⚠️ Important Notes

- **Never delete migration files** that have already been run in production
- **Always test migrations** in development before running in production
- **Backup your database** before running migrations in production
- **Check dependencies** when adding new migrations

## 🛠️ Development

When adding new migrations:

1. Create both `.up.sql` and `.down.sql` files
2. Place them in the appropriate `up/` and `down/` folders
3. Use the next sequential version number
4. Ensure proper dependency order
5. Test both up and down migrations
