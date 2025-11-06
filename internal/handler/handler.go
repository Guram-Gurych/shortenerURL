package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

type Handler struct {
	service service.URLShortener
	baseURL string
	db      *sql.DB
}

func NewHandler(s service.URLShortener, baseURL string, db *sql.DB) *Handler {
	return &Handler{
		service: s,
		baseURL: baseURL,
		db:      db,
	}
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body could not be read", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Request body cannot be empty.", http.StatusBadRequest)
		return
	}

	originalURL := string(body)

	ctx := r.Context()
	id, err := h.service.CreateShortURL(ctx, originalURL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, id)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	originalURL, err := h.service.GetOriginalURL(ctx, id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) PostShorten(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req RequestJSON
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL field is missing", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	id, err := h.service.CreateShortURL(ctx, req.URL)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, id)
	resp := ResponseJSON{Result: shortURL}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetPing(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel

	err := h.db.PingContext(ctx)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
