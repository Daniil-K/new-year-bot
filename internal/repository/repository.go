package repository

import (
	"github.com/Daniil-K/new-year-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type Wish interface {
	Create(text string, userId int) error
	GetAll(userId int) ([]models.Wishes, error)
	GetAllRecep(userId int) ([]models.Wishes, error)
	//Delete(userId, itemId int) error
	//Update(userId, wishId int, newText string) error
}

type User interface {
	Create(name, url string, userId, chatId int) error
	GetAll() ([]models.Users, error)
}

type Santa interface {
	Create(santaId, userId int) error
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
