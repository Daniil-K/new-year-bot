package models

type Users struct {
	Id     int    `db:"user_id"`
	Name   string `db:"name"`
	Url    string `db:"url"`
	ChatId int    `db:"chat_id"`
}
