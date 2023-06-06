package main

import (
	"io"
	"math/rand"
	"net/http"
	"time"
)

var shortString = ""
var urlsMap = make(map[string]string)
var addedUrls = make(map[string]string)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func addShortOrGetLargeURL(w http.ResponseWriter, r *http.Request) {
	parsedURL := r.URL
	if parsedURL.Path != "" && r.Method == http.MethodGet {

		url := parsedURL.Path[1:]
		if urlsMap[url] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("Location: " + urlsMap[url]))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Not supported", http.StatusBadRequest)
		return
	}

	requestData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading requestData failed", http.StatusBadRequest)
		return
	}

	if len(requestData) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	inputString := string(requestData)
	if addedUrls[inputString] != "" {
		shortString = addedUrls[inputString]
	} else {
		shortString = createRandomStringFromInput(inputString)
	}

	urlsMap[shortString] = inputString
	content := "http://localhost:8080/" + shortString

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(content))
}

func createRandomStringFromInput(url string) string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	short := string(b)
	addedUrls[url] = short

	return short
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, addShortOrGetLargeURL)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
