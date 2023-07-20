package storage

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	DB *sql.DB
}

func NewDBStorage(dsn string) *DBStorage {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return &DBStorage{
		DB: db,
	}
}

func (s *DBStorage) Ping() error {
	if err := s.DB.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) Store(data ShortenedData) error {
	return nil
}

func (s *DBStorage) Get(key string) (ShortenedData, error) {
	return ShortenedData{}, nil
}

func (s *DBStorage) Close() error {
	return s.DB.Close()
}
