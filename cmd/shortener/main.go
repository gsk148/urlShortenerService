package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gsk148/urlShorteningService/internal/app/compress"
	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// BuildData хранит в себе информацию о текущей версии, дате и значении коммита
type BuildData struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

// Template - переменная, содержит в себе темплейт для stdout с информацией о текущей сборке
const Template = `	Build version: {{if .BuildVersion}} {{.BuildVersion}} {{else}} N/A {{end}}
	Build date: {{if .BuildDate}} {{.BuildDate}} {{else}} N/A {{end}}
	Build commit: {{if .BuildCommit}} {{.BuildCommit}} {{else}} N/A {{end}}`

func buildInfo() {
	d := &BuildData{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}
	t := template.Must(template.New("buildTags").Parse(Template))
	err := t.Execute(os.Stdout, *d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return
	}
}

func main() {
	buildInfo()
	cfg := config.Load()

	myLog := logger.NewLogger()
	store, err := storage.NewStorage(*cfg, *myLog)
	if err != nil {
		log.Fatal(err)
	}

	h := &handlers.Handler{
		ShortURLAddr: cfg.ShortURLAddr,
		Store:        store,
		Logger:       *myLog,
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
		r.Delete("/api/user/urls", h.DeleteURLs)
	})

	r.Post("/", h.ShortenerHandler)
	r.Get("/{id}", h.FindByShortLinkHandler)
	r.Get("/ping", h.PingHandler)

	r.HandleFunc("/debug/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	r.HandleFunc("/debug/pprof/*", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))

	err = http.ListenAndServe(cfg.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
