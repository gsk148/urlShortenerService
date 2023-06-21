package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

func main() {
	cfg := config.Load()
	store := storage.NewInMemoryStorage()

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        *store,
	}
	r := chi.NewRouter()
	r.Post(`/`, h.ShortenerHandler)
	r.Get(`/{id}`, h.FindByShortLinkHandler)
	err := http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
