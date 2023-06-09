package handlers

import (
	"github.com/gsk148/urlShorteningService/internal/app/utils/hasher"
	"io"
	"net/http"
)

var urlsMap = make(map[string]string)

func GeneralHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Reading requestBody failed", http.StatusBadRequest)
			return
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		if len(body) == 0 {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}
		encoded := hasher.CreateHash()
		urlsMap[encoded] = string(body)
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		url := scheme + "://" + r.Host + "/" + encoded
		w.Write([]byte(url))
	} else if r.Method == http.MethodGet {
		path := r.URL.Path
		if len(path) < 2 {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		path = path[1:]
		if _, ok := urlsMap[path]; ok {
			w.Header().Set("Location", urlsMap[path])
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		} else {
			http.Error(w, "Short url not found", http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}
}
