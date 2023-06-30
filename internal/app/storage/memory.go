package storage

import (
	"errors"
	"log"
)

type InMemoryStorage struct {
	data            map[string]string
	FileStoragePath string
}

func NewInMemoryStorage(fileStoragePath string) *InMemoryStorage {
	return &InMemoryStorage{
		data:            make(map[string]string),
		FileStoragePath: fileStoragePath,
	}
}

func (s *InMemoryStorage) Store(key string, value string) error {
	s.data[key] = value

	if s.FileStoragePath == "" {
		return nil
	}

	err := SaveShortURLToStorage(key, value, s.FileStoragePath)
	if err != nil {
		log.Println("failed to save short url to storage")
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
