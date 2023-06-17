package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/urlShorteningService/internal/app/storage"
	"github.com/gsk148/urlShorteningService/internal/app/utils/hasher"
)

type Handler struct {
	ShortURLAddr string
	Store        storage.InMemoryStorage
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
	err = h.Store.Store(encoded, string(body))
	if err != nil {
		return
	}
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	url := h.ShortURLAddr + "/" + encoded
	w.Write([]byte(url))
}

func (h *Handler) FindByShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}
	shortLink := chi.URLParam(r, "id")
	if len(shortLink) < 2 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	url, err := h.Store.Get(shortLink)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
