package main

import (
	"io"
	"math/rand"
	"net/http"
	"time"
)

var urlsMap = make(map[string]string)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handler(w http.ResponseWriter, r *http.Request) {
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
		encoded := createRandomStringFromInput(string(body))
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

func createRandomStringFromInput(url string) string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
