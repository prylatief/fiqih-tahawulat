package model

type DaftarKitab struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul"`
	Kategori  string `json:"kategori"`
	PathOrig  string `json:"path_orig"`
	Penulis   string `json:"penulis"`
	Publisher string `json:"publisher"`
}

type KontenKitab struct {
	ID          int    `json:"id"`
	KitabID     int    `json:"kitab_id"`
	NomorBagian int    `json:"nomor_bagian"`
	IsiTeks     string `json:"isi_teks"`
}

type SearchResult struct {
	KitabID     int     `json:"kitab_id"`
	Judul       string  `json:"judul"`
	Kategori    string  `json:"kategori"`
	Penulis     string  `json:"penulis"`
	Publisher   string  `json:"publisher"`
	SectionID   int     `json:"section_id"`
	NomorBagian int     `json:"nomor_bagian"`
	IsiTeks     string  `json:"isi_teks"`
	Rank        float64 `json:"rank"`
	Highlight   string  `json:"highlight,omitempty"`
}

type KontenResponse struct {
	KitabID  int          `json:"kitab_id"`
	Judul    string       `json:"judul"`
	Penulis  string       `json:"penulis"`
	Sections []KontenKitab `json:"sections"`
}

type PaginatedResponse[T any] struct {
	Data       T          `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type KategoriResponse struct {
	Data  []string `json:"data"`
	Total int      `json:"total"`
}

type SingleResponse[T any] struct {
	Data T `json:"data"`
}

type KeywordRecommendation struct {
	Keyword  string `json:"keyword"`
	DocCount int64  `json:"doc_count"`
	Score    float64 `json:"score"`
}

type KeywordRecommendationResponse struct {
	Seed     string                 `json:"seed"`
	Keywords []KeywordRecommendation `json:"keywords"`
}
