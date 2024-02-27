package storage

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/gsk148/urlShorteningService/internal/app/api"
)

// ErrURLExists structure of special error
type ErrURLExists struct{}

// Error returns string message
func (e *ErrURLExists) Error() string {
	return "URL already exists"
}

// DBStorage structure of DBStorage
type DBStorage struct {
	DB     *sql.DB
	logger zap.SugaredLogger
}

// NewDBStorage return DBStorage object
func NewDBStorage(dsn string, logger zap.SugaredLogger) (*DBStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS shortener (
            id SERIAL PRIMARY KEY,
            user_id TEXT NOT NULL,
            uuid TEXT NOT NULL,
            short_url TEXT NOT NULL UNIQUE,
            original_url TEXT NOT NULL,
            is_deleted BOOLEAN
        );
		CREATE UNIQUE INDEX  IF NOT EXISTS shortener_original_url_uindex
		    on public.shortener (original_url);
    `)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		DB:     db,
		logger: logger,
	}, nil
}

// Ping ping db
func (s *DBStorage) Ping() error {
	if err := s.DB.PingContext(context.Background()); err != nil {
		return err
	}
	return nil
}

// Store saves data to DB and return error if already exists and short url if not
func (s *DBStorage) Store(data api.ShortenedData) (api.ShortenedData, error) {
	result, err := s.DB.ExecContext(context.Background(),
		"INSERT INTO shortener (uuid, user_id, short_url, original_url, is_deleted) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (original_url) DO NOTHING",
		data.UUID, data.UserID, data.ShortURL, data.OriginalURL, data.IsDeleted)
	if err != nil {
		return api.ShortenedData{}, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return api.ShortenedData{}, err
	}

	if affectedRows == 0 {
		row := s.DB.QueryRowContext(context.Background(),
			"SELECT uuid, user_id, short_url, original_url FROM shortener WHERE original_url = $1", data.OriginalURL)
		var existingData api.ShortenedData
		err := row.Scan(&existingData.UUID, &existingData.UserID, &existingData.ShortURL, &existingData.OriginalURL)
		if err != nil {
			return api.ShortenedData{}, err
		}
		return existingData, &ErrURLExists{}
	}

	return data, nil
}

// Get returns full url by short url
func (s *DBStorage) Get(key string) (api.ShortenedData, error) {
	var (
		uuid        string
		userID      string
		shortURL    string
		originalURL string
		isDeleted   bool
	)

	row := s.DB.QueryRowContext(context.Background(),
		"SELECT uuid, user_id, short_url, original_url, is_deleted FROM shortener WHERE short_url = $1", key)
	err := row.Scan(&uuid, &userID, &shortURL, &originalURL, &isDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return api.ShortenedData{}, errors.New("key not found: " + key)
		} else {
			return api.ShortenedData{}, err
		}
	}
	return api.ShortenedData{
		UserID:      userID,
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		IsDeleted:   isDeleted,
	}, nil
}

// Close return nil if ok or error
func (s *DBStorage) Close() error {
	return s.DB.Close()
}

// GetBatchByUserID returns batches of short urls by provided userID
func (s *DBStorage) GetBatchByUserID(userID string) ([]api.ShortenedData, error) {
	var (
		entity api.ShortenedData
		result []api.ShortenedData
	)
	query := "select short_url, original_url from shortener where user_id=$1"
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&entity.ShortURL, &entity.OriginalURL)
		if err != nil {
			break
		}
		result = append(result, entity)
	}
	if len(result) == 0 {
		return nil, errors.New("no batches by provided userID")
	}
	return result, nil
}

// DeleteByUserIDAndShort delete full url from db by userID and short url
func (s *DBStorage) DeleteByUserIDAndShort(userID string, short string) error {
	query := "UPDATE shortener SET is_deleted=true WHERE user_id=$1 AND short_url=$2"
	rows, err := s.DB.Exec(query, userID, short)
	if err != nil {
		return err
	}
	if r, err := rows.RowsAffected(); err != nil || r == 0 {
		s.logger.Info("0 rows affected in delete")
		return err
	}
	s.logger.Infof("Marked as deleted link %s", short)
	return nil
}

// GetStatistic - return num of saved urls and users
func (s *DBStorage) GetStatistic() *api.Statistic {
	var st api.Statistic
	query := "SELECT count(DISTINCT user_id), count(*) FROM shortener"
	res := s.DB.QueryRow(query)
	err := res.Scan(&st.Users, &st.URLs)
	if err != nil {
		return nil
	}
	return &st
}
