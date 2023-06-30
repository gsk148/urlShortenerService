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

//type Consumer struct {
//	file    *os.File
//	decoder *json.Decoder
//}

//func NewConsumer(fileName string) (*Consumer, error) {
//	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Consumer{
//		file:    file,
//		decoder: json.NewDecoder(file),
//	}, nil
//}
//
//func (c *Consumer) ReadEventFromFile() (*Event, error) {
//	event := &Event{}
//
//	if err := c.decoder.Decode(&event); err != nil {
//		return nil, err
//	}
//
//	return event, nil
//}
//
//func (c *Consumer) ReadEvent() (*Event, error) {
//	event := &Event{}
//	if err := c.decoder.Decode(&event); err != nil {
//		return nil, err
//	}
//
//	return event, nil
//}

func (p *Producer) SaveToFileStorage(shortURL, originalURL string) error {
	event := Event{
		UUID:        uuid.NewString(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	err := p.WriteEvent(&event)
	if err != nil {
		return err
	}
	return nil
}
