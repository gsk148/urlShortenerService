package storage

import (
	"errors"
)

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Store(key string, value string) {
	s.data[key] = value
}

func (s *InMemoryStorage) Get(key string) (string, error) {
	value, exists := s.data[key]
	if !exists {
		return "", errors.New("key not found: " + key)
	}
	return value, nil
}
