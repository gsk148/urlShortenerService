package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/gsk148/urlShorteningService/internal/app/api"
)

// FileStorage structure of FileStorage
type FileStorage struct {
	inMemoryData *InMemoryStorage
	filePath     string
}

// NewFileStorage return NewFileStorage object
func NewFileStorage(filename string) (*FileStorage, error) {
	inMemoryData := NewInMemoryStorage()
	fs := FileStorage{
		inMemoryData: inMemoryData,
		filePath:     filename,
	}

	err := readFromFile(fs)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}

func readFromFile(fs FileStorage) error {
	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()

		var sd api.ShortenedData
		if err := json.Unmarshal(line, &sd); err != nil {
			return err
		}

		_, err := fs.inMemoryData.Store(sd)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Store data and return error if already exists and short url if not
func (s *FileStorage) Store(data api.ShortenedData) (api.ShortenedData, error) {
	s.inMemoryData.data[data.ShortURL] = data
	err := s.Save()
	return api.ShortenedData{}, err
}

// Get returns full url by short url
func (s *FileStorage) Get(key string) (api.ShortenedData, error) {
	data, exists := s.inMemoryData.data[key]
	if !exists {
		return api.ShortenedData{}, errors.New("key not found: " + key)
	}
	return data, nil
}

// Save data to file storage
func (s *FileStorage) Save() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, v := range s.inMemoryData.data {
		line, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = writer.Write(append(line, '\n'))
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Ping return nil
func (s *FileStorage) Ping() error {
	return nil
}

// Close return nil if ok or error
func (s *FileStorage) Close() error {
	return s.Save()
}

// GetBatchByUserID returns batches of short urls by provided userID
func (s *FileStorage) GetBatchByUserID(userID string) ([]api.ShortenedData, error) {
	var data []api.ShortenedData
	data = append(data, api.ShortenedData{})

	return data, nil
}

// DeleteByUserIDAndShort return error
func (s *FileStorage) DeleteByUserIDAndShort(userID string, shortURL string) error {
	return errors.New("Error")
}

// GetStatistic - returns num of saved urls and users
func (s *FileStorage) GetStatistic() *api.Statistic {
	return nil
}
