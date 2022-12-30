package repository

import (
	"github.com/Daniil-K/new-year-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type Wish interface {
	Create(text, userId string) error
	GetAll(userId string) ([]models.Wishes, error)
	GetAllRecep(userId string) ([]models.Wishes, error)
	Delete(userId, itemId string) error
}

type User interface {
	Create(name, url, userId, chatId string) error
	GetAll() ([]models.Users, error)
}

type Santa interface {
	Create(santaId, userId string) error
	ClearAll() error
}

type Repository struct {
	Wish
	User
	Santa
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Wish:  NewWishPostgres(db),
		User:  NewUserPostgres(db),
		Santa: NewSantaPostgres(db),
	}
}
