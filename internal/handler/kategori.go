package handler

import (
	"net/http"

	"github.com/brolatief/fiqih-tahawulat/internal/model"
)

func (h *Handler) ListKategori(w http.ResponseWriter, r *http.Request) {
	data, err := h.svc.ListKategori(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.KategoriResponse{
		Data:  data,
		Total: len(data),
	})
}
