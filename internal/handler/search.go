package handler

import (
	"net/http"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
	"github.com/brolatief/fiqih-tahawulat/internal/repository"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "missing_query", "query parameter 'q' is required")
		return
	}

	page := parseIntDefault(r.URL.Query().Get("page"), 1)
	limit := parseIntDefault(r.URL.Query().Get("limit"), 10)
	limit = clampInt(limit, 1, 100)

	filter := repository.SearchFilter{
		Query:     query,
		Kategori:  r.URL.Query().Get("kategori"),
		Kitab:     r.URL.Query().Get("kitab"),
		Penulis:   r.URL.Query().Get("penulis"),
		Page:      page,
		Limit:     limit,
		Fuzzy:     r.URL.Query().Get("fuzzy") == "true",
		Highlight: r.URL.Query().Get("highlight") != "false",
	}

	data, total, err := h.svc.Search(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[[]model.SearchResult]{
		Data: data,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages(total, limit),
		},
	})
}
