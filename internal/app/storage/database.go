package storage

import (
	"context"
	"database/sql"
)

type DBStorage struct {
	DB *sql.DB
}

func NewDBStorage(dsn string) *DBStorage {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS shortener (
            id SERIAL PRIMARY KEY,
            uuid TEXT NOT NULL,
            short_url TEXT NOT NULL UNIQUE,
            original_url TEXT NOT NULL
        );
    `)
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
