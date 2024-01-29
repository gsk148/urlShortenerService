package config

import (
	"encoding/json"
	"flag"
	"os"
)

// Config contains environment variables which should be set
type Config struct {
	ServerAddr      string `json:"server_address"`
	ShortURLAddr    string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
	StorageType     string
	Config          string `env:"CONFIG"`
}

// Load gets env vars from arguments or environment
func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "The starting server address (format: host:port)")
	flag.StringVar(&cfg.ShortURLAddr, "b", "http://localhost:8080", "Returned address: net address host:port")
	flag.StringVar(&cfg.StorageType, "storage", "file", "type of storage to use (memory/file)")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Database host")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "Enable HTTPS server mode")
	flag.StringVar(&cfg.Config, "c", "", "JSON config file")
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

	if cfg.DatabaseDSN != "" {
		cfg.StorageType = "db"
	}

	// Read and parse JSON file if flag -c with value exists
	jsonFileData, err := os.ReadFile(cfg.Config)
	if err != nil {
		return cfg
	}

	var jsonCfg Config
	if err = json.Unmarshal(jsonFileData, &jsonCfg); err != nil {
		return cfg
	}

	if cfg.ServerAddr == "" {
		cfg.ServerAddr = jsonCfg.ServerAddr
	}
	if cfg.ShortURLAddr == "" {
		cfg.ShortURLAddr = jsonCfg.ShortURLAddr
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = jsonCfg.FileStoragePath
	}
	if cfg.DatabaseDSN == "" {
		cfg.DatabaseDSN = jsonCfg.DatabaseDSN
	}
	if !cfg.EnableHTTPS || jsonCfg.EnableHTTPS {
		cfg.EnableHTTPS = jsonCfg.EnableHTTPS
	}

	return cfg
}
