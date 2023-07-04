package storage

import (
	"errors"
)

type InMemoryStorage struct {
	data            map[string]string
	fileStoragePath string
}

func NewInMemoryStorage(fileStoragePath string) *InMemoryStorage {
	return &InMemoryStorage{
		data:            make(map[string]string),
		fileStoragePath: fileStoragePath,
	}
}

func (s *InMemoryStorage) Store(key string, value string) error {
	s.data[key] = value

	if s.fileStoragePath == "" {
		return nil
	}

	err := SaveShortURLToStorage(key, value, s.fileStoragePath)
	if err != nil {
		return errors.New("failed to save short url to storage")
	}

	return nil
}

func (s *InMemoryStorage) Get(key string) (string, error) {
	value, exists := s.data[key]
	if !exists {
		return "", errors.New("key not found: " + key)
	}
	return value, nil
}

func SaveShortURLToStorage(shortURL, originalURL, fileName string) error {
	producer, err := NewProducer(fileName)
	if err != nil {
		return err
	}
	err = producer.SaveToFileStorage(shortURL, originalURL)
	if err != nil {
		return err
	}
	return nil
}
