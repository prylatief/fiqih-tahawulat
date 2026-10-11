package repository

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
)

const (
	listKitabSQL   = `SELECT id, judul, kategori, COALESCE(path_orig,''), COALESCE(penulis,''), COALESCE(publisher,'') FROM daftar_kitab WHERE (CAST($1 AS TEXT) IS NULL OR judul ILIKE '%'||$1||'%') AND (CAST($2 AS TEXT) IS NULL OR kategori = $2) ORDER BY id LIMIT $3 OFFSET $4`
	countKitabSQL  = `SELECT COUNT(*) FROM daftar_kitab WHERE (CAST($1 AS TEXT) IS NULL OR judul ILIKE '%'||$1||'%') AND (CAST($2 AS TEXT) IS NULL OR kategori = $2)`
	getKitabByIDSQL = `SELECT id, judul, kategori, COALESCE(path_orig,''), COALESCE(penulis,''), COALESCE(publisher,'') FROM daftar_kitab WHERE id = $1`

	listKontenSQL  = `SELECT id, kitab_id, nomor_bagian, isi_teks FROM konten_kitab WHERE kitab_id = $1 ORDER BY nomor_bagian LIMIT $2 OFFSET $3`
	countKontenSQL = `SELECT COUNT(*) FROM konten_kitab WHERE kitab_id = $1`

	searchKontenBaseSQL  = `SELECT dk.id, dk.judul, dk.kategori, COALESCE(dk.penulis,''), COALESCE(dk.publisher,''), kk.id, kk.nomor_bagian, LEFT(kk.isi_teks, 50) AS preview, ts_rank(kk.search_vector, websearch_to_tsquery('arabic', $1)) AS rank FROM konten_kitab kk JOIN daftar_kitab dk ON dk.id = kk.kitab_id WHERE kk.search_vector @@ websearch_to_tsquery('arabic', $1)`
	countSearchBaseSQL   = `SELECT COUNT(*) FROM konten_kitab kk JOIN daftar_kitab dk ON dk.id = kk.kitab_id WHERE kk.search_vector @@ websearch_to_tsquery('arabic', $1)`

	listKategoriSQL        = `SELECT DISTINCT kategori FROM daftar_kitab WHERE kategori IS NOT NULL ORDER BY kategori`
	streamKontenChunkedSQL = `SELECT kk.id, kk.kitab_id, kk.nomor_bagian, kk.isi_teks, dk.judul, dk.kategori, COALESCE(dk.penulis,''), COALESCE(dk.publisher,'') FROM konten_kitab kk JOIN daftar_kitab dk ON dk.id = kk.kitab_id WHERE kk.kitab_id = $1 AND kk.id > $2 ORDER BY kk.id LIMIT $3`
	listKitabIDsSQL        = `SELECT DISTINCT kitab_id FROM konten_kitab ORDER BY kitab_id`
)
