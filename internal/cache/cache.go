package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type NoopCache struct{}

func (NoopCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}
func (NoopCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	return nil
}
func (NoopCache) Del(ctx context.Context, keys ...string) error {
	return nil
}

func KeyListKitab(judul, kategori string, page, limit int) string {
	return "ks:kitab:list:" + judul + ":" + kategori + ":" + itoa(page) + ":" + itoa(limit)
}

func KeyGetKitab(id int) string {
	return "ks:kitab:" + itoa(id)
}

func KeyKonten(kitabID, page, limit int) string {
	return "ks:konten:" + itoa(kitabID) + ":" + itoa(page) + ":" + itoa(limit)
}

func KeySearch(query, kategori, kitab, penulis string, page, limit int) string {
	return "ks:search:" + query + ":" + kategori + ":" + kitab + ":" + penulis + ":" + itoa(page) + ":" + itoa(limit)
}

func KeyKategori() string {
	return "ks:kategori"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
