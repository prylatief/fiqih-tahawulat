package handler

import (
	"net/http"
	"strconv"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) GetKonten(w http.ResponseWriter, r *http.Request) {
	kitabID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_id", "invalid kitab ID")
		return
	}

	page := 1
	limit := 10

	if halStr := chi.URLParam(r, "hal"); halStr != "" {
		page, err = strconv.Atoi(halStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_hal", "invalid page number")
			return
		}
		limit = 1
	} else {
		page = parseIntDefault(r.URL.Query().Get("page"), 1)
		limit = parseIntDefault(r.URL.Query().Get("limit"), 10)
		limit = clampInt(limit, 1, 100)
	}

	data, total, err := h.svc.GetKonten(r.Context(), kitabID, page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	kitab, _ := h.svc.GetKitab(r.Context(), kitabID)
	judul := ""
	penulis := ""
	if kitab != nil {
		judul = kitab.Judul
		penulis = kitab.Penulis
	}

	if lang := r.URL.Query().Get("lang"); lang != "" && len(data) > 0 {
		for i := range data {
			if translated := translateText(r.Context(), data[i].IsiTeks, lang); translated != "" {
				data[i].IsiTeks = translated
			}
		}
	}

	writeJSON(w, http.StatusOK, model.PaginatedResponse[model.KontenResponse]{
		Data: model.KontenResponse{
			KitabID:  kitabID,
			Judul:    judul,
			Penulis:  penulis,
			Sections: data,
		},
		Pagination: model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages(total, limit),
		},
	})
}
