package db

import (
	"context"
	"database/sql"
	"github.com/IDOMATH/portfolio/types"
	"time"
)

type PostgresGuestbookStore struct {
	DB *sql.DB
}

func NewPostgresGuestbookStore(db *sql.DB) *PostgresGuestbookStore {
	return &PostgresGuestbookStore{
		DB: db,
	}
}

func (s *PostgresGuestbookStore) InsertGuestbookSignature(signature types.GuestbookSignature) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into guestbook (name, is_approved, created_at)
				  values ($1, $2, $3) returning id`

	err := s.DB.QueryRowContext(ctx, statement,
		signature.Name,
		signature.IsApproved,
		signature.CreatedAt).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *PostgresGuestbookStore) ApproveGuestbookSignature(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
		update guestbook set is_approved = true
		where id = $2`

	_, err := s.DB.ExecContext(ctx, statement, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresGuestbookStore) DenyGuestbookSignature(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
		delete from guestbook where id = $1`

	_, err := s.DB.ExecContext(ctx, statement, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresGuestbookStore) GetApprovedGuestbookSignatures() ([]types.GuestbookSignature, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var signatures []types.GuestbookSignature

	query := `
		select id, name, email, created_at
		from guestbook
		where is_approved = true`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return signatures, err
	}
	defer rows.Close()

	for rows.Next() {
		var signature types.GuestbookSignature
		err := rows.Scan(
			&signature.Id,
			&signature.Name,
			&signature.CreatedAt,
		)
		if err != nil {
			return signatures, err
		}
		signatures = append(signatures, signature)
	}

	if err = rows.Err(); err != nil {
		return signatures, err
	}

	return signatures, nil
}

func (s *PostgresGuestbookStore) GetAllGuestbookSignatures() ([]types.GuestbookSignature, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var signatures []types.GuestbookSignature

	query := `
		select id, name, email, created_at
		from guestbook`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return signatures, err
	}
	defer rows.Close()

	for rows.Next() {
		var signature types.GuestbookSignature
		err := rows.Scan(
			&signature.Id,
			&signature.Name,
			&signature.IsApproved,
			&signature.CreatedAt,
		)
		if err != nil {
			return signatures, err
		}
		signatures = append(signatures, signature)
	}

	if err = rows.Err(); err != nil {
		return signatures, err
	}

	return signatures, nil
}
