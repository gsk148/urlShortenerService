package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type ErrURLExists struct{}

func (e *ErrURLExists) Error() string {
	return "URL already exists"
}

type DBStorage struct {
	DB *sql.DB
}

func NewDBStorage(dsn string) (*DBStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS shortener (
            id SERIAL PRIMARY KEY,
            uuid TEXT NOT NULL,
            short_url TEXT NOT NULL UNIQUE,
            original_url TEXT NOT NULL
        );
		CREATE UNIQUE INDEX  IF NOT EXISTS shortener_original_url_uindex
		    on public.shortener (original_url);
    `)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		DB: db,
	}, nil
}

func (s *DBStorage) Ping() error {
	if err := s.DB.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) Store(data ShortenedData) (ShortenedData, error) {
	result, err := s.DB.ExecContext(context.Background(),
		"INSERT INTO shortener (uuid, short_url, original_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING",
		data.UUID, data.ShortURL, data.OriginalURL)
	if err != nil {
		return ShortenedData{}, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return ShortenedData{}, err
	}

	if affectedRows == 0 {
		row := s.DB.QueryRowContext(context.Background(),
			"SELECT uuid, short_url, original_url FROM shortener WHERE original_url = $1", data.OriginalURL)
		var existingData ShortenedData
		err := row.Scan(&existingData.UUID, &existingData.ShortURL, &existingData.OriginalURL)
		if err != nil {
			return ShortenedData{}, err
		}
		return existingData, &ErrURLExists{}
	}

	return data, nil
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
