// go: build -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=`git rev-parse HEAD`" cmd/shortener/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gsk148/urlShorteningService/internal/app/config"
	"github.com/gsk148/urlShorteningService/internal/app/handlers"
	"github.com/gsk148/urlShorteningService/internal/app/logger"
	"github.com/gsk148/urlShorteningService/internal/app/storage"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func runSrv() (*http.Server, error) {
	cfg := config.Load()

	myLog := logger.NewLogger()
	store, err := storage.NewStorage(*cfg, *myLog)
	if err != nil {
		log.Fatal(err)
	}

	handler := &handlers.Handler{
		BaseURL:       cfg.BaseURL,
		TrustedSubnet: cfg.TrustedSubnet,
		Store:         store,
		Logger:        *myLog,
	}

	router := handler.InitRoutes()

	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	if cfg.EnableHTTPS {
		return srv, http.ListenAndServeTLS(cfg.ServerAddr, "internal/app/cert/server.crt", "internal/app/cert/server.key", router)
	}
	return srv, http.ListenAndServe(cfg.ServerAddr, router)
}

func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	srv, err := runSrv()
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe error: %v", err)
		}
	}()

	sig := <-sigint
	log.Printf("Received signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown error: %v", err)
	}
}
