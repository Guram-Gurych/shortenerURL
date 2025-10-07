package handler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockService struct {
	CreateShortURLFunc func(originalURL string) (string, error)
	GetOriginalURLFunc func(id string) (string, error)
}

func (m *MockService) CreateShortURL(originalURL string) (string, error) {
	return m.CreateShortURLFunc(originalURL)
}

func (m *MockService) GetOriginalURL(id string) (string, error) {
	return m.GetOriginalURLFunc(id)
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
			CreateShortURLFunc: func(originalURL string) (string, error) {
				return test.mockID, test.mockError
			},
		}

		baseURL := "http://localhost:8080"
		handler := NewHandler(mockService, baseURL)

		handler.PostHandler(recorder, req)

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
		mockOriginalURL  string
		mockError        error
		expectedStatus   int
		expectedLocation string
	}
	tests := []testCase{
		{
			name:             "Успешный редирект",
			requestURL:       "/shortID123",
			mockOriginalURL:  "https://practicum.yandex.ru/",
			mockError:        nil,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://practicum.yandex.ru/",
		},
		{
			name:             "URL не найден",
			requestURL:       "/notFoundID",
			mockOriginalURL:  "",
			mockError:        errors.New("URL not found"),
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
		{
			name:             "ID не указан в пути",
			requestURL:       "/",
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
		{
			name:             "Неверный HTTP-метод",
			requestURL:       "/ID",
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedLocation: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestMethod := http.MethodGet
			if test.name == "Неверный HTTP-метод" {
				requestMethod = http.MethodPost
			}

			req := httptest.NewRequest(requestMethod, test.requestURL, nil)
			recorder := httptest.NewRecorder()

			mockService := &MockService{
				GetOriginalURLFunc: func(id string) (string, error) {
					return test.mockOriginalURL, test.mockError
				},
			}

			handler := NewHandler(mockService, "http://localhost:8080")

			handler.GetHandler(recorder, req)

			res := recorder.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedStatus, res.StatusCode, "Статус-код не совпадает")

			if test.expectedLocation != "" {
				assert.Equal(t, test.expectedLocation, res.Header.Get("Location"), "Заголовок Location не совпадает")
			}
		})
	}
}
