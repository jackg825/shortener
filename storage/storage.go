package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(databaseURL string) *Storage {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	return &Storage{conn: conn}
}

func (s *Storage) Close() {
	s.conn.Close(context.Background())
}

func (s *Storage) SaveShortURL(shortURL, longURL, userId string) (string, error) {
	var sURL string
	err := s.conn.QueryRow(context.Background(), "SELECT short_url FROM urls WHERE short_url = $1 AND long_url = $2 AND user_id = $3", shortURL, longURL, userId).Scan(&sURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Pair of short and long URL not found, insert new
			_, err = s.conn.Exec(context.Background(), "INSERT INTO urls (short_url, long_url, user_id) VALUES ($1, $2, $3)", shortURL, longURL, userId)
			if err != nil {
				return "", err
			}
		} else {
			// Handle unexpected error
			return "", err
		}
	}
	// If the pair of short and long URL already exists, do nothing

	return shortURL, nil
}

func (s *Storage) GetLongURL(shortURL string) (string, error) {
	var longURL string
	err := s.conn.QueryRow(context.Background(), "SELECT long_url FROM urls WHERE short_url = $1", shortURL).Scan(&longURL)
	if err != nil {
		log.Printf("Failed to get URL: %v", err)
		return "", err
	}
	return longURL, nil
}
