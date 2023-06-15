package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
)

func main() {
	config.ParseAddresses()
	r := chi.NewRouter()
	r.Post(`/`, handlers.CreateShortLinkHandler)
	r.Get(`/{id}`, handlers.FindByShortLinkHandler)
	err := http.ListenAndServe(config.GetSrvAddr(), r)
	if err != nil {
		panic(err)
	}
}
