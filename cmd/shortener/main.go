package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gsk148/urlShorteningService/internal/app/comperess"
	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

func main() {
	cfg := config.Load()
	store := storage.NewInMemoryStorage()

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        *store,
	}
	logger.NewLogger()
	r := chi.NewRouter()
	r.Post(`/`, comperess.CompressGzip(logger.WithLogging(http.HandlerFunc(h.ShortenerHandler))))
	r.Post(`/api/shorten`, comperess.CompressGzip(logger.WithLogging(http.HandlerFunc(h.ShortenerAPIHandler))))
	r.Get(`/{id}`, comperess.CompressGzip(logger.WithLogging(http.HandlerFunc(h.FindByShortLinkHandler))))
	err := http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
