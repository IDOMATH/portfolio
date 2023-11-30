package db

import (
	"context"
	"database/sql"
	"github.com/IDOMATH/portfolio/types"
	"time"
)

type FitnessStore interface {
	Dropper
}

type PostgresFitnessStore struct {
	DB *sql.DB
	FitnessStore
}

func NewPostgresFitnessStore(db *sql.DB) *PostgresFitnessStore {
	return &PostgresFitnessStore{
		DB: db,
	}
}

func (s *PostgresFitnessStore) InsertFitnessRecap(recap types.FitnessRecap) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	statement := `insert into fitness_recaps (weight, distance, date)
				  values ($1, $2, $3) returning id`

	err := s.DB.QueryRowContext(ctx, statement,
		recap.TenthsOfAPound,
		recap.HundredthsOfAMile,
		recap.Date).Scan(&newId)

	if err != nil {
		return 0, err
	}
	return newId, nil
}

func (s *PostgresFitnessStore) GetAllFitnessRecaps() ([]types.FitnessRecap, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var recaps []types.FitnessRecap

	query := `
		select weight, distance, date
		from fitness_recaps
		order by date`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return recaps, err
	}

	defer rows.Close()
	for rows.Next() {
		var row types.FitnessRecap
		err := rows.Scan(
			&row.TenthsOfAPound,
			&row.HundredthsOfAMile,
			&row.Date)

		if err != nil {
			return recaps, err
		}
		recaps = append(recaps, row)
	}

	if err = rows.Err(); err != nil {
		return recaps, err
	}

	return recaps, nil
}
