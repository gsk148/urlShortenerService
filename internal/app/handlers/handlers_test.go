package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

func getTestHandler(store storage.Storage) *Handler {
	myLog := logger.NewLogger()
	handler := &Handler{
		BaseURL:       "http://localhost:8080",
		TrustedSubnet: "127.0.0.1/24",
		Store:         store,
		Logger:        *myLog,
	}

	return handler
}

func TestInitRoutes(t *testing.T) {
	t.Run("success generate routes", func(t *testing.T) {

		handler := getTestHandler(storage.NewInMemoryStorage())
		routes := handler.InitRoutes()

		assert.NotNil(t, routes)
	})
}

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
		requestData   string
		want          want
	}{
		{
			name:          "create short link success test",
			requestMethod: http.MethodPost,
			requestPath:   "/",
			requestData:   "https://practicum.yandex.ru/",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name:          "create short link not valid method",
			requestMethod: http.MethodGet,
			requestPath:   "/",
			requestData:   "https://sports.ru/",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "create short link not valid data",
			requestMethod: http.MethodPost,
			requestPath:   "/",
			requestData:   "",
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	h := getTestHandler(storage.NewInMemoryStorage())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader(test.requestData))
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

func TestFailCreateShortLink(t *testing.T) {
	t.Run("Create short link", func(t *testing.T) {
		t.Run("fail add to store", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockStorage(ctrl)

			data := api.ShortenedData{}
			dbMock.EXPECT().Store(gomock.Any()).Return(data, new(storage.ErrURLExists))

			handler := getTestHandler(dbMock)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/"))

			w := httptest.NewRecorder()
			handler.Shorten(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusConflict, w.Code)
			defer res.Body.Close()
		})
	})
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

	handler := getTestHandler(storage.NewInMemoryStorage())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			w := httptest.NewRecorder()
			handler.FindByShortLink(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestDeletedFindByShortLink(t *testing.T) {
	t.Run("DeletedFindByShortLink", func(t *testing.T) {
		t.Run("success find active short link", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockStorage(ctrl)
			mockedDBResult := api.ShortenedData{
				OriginalURL: "praktikum.yandex.ru",
				IsDeleted:   false,
			}
			dbMock.EXPECT().Get(gomock.Any()).Return(mockedDBResult, nil)

			handler := getTestHandler(dbMock)
			request := httptest.NewRequest(http.MethodGet, "/ngaCAPJ", nil)
			w := httptest.NewRecorder()
			handler.FindByShortLink(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
			defer res.Body.Close()
		})
		t.Run("success find deleted short link", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockStorage(ctrl)
			mockedDBResult := api.ShortenedData{
				OriginalURL: "praktikum.yandex.ru",
				IsDeleted:   true,
			}
			dbMock.EXPECT().Get(gomock.Any()).Return(mockedDBResult, nil)

			handler := getTestHandler(dbMock)
			request := httptest.NewRequest(http.MethodGet, "/ngaCAPJ", nil)
			w := httptest.NewRecorder()
			handler.FindByShortLink(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusGone, res.StatusCode)
			defer res.Body.Close()
		})
	})
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
			name:          "success case",
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
			name:          "not valid content-type",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten",
			contentType:   "text/plain",
			body:          map[string]string{"url": "https://practicum.yandex.ru"},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:          "not valid content",
			requestMethod: http.MethodPost,
			requestPath:   "/api/shorten",
			contentType:   "text/plain",
			body:          map[string]string{},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	handler := getTestHandler(storage.NewInMemoryStorage())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			request := httptest.NewRequest(test.requestMethod, test.requestPath, bytes.NewBuffer(body))
			request.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()
			handler.ShortenAPI(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestFailCreateShortLinkApi(t *testing.T) {
	t.Run("fail add to store by api", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := storage.NewMockStorage(ctrl)

		data := api.ShortenedData{}
		dbMock.EXPECT().Store(gomock.Any()).Return(data, new(storage.ErrURLExists))

		handler := getTestHandler(dbMock)

		body, _ := json.Marshal(map[string]string{"url": "https://practicum.yandex.ru"})
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
		request.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		handler.ShortenAPI(w, request)

		res := w.Result()
		assert.Equal(t, http.StatusConflict, w.Code)
		defer res.Body.Close()
	})
}

func TestPingHandler(t *testing.T) {
	t.Run("Ping", func(t *testing.T) {
		t.Run("fail case", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockStorage(ctrl)
			dbMock.EXPECT().Ping().Return(errors.New("err"))
			handler := getTestHandler(dbMock)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)

			w := httptest.NewRecorder()
			handler.Ping(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			defer res.Body.Close()
		})
		t.Run("success case", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dbMock := storage.NewMockStorage(ctrl)
			dbMock.EXPECT().Ping().Return(error(nil))
			handler := getTestHandler(dbMock)

			request := httptest.NewRequest(http.MethodGet, "/ping", nil)

			w := httptest.NewRecorder()
			handler.Ping(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusOK, w.Code)
			defer res.Body.Close()
		})
	})
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

	handler := getTestHandler(storage.NewInMemoryStorage())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader(test.requestBody))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			handler.DeleteURLs(w, request)

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

	handler := getTestHandler(storage.NewInMemoryStorage())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.requestMethod, test.requestPath, strings.NewReader(test.requestBody))
			w := httptest.NewRecorder()
			handler.BatchShortenAPI(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("content-type"))
			defer res.Body.Close()
		})
	}
}

func TestFindUserURLS(t *testing.T) {
	t.Run("get user's urls", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h := getTestHandler(storage.NewInMemoryStorage())
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			h.FindUserURLS(w, request)

			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			defer res.Body.Close()
		})
	})

	t.Run("success get user's urls", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := storage.NewMockStorage(ctrl)

		dbMock.EXPECT().GetBatchByUserID(gomock.Any()).Return(nil, errors.New("error"))

		handler := getTestHandler(dbMock)
		request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		w := httptest.NewRecorder()
		handler.FindUserURLS(w, request)

		res := w.Result()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		defer res.Body.Close()
	})
}

func TestGetStats(t *testing.T) {
	t.Run("forbidden get stats", func(t *testing.T) {
		handler := getTestHandler(storage.NewInMemoryStorage())
		request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		w := httptest.NewRecorder()
		handler.GetStats(w, request)

		res := w.Result()
		assert.Equal(t, http.StatusForbidden, res.StatusCode)
		defer res.Body.Close()
	})

	t.Run("fail get stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := storage.NewMockStorage(ctrl)

		dbMock.EXPECT().GetStatistic().Return(nil).AnyTimes()

		h := getTestHandler(dbMock)
		request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		request.Header.Set("X-Real-IP", "127.0.0.1")
		w := httptest.NewRecorder()
		h.GetStats(w, request)

		res := w.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		defer res.Body.Close()
	})

	t.Run("success get stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbMock := storage.NewMockStorage(ctrl)

		stats := &api.Statistic{}
		dbMock.EXPECT().GetStatistic().Return(stats).AnyTimes()

		handler := getTestHandler(dbMock)
		request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
		request.Header.Set("X-Real-IP", "127.0.0.1")
		w := httptest.NewRecorder()
		handler.GetStats(w, request)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()
	})
}
