package service

import (
	"context"

	"github.com/brolatief/fiqih-tahawulat/internal/cache"
	"github.com/brolatief/fiqih-tahawulat/internal/model"
	"github.com/brolatief/fiqih-tahawulat/internal/repository"
)

type Service interface {
	ListKitab(ctx context.Context, filter repository.KitabFilter) ([]model.DaftarKitab, int, error)
	GetKitab(ctx context.Context, id int) (*model.DaftarKitab, error)
	GetKonten(ctx context.Context, kitabID, page, limit int) ([]model.KontenKitab, int, error)
	Search(ctx context.Context, filter repository.SearchFilter) ([]model.SearchResult, int, error)
	SearchDirect(ctx context.Context, filter repository.SearchFilter) ([]model.SearchResult, int, error)
	RecommendKeywords(ctx context.Context, seed string, limit int) ([]model.KeywordRecommendation, error)
	ListKategori(ctx context.Context) ([]string, error)
}

type KitabService struct {
	repo repository.KitabRepository
}

func NewService(repo repository.KitabRepository, c cache.Cache) *KitabService {
	return &KitabService{repo: repo}
}

func (s *KitabService) ListKitab(ctx context.Context, filter repository.KitabFilter) ([]model.DaftarKitab, int, error) {
	return s.repo.ListKitab(ctx, filter)
}

func (s *KitabService) GetKitab(ctx context.Context, id int) (*model.DaftarKitab, error) {
	return s.repo.GetKitabByID(ctx, id)
}

func (s *KitabService) GetKonten(ctx context.Context, kitabID, page, limit int) ([]model.KontenKitab, int, error) {
	return s.repo.GetKontenByKitabID(ctx, kitabID, page, limit)
}

func (s *KitabService) Search(ctx context.Context, filter repository.SearchFilter) ([]model.SearchResult, int, error) {
	return s.repo.SearchKonten(ctx, filter)
}

func (s *KitabService) SearchDirect(ctx context.Context, filter repository.SearchFilter) ([]model.SearchResult, int, error) {
	return s.repo.SearchKonten(ctx, filter)
}

func (s *KitabService) RecommendKeywords(ctx context.Context, seed string, limit int) ([]model.KeywordRecommendation, error) {
	return s.repo.RecommendKeywords(ctx, seed, limit)
}

func (s *KitabService) ListKategori(ctx context.Context) ([]string, error) {
	return s.repo.ListKategori(ctx)
}

