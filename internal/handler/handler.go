package handler

import (
	"fmt"
	"github.com/Guram-Gurych/shortenerURL.git/internal/service"
	"io"
	"net/http"
	"strings"
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

func (h *Handler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Разрешены только Post запросы", http.StatusMethodNotAllowed)
		return
	}

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

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Разрешены только Get запросы", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
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
