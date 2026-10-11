# Fiqih Tahawulat API — Habib Abu Bakar

REST API + MCP server untuk koleksi digital **38 kitab Habib Abu Bakar bin Ali Al-Aththas al-Husayni** (5.963 halaman), dengan full-text search Bahasa Arab.

## Stack

- **Go 1.22+** — `chi/v5` router
- **PostgreSQL** — `pgx/v5` + tsvector Arabic FTS
- **MCP** — Model Context Protocol server untuk integrasi AI agent

## Quick Start

```bash
# Clone
gh repo clone prylatief/fiqih-tahawulat
cd fiqih-tahawulat

# Restore database
gunzip < fiqih_tahawulat.sql.gz | sudo -u postgres psql

# Run API
KUTUB_DATABASE_URL="postgres://postgres:postgres@localhost:5432/fiqih_tahawulat?sslmode=disable" \
    go run ./cmd/api
```

Server akan berjalan di `http://localhost:3000`.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/kitab` | List semua kitab (pagination + filter `?judul=` & `?kategori=`) |
| `GET` | `/v1/kitab/{id}` | Detail metadata kitab |
| `GET` | `/v1/kitab/{id}/konten` | Konten per halaman (`?page=&limit=`) |
| `GET` | `/v1/kitab/{id}/konten/{hal}` | Halaman spesifik |
| `GET` | `/v1/search?q=` | Full-text search Arabic (ranked) |
| `GET` | `/v1/kategori` | Daftar kategori |
| `GET` | `/health` | Health check |

## Search

Search menggunakan PostgreSQL `websearch_to_tsquery('arabic', ...)` dengan ranking `ts_rank`:

```bash
# Cari kata "تحولات"
curl 'http://localhost:3000/v1/search?q=%D8%AA%D8%AD%D9%88%D9%84%D8%A7%D8%AA&limit=5'

# Filter kategori
curl 'http://localhost:3000/v1/search?q=%D8%AA%D8%AD%D9%88%D9%84&kategori=Fiqih%20Tahawulat'

# Filter kitab
curl 'http://localhost:3000/v1/search?q=%D8%AA%D8%AD%D9%88%D9%84&kitab=%D8%A7%D9%84%D8%AA%D9%84%D9%8A%D8%AF'
```

## MCP Server

Untuk integrasi dengan AI agent (Claude Desktop, opencode, dll):

```bash
# Stdio transport (default)
go run ./cmd/mcp

# Atau build dulu
make build-mcp
./bin/fiqih-tahawulat-mcp
```

### MCP Tools

| Tool | Description |
|------|-------------|
| `list_kitab` | Browse kitab catalog |
| `get_kitab` | Get kitab metadata by ID |
| `get_konten` | Read content sections |
| `search_kitab` | Full-text search across all kitab |
| `list_kategori` | List categories |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KUTUB_DATABASE_URL` | `postgres://postgres:@localhost:5432/fiqih_tahawulat?sslmode=disable` | PostgreSQL DSN |
| `KUTUB_ADDR` | `:3000` | Listen address |
| `KUTUB_ALLOWED_ORIGINS` | `http://localhost:3000` | CORS origins |

## Development

```bash
make run       # Run API server
make build     # Build binary
make build-mcp # Build MCP server
make test      # Run tests
```

## Menambah Kitab Baru

```bash
./update.sh <google-drive-folder-url>
```

Atau manual: taruh PDF di folder, jalankan `python3 extract.py && python3 import_db5.py`.

## Dataset

- **38 kitab** Habib Abu Bakar
- **5.963 halaman** konten
- Format: Arabic text + PostgreSQL tsvector
- 1 kitab hasil OCR (بين يدي الدجال)

## License

Hak cipta isi kitab milik pengarang. Source code API: MIT.
