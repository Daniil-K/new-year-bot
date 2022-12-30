package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	commandStart      = "start"
	commandGifts      = "start_gifts"
	commandSanta      = "start_santa"
	commandClearSanta = "clear_santa"

	myWishes    = "Мои пожелания"
	recepWishes = "Пожелания получателя"

	adminId = 398382229
)

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(myWishes),
		tgbotapi.NewKeyboardButton(recepWishes),
	),
)

var wishKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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

	case commandClearSanta:
		if message.From.ID == adminId {
			err := b.clearSanta(message)
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

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	err := b.repos.Wish.Delete(strconv.FormatInt(callback.Message.Chat.ID, 10), callback.Message.Text)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, b.messages.UnableToDelete)
		_, err := b.bot.Send(msg)
		log.Println(err)
		return err
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, b.messages.SuccessDelete)
	_, err = b.bot.Send(msg)

	msgDelete := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	_, err = b.bot.Send(msgDelete)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
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

	}

	if message.Text != myWishes && message.Text != recepWishes {
		err = b.addWishes(message)
		if err != nil {
			return err
		}
	}

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

	err := b.repos.User.Create(userName, userUrl, strconv.FormatInt(userId, 10), strconv.FormatInt(chatId, 10))
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

		err = b.repos.Santa.Create(string(rune(user.Id)), string(rune(users[random].Id)))
		if err != nil {
			log.Println(err)
			return err
		}

		notUserUsersId[n] = notUserUsersId[len(notUserUsersId)-1]
		notUserUsersId[len(notUserUsersId)-1] = 100
		notUserUsersId = notUserUsersId[:len(notUserUsersId)-1]

		msg := tgbotapi.NewMessage(int64(user.ChatId), "")

		msg.Text = fmt.Sprintf(
			"Тебе выпадает: %s!\nСсылка на профиль: %s \nНе забудь про подарок!",
			users[random].Name,
			users[random].Url,
		)

		_, err = b.bot.Send(msg)

	}

	return err
}

func (b *Bot) clearSanta(message *tgbotapi.Message) error {
	err := b.repos.Santa.ClearAll()
	if err != nil {
		log.Println(err)
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.ClearSanta)
	_, err = b.bot.Send(msg)
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

func (b *Bot) myWishesSection(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Мои пожелания")
	_, err := b.bot.Send(msg)

	wishes, err := b.repos.Wish.GetAll(strconv.FormatInt(message.From.ID, 10))
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
		msg.ReplyMarkup = wishKeyboard
		_, err = b.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func (b *Bot) recepWishesSection(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, " Пожелания твоего получателя")
	_, err := b.bot.Send(msg)

	wishes, err := b.repos.Wish.GetAllRecep(strconv.FormatInt(message.From.ID, 10))
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
		_, err = b.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

	return err
}

func (b *Bot) addWishes(message *tgbotapi.Message) error {
	if message.IsCommand() {
		return nil
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AddWish)

	err := b.repos.Wish.Create(message.Text, strconv.FormatInt(message.From.ID, 10))
	if err != nil {
		log.Println(err)
		msg.Text = b.messages.UnableToSave
	}

	_, err = b.bot.Send(msg)
	return err
}
