package storage

import "github.com/gsk148/urlShorteningService/internal/app/config"

type ShortenedData struct {
	UserID      string `json:"userID"`
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage interface {
	Store(data ShortenedData) (ShortenedData, error)
	Get(key string) (ShortenedData, error)
	Ping() error
	Close() error
	GetBatchByUserID(userID string) ([]ShortenedData, error)
}

func NewStorage(cfg config.Config) (Storage, error) {
	switch cfg.StorageType {
	case "memory":
		return NewInMemoryStorage(), nil
	case "file":
		return NewFileStorage(cfg.FileStoragePath)
	case "db":
		return NewDBStorage(cfg.DatabaseDSN)
	default:
		return NewInMemoryStorage(), nil
	}
}
