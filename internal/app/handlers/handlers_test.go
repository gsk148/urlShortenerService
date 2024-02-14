package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

func TestCreateShortLinkHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "create short link success test",
			requestMethod: http.MethodPost,
			requestPath:   "/",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name:          "create short link not valid method",
			requestMethod: http.MethodGet,
			requestPath:   "/",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader("https://practicum.yandex.ru/"))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.Shorten(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestFindByShortLinkHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "find url not valid method",
			requestMethod: http.MethodPost,
			requestPath:   "/arzKvEKh",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "not valid length url test",
			requestMethod: http.MethodGet,
			requestPath:   "/T",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.FindByShortLink(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestShorterApiHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		contentType   string
		body          map[string]string
		want          want
	}{
		{
			name:          "api short link success test",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten",
			contentType:   "application/json",
			body:          map[string]string{"url": "https://practicum.yandex.ru"},
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name:          "api short link not valid content-type",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten",
			contentType:   "text/plain",
			body:          map[string]string{"url": "https://practicum.yandex.ru"},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(test.requestMethod, test.requestPath, bytes.NewBuffer(body))
			request.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()
			h.ShortenAPI(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestPingHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "success ping test",
			requestMethod: http.MethodGet,
			requestPath:   "/ping",
			want: want{
				code:        200,
				contentType: "",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.Ping(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestDeleteURLs(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		requestBody   string
		want          want
	}{
		{
			name:          "failed delete urls test",
			requestMethod: http.MethodDelete,
			requestPath:   "/api/user/urls",
			requestBody:   "",
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name:          "success delete urls test",
			requestMethod: http.MethodDelete,
			requestPath:   "/api/user/urls",
			requestBody:   "[\"6qxTVvsy\", \"RTfd56hn\", \"Jlfd67ds\"]",
			want: want{
				code:        202,
				contentType: "",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader(test.requestBody))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.DeleteURLs(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestBatchShortenerAPIHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		requestBody   string
		want          want
	}{
		{
			name:          "fail batch shortner api test",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten/batch",
			requestBody:   "",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "success batch shortner api test",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten/batch",
			requestBody:   "[{\"correlation_id\":\"first\",\"original_url\":\"https://ya.ru\"},{\"correlation_id\":\"second\",\"original_url\":\"https://rambler.ru\"}]",
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader(test.requestBody))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.BatchShortenAPI(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestFindUserURLS(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "fail batch shortner api test",
			requestMethod: http.MethodGet,
			requestPath:   "/api/user/urls",
			want: want{
				code:        200,
				contentType: "application/json",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.FindUserURLS(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestGetStats(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name          string
		requestMethod string
		requestPath   string
		want          want
	}{
		{
			name:          "fail get stats test",
			requestMethod: http.MethodGet,
			requestPath:   "/api/internal/stats",
			want: want{
				code:        403,
				contentType: "",
			},
		},
	}

	myLog := logger.NewLogger()
	h := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "",
		Store:         storage.NewInMemoryStorage(),
		Logger:        *myLog,
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			w := httptest.NewRecorder()
			h.GetStats(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}
