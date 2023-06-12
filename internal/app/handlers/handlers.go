package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/utils/hasher"
)

var urlsMap = make(map[string]string)

func CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
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
	urlsMap[encoded] = string(body)
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	url := config.GetBaseURL() + "/" + encoded
	w.Write([]byte(url))
}

func FindByShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}
	shortLink := chi.URLParam(r, "id")
	if len(shortLink) < 2 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if _, ok := urlsMap[shortLink]; ok {
		w.Header().Set("Location", urlsMap[shortLink])
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		http.Error(w, "Short url not found", http.StatusBadRequest)
		return
	}
}
