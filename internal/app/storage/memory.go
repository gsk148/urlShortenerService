package storage

import (
	"errors"
)

type InMemoryStorage struct {
	data map[string]ShortenedData
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]ShortenedData),
	}
}

func (s *InMemoryStorage) Store(data ShortenedData) (ShortenedData, error) {
	s.data[data.ShortURL] = data
	return ShortenedData{}, nil
}

func (s *InMemoryStorage) Get(key string) (ShortenedData, error) {
	value, exists := s.data[key]
	if !exists {
		return ShortenedData{}, errors.New("key not found: " + key)
	}
	return value, nil
}

func (s *InMemoryStorage) Ping() error {
	return nil
}

func (s *InMemoryStorage) Close() error {
	return nil
}

func (s *InMemoryStorage) GetBatchByUserID(userID string) ([]ShortenedData, error) {
	var data []ShortenedData
	data = append(data, ShortenedData{})

	return data, nil
}

func (s *InMemoryStorage) DeleteByUserIDAndShort(userID string, shortURL string) error {
	return errors.New("Error")
}
