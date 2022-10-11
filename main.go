package main

import (
	"log"

	// "poetryLibrary/parse"

	// "github.com/geziyor/geziyor"
	// "github.com/geziyor/geziyor/export"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"poetryLibrary/handlers"
	"poetryLibrary/keyboards"
	. "poetryLibrary/utils"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(MustToken())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			log.Print("Callback Принят! 🆗🆗🆗")
			if update.CallbackData() != keyboards.AUTHOR && update.CallbackData() != keyboards.TITLE {
				handlers.SendPoemCallbackHandler(bot, update.CallbackQuery)
			} else {
				handlers.SearchCallbackHandler(bot, update.CallbackQuery)
			}
		}
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			username := update.Message.Chat.UserName

			if update.Message.IsCommand() {
				switch update.Message.Text {
				case "/start":
					handlers.StartHandler(bot, chatID, username)
				}
			} else if update.Message.Text == "Получить стих 📜" {
				handlers.GetRndPoemHandler(bot, chatID)
			} else if update.Message.Text == "Поиск 🔍" {
				handlers.SearchHandler(bot, chatID, "_MagicWORD_")
			} else {
				handlers.SearchHandler(bot, chatID, update.Message.Text)
			}
		}

	}
}
