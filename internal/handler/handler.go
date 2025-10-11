package handler

import (
	"fmt"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Handler struct {
	service service.URLShortener
	baseURL string
}

func NewHandler(s service.URLShortener, baseURL string) *Handler {
	return &Handler{
		service: s,
		baseURL: baseURL,
	}
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Тело запроса не может быть пустым", http.StatusBadRequest)
		return
	}

	originalURL := string(body)

	id, err := h.service.CreateShortURL(originalURL)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
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
		http.Error(w, "ID не может быть пустым", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetOriginalURL(id)
	if err != nil {
		http.Error(w, "URL не найден", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
