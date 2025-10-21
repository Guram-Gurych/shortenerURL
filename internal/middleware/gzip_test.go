package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		require.NoError(t, err)
	})

	handlerToTest := GzipMiddleware(dummyHandler)
	srv := httptest.NewServer(handlerToTest)
	defer srv.Close()

	type testCase struct {
		name                     string
		requestBody              string
		compressRequest          bool
		headers                  map[string]string
		expectedStatusCode       int
		expectCompressedResponse bool
		expectedResponseBody     string
	}

	requestBody := `{"url": "https://practicum.yandex.ru"}`

	tests := []testCase{
		{
			name:                     "client sends compressed data",
			requestBody:              requestBody,
			compressRequest:          true,
			headers:                  map[string]string{"Content-Encoding": "gzip", "Content-Type": "application/json"},
			expectedStatusCode:       http.StatusOK,
			expectCompressedResponse: false,
			expectedResponseBody:     requestBody,
		},
		{
			name:                     "client accepts compressed data",
			requestBody:              requestBody,
			compressRequest:          false,
			headers:                  map[string]string{"Accept-Encoding": "gzip", "Content-Type": "application/json"},
			expectedStatusCode:       http.StatusOK,
			expectCompressedResponse: true,
			expectedResponseBody:     requestBody,
		},
		{
			name:                     "no compression",
			requestBody:              requestBody,
			compressRequest:          false,
			headers:                  map[string]string{"Content-Type": "application/json"},
			expectedStatusCode:       http.StatusOK,
			expectCompressedResponse: false,
			expectedResponseBody:     requestBody,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body io.Reader
			if test.compressRequest {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				_, err := gz.Write([]byte(test.requestBody))
				require.NoError(t, err)
				err = gz.Close()
				require.NoError(t, err)
				body = &buf
			} else {
				body = bytes.NewBufferString(test.requestBody)
			}

			req, err := http.NewRequest(http.MethodPost, srv.URL, body)
			require.NoError(t, err)

			for key, value := range test.headers {
				req.Header.Set(key, value)
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)

			var respBody []byte
			if test.expectCompressedResponse {
				assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

				gz, err := gzip.NewReader(resp.Body)
				require.NoError(t, err)
				respBody, err = io.ReadAll(gz)
				require.NoError(t, err)
			} else {
				assert.Empty(t, resp.Header.Get("Content-Encoding"))
				respBody, err = io.ReadAll(resp.Body)
				require.NoError(t, err)
			}

			assert.JSONEq(t, test.expectedResponseBody, string(respBody))
		})
	}
}
