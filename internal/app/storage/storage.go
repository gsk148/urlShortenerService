package storage

import (
	"go.uber.org/zap"

	"github.com/gsk148/urlShorteningService/internal/app/api"
	"github.com/gsk148/urlShorteningService/internal/app/config"
)

// Storage interface with included needed methods
type Storage interface {
	Store(data api.ShortenedData) (api.ShortenedData, error)
	Get(key string) (api.ShortenedData, error)
	Ping() error
	Close() error
	GetBatchByUserID(userID string) ([]api.ShortenedData, error)
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
