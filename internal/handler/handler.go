package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/brolatief/fiqih-tahawulat/internal/service"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 4096)
		return &b
	},
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, map[string]string{"error": code, "message": msg})
}

func parseIntDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return i
}

func clampInt(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func totalPages(total, limit int) int {
	if total == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
