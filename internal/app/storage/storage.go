package storage

import "github.com/gsk148/urlShorteningService/internal/app/config"

type ShortenedData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Store(data ShortenedData) error
	Get(key string) (ShortenedData, error)
	Ping() error
	Close() error
}

func NewStorage(cfg config.Config) Storage {
	switch cfg.StorageType {
	case "memory":
		return NewInMemoryStorage()
	case "file":
		return NewFileStorage(cfg.FileStoragePath)
	case "db":
		return NewDBStorage(cfg.DatabaseDSN)
	default:
		return NewInMemoryStorage()
	}
}
