package main

import (
	"io"
	"log"
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

func addShortOrGetLargeUrl(w http.ResponseWriter, r *http.Request) {

	parsedUrl := r.URL
	if parsedUrl.Path != "" && r.Method == http.MethodGet {

		url := parsedUrl.Path[1:]
		if urlsMap[url] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("Location: " + urlsMap[url]))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Allowed only POST method"))
		return
	}

	requestData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
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
	mux.HandleFunc(`/`, addShortOrGetLargeUrl)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
