package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
	"github.com/gsk148/urlShorteningService/internal/app/utils/hasher"
)

type Handler struct {
	ShortURLAddr string
	Store        storage.Storage
}

func (h *Handler) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
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

	encoded := hasher.CreateHash()

	err = h.Store.Store(storage.ShortenedData{
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: string(body),
	})
	if err != nil {
		return
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	url := h.ShortURLAddr + "/" + encoded
	_, err = w.Write([]byte(url))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) FindByShortLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Location", data.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) ShortenerAPIHandler(w http.ResponseWriter, r *http.Request) {
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

	encoded := hasher.CreateHash()
	err = h.Store.Store(storage.ShortenedData{
		UUID:        uuid.New().String(),
		ShortURL:    encoded,
		OriginalURL: request.URL,
	})
	if err != nil {
		return
	}

	var response api.ShortenResponse
	response.Result = h.ShortURLAddr + "/" + encoded
	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Marshaling response failed", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PingHandler(res http.ResponseWriter, req *http.Request) {
	if err := h.Store.Ping(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
