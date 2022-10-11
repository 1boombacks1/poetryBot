package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"poetryLibrary/keyboards"
	"poetryLibrary/messages"
	"poetryLibrary/models"
	"poetryLibrary/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendPoemCallbackHandler(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, "Принято!🤙")
	if _, err := bot.Request(callback); err != nil {
		log.Printf("Коллбэк не обработан!☢️ %s", err.Error())
		return
	}

	poems, _ := utils.GetPoems()
	users, _ := utils.GetUsers()
	author := strings.Split(callbackQuery.Data, "//")[0]
	title := strings.Split(callbackQuery.Data, "//")[1]
	pattern := "<b>Автор: %s</b>\n<b>Название: %s</b>\n\n%s"

	msg := tgbotapi.NewMessage(callbackQuery.From.ID, "")
	msg.ParseMode = "html"
	msg.ReplyMarkup = keyboards.StartKeyboard

	var poemID int

	for i, poem := range poems {
		if strings.Contains(poem.Author, author) && strings.Contains(poem.Title, title) {
			msg.Text = fmt.Sprintf(pattern, poem.Author, poem.Title, poem.Text)
			poemID = i
			break
		}
	}

	for i, user := range users {
		if callbackQuery.From.ID == user.ChatID {
			users[i].ReadPoems = append(users[i].ReadPoems, poemID)
		}
	}

	utils.WriteUsers(users)
	bot.Send(msg)
}

func SearchCallbackHandler(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(callbackQuery.ID, "Принято!🤙")
	if _, err := bot.Request(callback); err != nil {
		log.Printf("Коллбэк не обработан!☢️ %s", err.Error())
		return
	}

	text := ""
	chatID := callbackQuery.From.ID

	if callbackQuery.Data == keyboards.AUTHOR {
		utils.ChangeUserState(chatID, models.SearchAuthor)
		text = "автора ✍️"
	} else if callbackQuery.Data == keyboards.TITLE {
		utils.ChangeUserState(chatID, models.SearchTitle)
		text = "заголовок 🔤"
	}

	delmsg := tgbotapi.NewDeleteMessage(chatID, callbackQuery.Message.MessageID)
	bot.Send(delmsg)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.CallbackMessage, text))
	bot.Send(msg)
}

func SearchHandler(bot *tgbotapi.BotAPI, chatID int64, pattern string) {
	if pattern == "_MagicWORD_" {
		msg := tgbotapi.NewMessage(chatID, messages.SearchMessage)
		msg.ReplyMarkup = keyboards.InlineSearchKeyboard
		bot.Send(msg)
	} else {
		users, _ := utils.GetUsers()
		poems, _ := utils.GetPoems()
		foundPoems := make([]models.Poem, 0)
		foundAuthors := make(map[string][]models.Poem)

		msg := tgbotapi.NewMessage(chatID, "")
		msg.ReplyMarkup = keyboards.SearchKeyboard
		msg.ParseMode = "html"

		for _, user := range users {
			if user.ChatID == chatID {
				switch user.State {
				case models.SearchAuthor:
					for _, poem := range poems {
						if strings.Contains(strings.ToLower(poem.Author), strings.ToLower(pattern)) {
							foundPoems = append(foundPoems, poem)
							foundAuthors[poem.Author] = append(foundAuthors[poem.Author], poem)
						}
					}

					answer, keyboardMarkup := utils.AnswerFormat(foundPoems, foundAuthors)

					msg.Text = answer
					if len(keyboardMarkup.InlineKeyboard) != 0 {
						msg.ReplyMarkup = keyboardMarkup
					}

					utils.ChangeUserState(chatID, models.Start)

				case models.SearchTitle:
					for _, poem := range poems {
						if strings.Contains(strings.ToLower(poem.Title), strings.ToLower(pattern)) {
							foundPoems = append(foundPoems, poem)
							foundAuthors[poem.Author] = append(foundAuthors[poem.Author], poem)
						}
					}

					answer, keyboardMarkup := utils.AnswerFormat(foundPoems, foundAuthors)
					if len(keyboardMarkup.InlineKeyboard) != 0 {
						msg.ReplyMarkup = keyboardMarkup
					}
					msg.Text = answer

					utils.ChangeUserState(chatID, models.Start)
				}
			}
		}
		bot.Send(msg)
	}
}

func GetRndPoemHandler(bot *tgbotapi.BotAPI, chatID int64) {
	users, err := utils.GetUsers()
	if err != nil {
		log.Printf("Не удалось получить пользователей!☢️\n%s", err.Error())
		return
	}

	poems, err := utils.GetPoems()
	if err != nil {
		errorMsg := tgbotapi.NewMessage(chatID, "Произошла ошибка!☢️")
		bot.Send(errorMsg)
		return
	}

	rand.Seed(time.Now().UnixMicro())
	var poemNumber int

	for i, user := range users {
		if user.ChatID == chatID {
			for {
				isHas := false
				poemNumber = rand.Intn(len(poems))
				for _, n := range user.ReadPoems {
					if poemNumber == n {
						isHas = true
						break
					}
				}
				if !isHas {
					break
				}
			}
			users[i].ReadPoems = append(users[i].ReadPoems, poemNumber)
		}
	}

	utils.WriteUsers(users)

	rndPoem := poems[poemNumber]
	poem := fmt.Sprintf("*Автор: %s*\n*Название: %s*\n\n%s", rndPoem.Author, rndPoem.Title, rndPoem.Text)

	msg := tgbotapi.NewMessage(chatID, poem)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func StartHandler(bot *tgbotapi.BotAPI, chatID int64, username string) {
	user := models.User{
		ChatID:    chatID,
		Username:  username,
		ReadPoems: make([]int, 0),
		State:     models.Start,
	}

	users, err := utils.GetUsers()
	if err != nil {
		log.Printf("Не удалось получить пользователей!☢️\n%s", err.Error())
		return
	}

	for _, user := range users {
		if chatID == user.ChatID {
			msg := tgbotapi.NewMessage(chatID, messages.ReauthorizationErr)
			bot.Send(msg)
			return
		}
	}

	users = append(users, user)

	if err = utils.WriteUsers(users); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, messages.StartMessage)
	msg.ReplyMarkup = keyboards.StartKeyboard
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка в отправке стартового сообщения!☢️")
	}
}
