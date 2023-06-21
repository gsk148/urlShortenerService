package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddr   string
	ShortURLAddr string
	StorageType  string
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "The starting server address (format: host:port)")
	flag.StringVar(&cfg.ShortURLAddr, "b", "http://localhost:8080", "Returned address: net address host:port")
	flag.StringVar(&cfg.StorageType, "s", "memory", "type of storage to use (memory)")
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.ServerAddr = envRunAddr
	}
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.ShortURLAddr = envBaseAddr
	}

	return cfg
}
