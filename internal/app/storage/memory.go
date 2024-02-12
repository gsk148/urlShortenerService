package storage

import (
	"errors"

	"github.com/gsk148/urlShorteningService/internal/app/api"
)

// InMemoryStorage structure of InMemoryStorage
type InMemoryStorage struct {
	data map[string]api.ShortenedData
}

// NewInMemoryStorage return NewInMemoryStorage object
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]api.ShortenedData),
	}
}

// Store data and return error if already exists and short url if not
func (s *InMemoryStorage) Store(data api.ShortenedData) (api.ShortenedData, error) {
	s.data[data.ShortURL] = data
	return api.ShortenedData{}, nil
}

// Get returns full url by short url
func (s *InMemoryStorage) Get(key string) (api.ShortenedData, error) {
	value, exists := s.data[key]
	if !exists {
		return api.ShortenedData{}, errors.New("key not found: " + key)
	}
	return value, nil
}

// Ping return nil
func (s *InMemoryStorage) Ping() error {
	return nil
}

// Close return nil
func (s *InMemoryStorage) Close() error {
	return nil
}

// GetBatchByUserID returns batches of short urls by provided userID
func (s *InMemoryStorage) GetBatchByUserID(userID string) ([]api.ShortenedData, error) {
	var data []api.ShortenedData
	data = append(data, api.ShortenedData{})

	return data, nil
}

// DeleteByUserIDAndShort return error
func (s *InMemoryStorage) DeleteByUserIDAndShort(userID string, shortURL string) error {
	return errors.New("Error")
}

// GetStatistic - return num of saved urls and users
func (s *InMemoryStorage) GetStatistic() *api.Statistic {
	return nil
}
