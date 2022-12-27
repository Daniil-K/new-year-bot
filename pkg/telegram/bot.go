package telegram

import (
	"github.com/Daniil-K/new-year-bot/internal/repository"
	"github.com/Daniil-K/new-year-bot/pkg/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Bot struct {
	bot *tgbotapi.BotAPI

	messages config.Messages

	repos *repository.Repository
}

func NewBot(bot *tgbotapi.BotAPI, messages config.Messages, repos *repository.Repository) *Bot {
	return &Bot{bot: bot, messages: messages, repos: repos}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // If we got a message
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		err := b.handleMessage(update.Message)
		if err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u), nil
}
