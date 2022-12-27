package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"time"
)

const (
	commandStart = "start"
	commandGifts = "start_gifts"
	commandSanta = "start_santa"

	myWishes    = "Мои пожелания"
	recepWishes = "Пожелания получателя"
	resume      = "Назад"

	adminId = 398382229
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои пожелания"),
		tgbotapi.NewKeyboardButton("Пожелания получателя"),
	),
)

var resumeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Назад"),
	),
)

var wishesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "delete"),
	),
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды ")

	switch message.Command() {
	case commandStart:
		err := b.startMessage(message)
		return err

	case commandSanta:
		if message.From.ID == adminId {
			err := b.startSanta()
			return err
		}
		return nil

	case commandGifts:
		if message.From.ID == adminId {
			err := b.startGifts()
			return err
		}
		return nil

	default:
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	log.Printf("Id chat: [%d]", message.Chat.ID)

	var err error

	switch message.Text {
	case myWishes:
		err = b.myWishesSection(message)
		if err != nil {
			return err
		}

	case recepWishes:
		err = b.recepWishesSection(message)
		if err != nil {
			return err
		}

	case resume:
		err = b.mainMenu(message)
		if err != nil {
			return err
		}
	}

	if message.Text != myWishes && message.Text != resume && message.Text != recepWishes {
		err = b.addWishesSection(message)
		if err != nil {
			return err
		}
	}

	return err
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, callback.Data)

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) startMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	userName := fmt.Sprintf("%s %s", message.From.FirstName, message.From.LastName)
	userUrl := fmt.Sprintf("https://t.me/%s", message.From.UserName)
	userId := message.From.ID
	chatId := message.Chat.ID

	msg.Text = fmt.Sprintf("Привет %s!\nРад приветствовать в новогоднем боте 'Тайный Санта' \nДавай начнём!", userName)
	msg.ReplyMarkup = startKeyboard
	_, _ = b.bot.Send(msg)

	err := b.repos.User.Create(userName, userUrl, int(userId), int(chatId))
	if err != nil {
		log.Println(err)
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")
	msg.Text = b.messages.Start

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) startSanta() error {
	notUserUsersId := []int{}

	users, err := b.repos.User.GetAll()
	if err != nil {
		log.Println(err)
		return err
	}

	for i := 0; i < len(users); i++ {
		notUserUsersId = append(notUserUsersId, i)
	}

	for _, user := range users {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(notUserUsersId))
		random := notUserUsersId[n]

		if user.Id == users[random].Id {
			rand.Seed(time.Now().UnixNano())
			n = rand.Intn(len(notUserUsersId))
			random = notUserUsersId[n]
		}

		err = b.repos.Santa.Create(user.Id, users[random].Id)
		if err != nil {
			log.Println(err)
			return err
		}

		notUserUsersId[n] = notUserUsersId[len(notUserUsersId)-1]
		notUserUsersId[len(notUserUsersId)-1] = 100
		notUserUsersId = notUserUsersId[:len(notUserUsersId)-1]

		msg := tgbotapi.NewMessage(int64(user.ChatId), "")

		msg.Text = fmt.Sprintf(
			"Твой получатель: %s!\nСсылка на его профиль: %s \nНе забудь про подарок!",
			users[random].Name,
			users[random].Url,
		)

		_, err = b.bot.Send(msg)

	}
	/*
		msg := tgbotapi.NewMessage(message.Chat.ID, "")
		userName := fmt.Sprintf("%s %s", message.From.FirstName, message.From.LastName)

		msg.Text = fmt.Sprintf("Твой получатель: %s!\nСсылка на его профиль: %s \nНе забудь про подарок!", userName, userName)

		_, err = b.bot.Send(msg)

	*/
	return err
}

func (b *Bot) startGifts() error {
	users, err := b.repos.User.GetAll()
	if err != nil {
		log.Println(err)
	}

	for _, user := range users {
		msg := tgbotapi.NewMessage(int64(user.ChatId), b.messages.StartGifts)
		_, err = b.bot.Send(msg)
		return err
	}

	return err
}

func (b *Bot) mainMenu(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	msg.Text = fmt.Sprintf("Ты находишься в главном меню!")
	msg.ReplyMarkup = startKeyboard

	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) myWishesSection(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Мои пожелания")
	msg.ReplyMarkup = resumeKeyboard
	_, err := b.bot.Send(msg)

	wishes, err := b.repos.Wish.GetAll(int(message.From.ID))
	if err != nil {
		log.Println(err)
	}

	if len(wishes) == 0 {
		msg = tgbotapi.NewMessage(message.Chat.ID, b.messages.EmptyWishes)
		_, err = b.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

	for _, wish := range wishes {
		msg = tgbotapi.NewMessage(message.Chat.ID, wish.Text)
		msg.ReplyMarkup = wishesKeyboard
		_, err = b.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func (b *Bot) recepWishesSection(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, " Пожелания твоего получателя")
	msg.ReplyMarkup = resumeKeyboard

	_, err := b.repos.Wish.GetAllRecep(int(message.From.ID))
	if err != nil {
		log.Println(err)
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) addWishesSection(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AddWish)

	err := b.repos.Wish.Create(message.Text, int(message.From.ID))
	if err != nil {
		log.Println(err)
		msg.Text = b.messages.UnableToSave
	}

	_, err = b.bot.Send(msg)
	return err
}
