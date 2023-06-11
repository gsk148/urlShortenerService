package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"net/http"
)

func main() {
	addr := flag.String("a", ":8080", "Server address")
	r := chi.NewRouter()
	r.Post(`/`, handlers.CreateShortLinkHandler)
	r.Get(`/{id}`, handlers.FindByShortLinkHandler)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}
