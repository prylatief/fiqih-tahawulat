package handler

import (
	"net/http"
	"strconv"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
	"github.com/brolatief/fiqih-tahawulat/internal/repository"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) ListKitab(w http.ResponseWriter, r *http.Request) {
	page := parseIntDefault(r.URL.Query().Get("page"), 1)
	limit := parseIntDefault(r.URL.Query().Get("limit"), 20)
	limit = clampInt(limit, 1, 100)

	filter := repository.KitabFilter{
		Judul:    r.URL.Query().Get("judul"),
		Kategori: r.URL.Query().Get("kategori"),
		Page:     page,
		Limit:    limit,
	}

	data, total, err := h.svc.ListKitab(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[[]model.DaftarKitab]{
		Data: data,
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages(total, limit),
		},
	})
}

func (h *Handler) GetKitab(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "invalid kitab ID")
		return
	}

	kitab, err := h.svc.GetKitab(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "kitab not found")
		return
	}

	writeJSON(w, http.StatusOK, model.SingleResponse[*model.DaftarKitab]{Data: kitab})
}
