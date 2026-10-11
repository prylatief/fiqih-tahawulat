package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

type PostgresRepo struct {
	pool     *pgxpool.Pool
	logLevel string
}

func NewPostgres(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool, logLevel: "1"}
}

func (r *PostgresRepo) ListKitab(ctx context.Context, filter KitabFilter) ([]model.DaftarKitab, int, error) {
	g, gctx := errgroup.WithContext(ctx)

	var kitab []model.DaftarKitab
	g.Go(func() error {
		offset := (filter.Page - 1) * filter.Limit
		rows, err := r.pool.Query(gctx, listKitabSQL, nullable(filter.Judul), nullable(filter.Kategori), filter.Limit, offset)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var k model.DaftarKitab
			if err := rows.Scan(&k.ID, &k.Judul, &k.Kategori, &k.PathOrig, &k.Penulis, &k.Publisher); err != nil {
				return err
			}
			kitab = append(kitab, k)
		}
		return rows.Err()
	})

	var total int
	g.Go(func() error {
		return r.pool.QueryRow(gctx, countKitabSQL, nullable(filter.Judul), nullable(filter.Kategori)).Scan(&total)
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	if kitab == nil {
		kitab = []model.DaftarKitab{}
	}
	return kitab, total, nil
}

func (r *PostgresRepo) GetKitabByID(ctx context.Context, id int) (*model.DaftarKitab, error) {
	var k model.DaftarKitab
	err := r.pool.QueryRow(ctx, getKitabByIDSQL, id).Scan(&k.ID, &k.Judul, &k.Kategori, &k.PathOrig, &k.Penulis, &k.Publisher)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *PostgresRepo) GetKontenByKitabID(ctx context.Context, kitabID, page, limit int) ([]model.KontenKitab, int, error) {
	g, gctx := errgroup.WithContext(ctx)

	var konten []model.KontenKitab
	g.Go(func() error {
		offset := (page - 1) * limit
		rows, err := r.pool.Query(gctx, listKontenSQL, kitabID, limit, offset)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var k model.KontenKitab
			if err := rows.Scan(&k.ID, &k.KitabID, &k.NomorBagian, &k.IsiTeks); err != nil {
				return err
			}
			konten = append(konten, k)
		}
		return rows.Err()
	})

	var total int
	g.Go(func() error {
		return r.pool.QueryRow(gctx, countKontenSQL, kitabID).Scan(&total)
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	if konten == nil {
		konten = []model.KontenKitab{}
	}
	return konten, total, nil
}

func (r *PostgresRepo) SearchKonten(ctx context.Context, filter SearchFilter) ([]model.SearchResult, int, error) {
	g, gctx := errgroup.WithContext(ctx)

	var results []model.SearchResult
	g.Go(func() error {
		query := buildSearchQuery(filter)
		rows, err := r.pool.Query(gctx, query, filter.Query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var s model.SearchResult
			if err := rows.Scan(&s.KitabID, &s.Judul, &s.Kategori, &s.Penulis, &s.Publisher, &s.SectionID, &s.NomorBagian, &s.IsiTeks, &s.Rank); err != nil {
				return err
			}
			results = append(results, s)
		}
		return rows.Err()
	})

	var total int
	g.Go(func() error {
		q := countSearchBaseSQL + buildSearchWhere(filter)
		return r.pool.QueryRow(gctx, q, filter.Query).Scan(&total)
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	if results == nil {
		results = []model.SearchResult{}
	}
	return results, total, nil
}

func buildSearchQuery(filter SearchFilter) string {
	q := searchKontenBaseSQL + buildSearchWhere(filter)
	offset := (filter.Page - 1) * filter.Limit
	q += fmt.Sprintf(" ORDER BY rank DESC LIMIT %d OFFSET %d", filter.Limit, offset)
	return q
}

func buildSearchWhere(filter SearchFilter) string {
	var clauses []string
	if filter.Kategori != "" {
		clauses = append(clauses, fmt.Sprintf(" AND dk.kategori = '%s'", strings.ReplaceAll(filter.Kategori, "'", "''")))
	}
	if filter.Kitab != "" {
		clauses = append(clauses, fmt.Sprintf(" AND dk.judul ILIKE '%%%s%%'", strings.ReplaceAll(filter.Kitab, "'", "''")))
	}
	if filter.Penulis != "" {
		clauses = append(clauses, fmt.Sprintf(" AND dk.penulis ILIKE '%%%s%%'", strings.ReplaceAll(filter.Penulis, "'", "''")))
	}
	return strings.Join(clauses, "")
}

func (r *PostgresRepo) RecommendKeywords(ctx context.Context, seed string, limit int) ([]model.KeywordRecommendation, error) {
	return nil, nil
}

func (r *PostgresRepo) ListKategori(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx, listKategoriSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kategories []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, err
		}
		kategories = append(kategories, k)
	}
	if kategories == nil {
		kategories = []string{}
	}
	return kategories, rows.Err()
}

func (r *PostgresRepo) StreamKontenChunked(ctx context.Context, kitabID, chunkSize int, fn func(model.SearchResult) error) error {
	lastID := 0
	for {
		rows, err := r.pool.Query(ctx, streamKontenChunkedSQL, kitabID, lastID, chunkSize)
		if err != nil {
			return err
		}

		count := 0
		for rows.Next() {
			var s model.SearchResult
			if err := rows.Scan(&s.SectionID, &s.KitabID, &s.NomorBagian, &s.IsiTeks, &s.Judul, &s.Kategori, &s.Penulis, &s.Publisher); err != nil {
				rows.Close()
				return err
			}
			lastID = s.SectionID
			count++
			if err := fn(s); err != nil {
				rows.Close()
				return err
			}
		}
		rows.Close()

		if count == 0 {
			break
		}
	}
	return nil
}

func (r *PostgresRepo) ListKitabIDs(ctx context.Context) ([]int, error) {
	rows, err := r.pool.Query(ctx, listKitabIDsSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func nullable(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
