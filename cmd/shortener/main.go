package main

import (
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.GeneralHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
