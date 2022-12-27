package repository

import (
	"fmt"
	"github.com/Daniil-K/new-year-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(name, url string, userId, chatId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (name, url, user_id, chat_id) values ($1, $2, $3, $4)", usersTable)

	tx.QueryRow(createItemQuery, name, url, userId, chatId)

	return tx.Commit()
}

func (r *UserPostgres) GetAll() ([]models.Users, error) {
	var users []models.Users
	query := fmt.Sprintf(`SELECT * FROM %s`, usersTable)
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}
