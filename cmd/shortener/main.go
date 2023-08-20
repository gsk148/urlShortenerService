package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gsk148/urlShorteningService/internal/app/compress"

	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

func main() {
	cfg := config.Load()

	store, err := storage.NewStorage(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger.NewLogger()

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
	}

	r := chi.NewRouter()

	r.Use(middleware.Compress(5,
		"application/javascript",
		"application/json",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml"))
	r.Use(compress.Middleware)
	r.Use(logger.WithLogging)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/shorten", h.ShortenerAPIHandler)
		r.Post("/api/shorten/batch", h.BatchShortenerAPIHandler)
		r.Get("/api/user/urls", h.FindUserURLS)
	})

	r.Post("/", h.ShortenerHandler)
	r.Get("/{id}", h.FindByShortLinkHandler)
	r.Get("/ping", h.PingHandler)

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
