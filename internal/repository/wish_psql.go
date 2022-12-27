package repository

import (
	"fmt"
	"github.com/Daniil-K/new-year-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type WishPostgres struct {
	db *sqlx.DB
}

func NewWishPostgres(db *sqlx.DB) *WishPostgres {
	return &WishPostgres{db: db}
}

func (r *WishPostgres) Create(text string, userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createItemQuery := fmt.Sprintf("INSERT INTO %s (text, user_id) values ($1, $2)", wishesTable)

	tx.QueryRow(createItemQuery, text, userId)

	return tx.Commit()
}

func (r *WishPostgres) GetAll(userId int) ([]models.Wishes, error) {
	var wishes []models.Wishes
	query := fmt.Sprintf(
		`SELECT * FROM %s WHERE user_id = $1`,
		wishesTable,
	)
	if err := r.db.Select(&wishes, query, userId); err != nil {
		return nil, err
	}

	return wishes, nil
}

func (r *WishPostgres) GetAllRecep(userId int) ([]models.Wishes, error) {
	var wishes []models.Wishes
	query := fmt.Sprintf(
		`SELECT w.text FROM %s w INNER JOIN %s u ON w.user_id = u.user_id
                                        INNER JOIN %s s ON s.recep_id = u.user_id
                                        WHERE s.santa_id = $1`,
		wishesTable,
		usersTable,
		santaTable,
	)
	if err := r.db.Select(&wishes, query, userId); err != nil {
		return nil, err
	}

	return wishes, nil
}

/*
func (r *WishItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}
*/
