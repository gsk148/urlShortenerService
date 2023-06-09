package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

	h := &Handler{
		ShortURLAddr: "http://localhost:8080",
		Store:        *storage.NewInMemoryStorage(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader("https://practicum.yandex.ru/"))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.ShortenerHandler(w, request)

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

	h := &Handler{
		ShortURLAddr: "http://localhost:8080",
		Store:        *storage.NewInMemoryStorage(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.FindByShortLinkHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}
