// Package handlers contains public API handlers
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/auth"
	"github.com/gsk148/urlShorteningService/internal/app/compress"
	"github.com/gsk148/urlShorteningService/internal/app/hashutil"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

// Handler structure of Handler
type Handler struct {
	BaseURL       string
	TrustedSubnet string
	Store         storage.Storage
	Logger        zap.SugaredLogger
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5,
		"application/javascript",
		"application/json",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml"))
	r.Use(compress.Middleware)
	r.Use(logger.WithLogging)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten", h.ShortenAPI)
		r.Post("/api/shorten/batch", h.BatchShortenAPI)
		r.Get("/api/user/urls", h.FindUserURLS)
		r.Delete("/api/user/urls", h.DeleteURLs)
	})

	r.Post("/", h.Shorten)
	r.Get("/{id}", h.FindByShortLink)
	r.Get("/ping", h.Ping)
	r.Get("/api/internal/stats", h.GetStats)

	r.HandleFunc("/debug/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	r.HandleFunc("/debug/pprof/*", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

	return r
}

// Shorten save provided in text/plain format full url and returns short
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading requestBody failed", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	encoded := hashutil.Encode(body)

	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storedData, err := h.Store.Store(api.ShortenedData{
		UserID:      userID,
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: string(body),
	})
	if err != nil {
		if errors.Is(err, &storage.ErrURLExists{}) {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			url := h.BaseURL + "/" + storedData.ShortURL
			_, err = w.Write([]byte(url))
			if err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	url := h.BaseURL + "/" + encoded
	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// FindByShortLink returns full url by provided id
func (h *Handler) FindByShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}
	shortLink := chi.URLParam(r, "id")
	data, err := h.Store.Get(shortLink)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if data.IsDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", data.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// ShortenAPI save provided in json format full url and returns short
func (h *Handler) ShortenAPI(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Not valid content type", http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading requestBody failed", http.StatusBadRequest)
		return
	}

	var request api.ShortenRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Unmarshalling request failed", http.StatusBadRequest)
		return
	}

	encoded := hashutil.Encode([]byte(request.URL))
	var response api.ShortenResponse
	response.Result = h.BaseURL + "/" + encoded
	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Marshaling response failed", http.StatusBadRequest)
		return
	}

	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	_, err = h.Store.Store(api.ShortenedData{
		UserID:      userID,
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: request.URL,
		IsDeleted:   false,
	})
	if err != nil {
		if errors.Is(err, &storage.ErrURLExists{}) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusConflict)

			_, err = w.Write(result)
			if err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// Ping makes test connection to storage
func (h *Handler) Ping(res http.ResponseWriter, req *http.Request) {
	if err := h.Store.Ping(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

// BatchShortenAPI saves array of provided urls
func (h *Handler) BatchShortenAPI(w http.ResponseWriter, r *http.Request) {
	var reqItems []api.BatchShortenRequestItem
	err := json.NewDecoder(r.Body).Decode(&reqItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respItems := make([]api.BatchShortenResponseItem, 0, len(reqItems))
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	for _, reqItem := range reqItems {
		shortURL := hashutil.Encode([]byte(reqItem.OriginalURL))

		_, err := h.Store.Store(api.ShortenedData{
			UserID:      userID,
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: reqItem.OriginalURL,
			IsDeleted:   false,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respItems = append(respItems, api.BatchShortenResponseItem{
			CorrelationID: reqItem.CorrelationID,
			ShortURL:      h.BaseURL + "/" + shortURL,
		})
	}

	respBytes, err := json.Marshal(respItems)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(respBytes)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// FindUserURLS returns array of all saved by user urls
func (h *Handler) FindUserURLS(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	batch, err := h.Store.GetBatchByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if len(batch) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	type UserLinksResponse struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	var result []UserLinksResponse

	for _, v := range batch {
		shortURL := h.BaseURL + "/" + v.ShortURL
		b := &UserLinksResponse{shortURL, v.OriginalURL}
		result = append(result, *b)
	}

	response, err := json.Marshal(result)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// DeleteURLs removes array of provided urls
func (h *Handler) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var inputArray []string
	userID, err := auth.GetUserToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	urls, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(urls, &inputArray)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputCh := addShortURLs(inputArray)
	go h.MarkAsDeleted(inputCh, userID)

	w.WriteHeader(http.StatusAccepted)
}

// MarkAsDeleted set flag deleted=true for provided short url
func (h *Handler) MarkAsDeleted(inputShort chan string, userID string) {
	for v := range inputShort {
		err := h.Store.DeleteByUserIDAndShort(userID, v)
		if err != nil {
			h.Logger.Warnf("Failed to mark deleted by short %s", v)
		}
	}
}

// GetStats returns count of urls and users
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	isTrusted, err := h.checkIPIsTrusted(realIP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isTrusted {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stat := h.Store.GetStatistic()
	if stat == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(stat)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) checkIPIsTrusted(clientIP string) (bool, error) {
	_, trustedIP, err := net.ParseCIDR(h.TrustedSubnet)
	if err != nil {
		return false, err
	}

	parsedIP := net.ParseIP(clientIP)
	if parsedIP == nil {
		return false, err
	}

	if !trustedIP.Contains(parsedIP) {
		return false, nil
	}
	return true, nil
}

func addShortURLs(input []string) chan string {
	inputCh := make(chan string, 10)

	go func() {
		defer close(inputCh)
		for _, url := range input {
			inputCh <- url
		}
	}()

	return inputCh
}
