package config

import (
	"flag"
	"os"
)

// Config contains environment variables which should be set
type Config struct {
	ServerAddr      string
	ShortURLAddr    string
	StorageType     string
	FileStoragePath string
	DatabaseDSN     string
	EnableHTTPS     string
}

// Load gets env vars from arguments or environment
func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "The starting server address (format: host:port)")
	flag.StringVar(&cfg.ShortURLAddr, "b", "http://localhost:8080", "Returned address: net address host:port")
	flag.StringVar(&cfg.StorageType, "storage", "file", "type of storage to use (memory/file)")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database host")
	flag.StringVar(&cfg.EnableHTTPS, "s", "", "Server would be run on TLS")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.ServerAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.ShortURLAddr = envBaseAddr
	}
	if envStorageType := os.Getenv("STORAGE_TYPE"); envStorageType != "" {
		cfg.StorageType = envStorageType
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		cfg.EnableHTTPS = envEnableHTTPS
	}

	if cfg.DatabaseDSN != "" {
		cfg.StorageType = "db"
	}

	return cfg
}
