package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SantaPostgres struct {
	db *sqlx.DB
}

func NewSantaPostgres(db *sqlx.DB) *SantaPostgres {
	return &SantaPostgres{db: db}
}

func (r *SantaPostgres) Create(santaId, userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (santa_id, recep_id) values ($1, $2)", santaTable)

	tx.QueryRow(createItemQuery, santaId, userId)

	return tx.Commit()
}
