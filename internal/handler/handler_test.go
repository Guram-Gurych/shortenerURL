package handler

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	CreateShortURLFunc func(ctx context.Context, originalURL string) (string, error)
	GetOriginalURLFunc func(ctx context.Context, id string) (string, error)
}

func (m *MockService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	return m.CreateShortURLFunc(ctx, originalURL)
}

func (m *MockService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	return m.GetOriginalURLFunc(ctx, id)
}

func TestPostHandler(t *testing.T) {
	type testCase struct {
		name           string
		requestBody    string
		mockID         string
		mockError      error
		expectedStatus int
		expectedBody   string
	}
	tests := []testCase{
		{
			name:           "Успешное создание",
			requestBody:    "https://google.com",
			mockID:         "E9wVbL1G",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/E9wVbL1G",
		},
		{
			name:           "Пустое тело запроса",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Ошибка от сервиса",
			requestBody:    "https://yandex.ru",
			mockError:      errors.New("не удалось сохранить"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		reqBody := strings.NewReader(test.requestBody)
		req := httptest.NewRequest(http.MethodPost, "/", reqBody)
		recorder := httptest.NewRecorder()

		mockService := &MockService{
			CreateShortURLFunc: func(ctx context.Context, originalURL string) (string, error) {
				return test.mockID, test.mockError
			},
		}

		baseURL := "http://localhost:8080"
		handler := NewHandler(mockService, baseURL, nil)

		handler.Post(recorder, req)

		res := recorder.Result()
		defer res.Body.Close()

		assert.Equal(t, test.expectedStatus, res.StatusCode, "Код ответа не совпадает")

		if test.expectedBody != "" {
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.expectedBody, string(body), "Тело ответа не совпадает")
		}
	}
}

func TestGetHandler(t *testing.T) {
	type testCase struct {
		name             string
		requestURL       string
		method           string
		mockOriginalURL  string
		mockError        error
		expectedStatus   int
		expectedLocation string
	}

	tests := []testCase{
		{
			name:             "Успешный редирект",
			requestURL:       "/shortID123",
			method:           http.MethodGet,
			mockOriginalURL:  "https://practicum.yandex.ru/",
			mockError:        nil,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://practicum.yandex.ru/",
		},
		{
			name:             "URL не найден",
			requestURL:       "/notFoundID",
			method:           http.MethodGet,
			mockOriginalURL:  "",
			mockError:        errors.New("URL not found"),
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
		{
			name:             "ID не указан в пути",
			requestURL:       "/",
			method:           http.MethodGet,
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
		{
			name:             "Неверный HTTP-метод",
			requestURL:       "/someID",
			method:           http.MethodPost,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.requestURL, nil)
			recorder := httptest.NewRecorder()

			mockService := &MockService{
				GetOriginalURLFunc: func(ctx context.Context, id string) (string, error) {
					if id == "shortID123" {
						return test.mockOriginalURL, test.mockError
					}
					return "", errors.New("URL not found")
				},
			}

			handler := NewHandler(mockService, "http://localhost:8080", nil)
			router := chi.NewRouter()
			router.Get("/{id}", handler.Get)

			router.ServeHTTP(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedStatus, res.StatusCode, "Статус-код не совпадает")

			if test.expectedLocation != "" {
				assert.Equal(t, test.expectedLocation, res.Header.Get("Location"), "Заголовок Location не совпадает")
			}
		})
	}
}

func TestPostShortenHandler(t *testing.T) {
	type testCase struct {
		name           string
		requestBody    string
		contentType    string
		mockID         string
		mockError      error
		expectedStatus int
		expectedBody   string
	}
	tests := []testCase{
		{
			name:           "Успешное создание JSON",
			requestBody:    `{"url": "https://practicum.yandex.ru"}`,
			contentType:    "application/json",
			mockID:         "E9wVbL1G",
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"result": "http://localhost:8080/E9wVbL1G"}`,
		},
		{
			name:           "Ошибка: Неверный Content-Type",
			requestBody:    `{"url": "https://google.com"}`,
			contentType:    "text/plain",
			mockError:      nil,
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   "",
		},
		{
			name:           "Ошибка: Невалидный JSON",
			requestBody:    `{"url": "https://malformed.com"`,
			contentType:    "application/json",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Ошибка: Пустое поле URL",
			requestBody:    `{"url": ""}`,
			contentType:    "application/json",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "Ошибка: Ошибка от сервиса",
			requestBody:    `{"url": "https://yandex.ru"}`,
			contentType:    "application/json",
			mockError:      errors.New("не удалось сохранить"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := strings.NewReader(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)
			req.Header.Set("Content-Type", test.contentType)

			recorder := httptest.NewRecorder()

			mockService := &MockService{
				CreateShortURLFunc: func(ctx context.Context, originalURL string) (string, error) {
					return test.mockID, test.mockError
				},
			}

			baseURL := "http://localhost:8080"
			handler := NewHandler(mockService, baseURL, nil)

			handler.PostShorten(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedStatus, res.StatusCode, "Код ответа не совпадает")

			if test.expectedBody != "" {
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.JSONEq(t, test.expectedBody, string(body), "Тело ответа не совпадает")
			}
		})
	}
}
