package storage

import (
	"context"
	"database/sql"
	"errors"

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
	_, err := s.DB.ExecContext(context.Background(),
		"INSERT INTO shortener (uuid, short_url, original_url) VALUES ($1, $2, $3)",
		data.UUID, data.ShortURL, data.OriginalURL)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) Get(key string) (ShortenedData, error) {
	var (
		uuid        string
		shortURL    string
		originalURL string
	)

	row := s.DB.QueryRowContext(context.Background(),
		"SELECT uuid, short_url, original_url FROM shortener WHERE short_url = $1", key)
	err := row.Scan(&uuid, &shortURL, &originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return ShortenedData{}, errors.New("key not found: " + key)
		} else {
			return ShortenedData{}, err
		}
	}
	return ShortenedData{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}, nil
}

func (s *DBStorage) Close() error {
	return s.DB.Close()
}
