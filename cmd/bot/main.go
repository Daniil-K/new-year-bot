package main

import (
	"github.com/Daniil-K/new-year-bot/internal/repository"
	"github.com/Daniil-K/new-year-bot/pkg/config"
	"github.com/Daniil-K/new-year-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: cfg.PostgresLogin,
		DBName:   cfg.PostgresDBName,
		SSLMode:  "disable",
		Password: cfg.PostgresPass,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	telegramBot := telegram.NewBot(bot, cfg.Messages, repos)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}

}
