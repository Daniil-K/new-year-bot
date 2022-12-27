package models

type Wishes struct {
	Text   string `db:"text"`
	UserId int    `db:"user_id"`
}
