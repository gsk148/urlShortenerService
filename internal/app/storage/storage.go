package storage

import (
	"go.uber.org/zap"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/config"
)

// ShortenedData model for url info
type ShortenedData struct {
	UserID      string `json:"userID"`
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	IsDeleted   bool   `json:"is_deleted"`
}

// Storage interface with included needed methods
type Storage interface {
	Store(data ShortenedData) (ShortenedData, error)
	Get(key string) (ShortenedData, error)
	Ping() error
	Close() error
	GetBatchByUserID(userID string) ([]ShortenedData, error)
	DeleteByUserIDAndShort(userID string, shortURL string) error
	GetStatistic() *api.Statistic
}

// NewStorage return NewStorage object
func NewStorage(cfg config.Config, logger zap.SugaredLogger) (Storage, error) {
	switch cfg.StorageType {
	case "memory":
		return NewInMemoryStorage(), nil
	case "file":
		return NewFileStorage(cfg.FileStoragePath)
	case "db":
		return NewDBStorage(cfg.DatabaseDSN, logger)
	default:
		return NewInMemoryStorage(), nil
	}
}
