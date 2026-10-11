package repository

import (
	"context"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
)

type KitabFilter struct {
	Judul    string
	Kategori string
	Page     int
	Limit    int
}

type SearchFilter struct {
	Query     string
	Kategori  string
	Kitab     string
	Penulis   string
	Page      int
	Limit     int
	Fuzzy     bool
	Highlight bool
}

type KitabRepository interface {
	ListKitab(ctx context.Context, filter KitabFilter) ([]model.DaftarKitab, int, error)
	GetKitabByID(ctx context.Context, id int) (*model.DaftarKitab, error)
	GetKontenByKitabID(ctx context.Context, kitabID, page, limit int) ([]model.KontenKitab, int, error)
	SearchKonten(ctx context.Context, filter SearchFilter) ([]model.SearchResult, int, error)
	RecommendKeywords(ctx context.Context, seed string, limit int) ([]model.KeywordRecommendation, error)
	ListKategori(ctx context.Context) ([]string, error)
	StreamKontenChunked(ctx context.Context, kitabID, chunkSize int, fn func(model.SearchResult) error) error
	ListKitabIDs(ctx context.Context) ([]int, error)
}
