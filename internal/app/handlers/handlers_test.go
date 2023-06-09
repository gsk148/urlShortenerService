package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeneralHandler(t *testing.T) {
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader("https://practicum.yandex.ru/"))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			GeneralHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}
