package storage

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func (p *Producer) SaveToFileStorage(shortURL, originalURL string) error {
	event := Event{
		UUID:        uuid.NewString(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	return p.WriteEvent(&event)
}
