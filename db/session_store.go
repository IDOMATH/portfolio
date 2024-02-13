package db

import (
	"context"
	"database/sql"
	"time"
)

type PostgresSessionStore struct {
	DB *sql.DB
}

func NewSessionStore(db *sql.DB) *PostgresSessionStore {
	return &PostgresSessionStore{
		DB: db,
	}
}

func (s *PostgresSessionStore) InsertSessionToken(token string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into auth_tokens (token, expires_at)
				  values ($1, $2) returning id`

	err := s.DB.QueryRowContext(ctx, statement,
		token,
		time.Now().Add(time.Hour*2)).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *PostgresSessionStore) CheckSessionToken(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var sessionToken string
	var expiresAt time.Time
	statement := `SELECT token, expires_at FROM auth_tokens 
				  WHERE token = $1 `

	err := s.DB.QueryRowContext(ctx, statement,
		token,
		time.Now().Add(time.Hour*2)).Scan(sessionToken, expiresAt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
