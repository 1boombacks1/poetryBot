package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	AUTHOR = "author"
	TITLE  = "title"
	TEXT   = "text"
)

const (
	authorFilter = "По автору ✍️"
	titleFilter  = "По названию 🔤"
	textFilter   = "По тексту 📜"
)

const (
	searchBtnText  = "Поиск 🔍"
	getPoemBtnText = "Получить стих 📜"
)

var StartKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(searchBtnText),
		tgbotapi.NewKeyboardButton(getPoemBtnText),
	),
)

var InlineSearchKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(authorFilter, AUTHOR),
		tgbotapi.NewInlineKeyboardButtonData(titleFilter, TITLE),
	),
)

var SearchKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(authorFilter),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(titleFilter),
	),
)
