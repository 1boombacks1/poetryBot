package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"poetryLibrary/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AnswerFormat(
	foundPoems []models.Poem,
	foundAuthors map[string][]models.Poem,
) (string, tgbotapi.InlineKeyboardMarkup) {

	keyboardMarkup := make([][]tgbotapi.InlineKeyboardButton, 0)
	answer := fmt.Sprintf("Найден %d стих, писатели:\n", len(foundPoems))
	authors := ""

	for author, poems := range foundAuthors {
		authors += fmt.Sprintf("<b>%s - %d</b>\n", author, len(poems))
		poemTitles := ""
		keyboardRow := make([]tgbotapi.InlineKeyboardButton, 0)

		for i, poem := range poems {
			// formatPoem := fmt.Sprintf("*Автор: %s*\n*Название: %s*\n\n%s", poem.Author, poem.Title, poem.Text)
			keyboardRow = append(keyboardRow,
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+1), "1"),
			)
			poemTitles += fmt.Sprintf("<i>%d) %s</i>\n", i+1, poem.Title)
			if (i+1)%6 == 0 {
				keyboardMarkup = append(keyboardMarkup, keyboardRow)
				keyboardRow = make([]tgbotapi.InlineKeyboardButton, 0)
			}
		}
		keyboardMarkup = append(keyboardMarkup, keyboardRow)
		authors += poemTitles
	}
	answer += authors
	if len(keyboardMarkup) == 0 {
		log.Print("ТУТ 0 БРО")
	}
	return answer, tgbotapi.NewInlineKeyboardMarkup(keyboardMarkup...)
}

func ChangeUserState(chatID int64, state int) {
	users, err := GetUsers()
	if err != nil {
		log.Print("Получили юзеров👤👤👤")
		return
	}
	for i, user := range users {
		if user.ChatID == chatID {
			log.Print("Юзер найден")
			users[i].State = state
		}
	}
	if err := WriteUsers(users); err != nil {
		return
	}
	log.Print("Записали пользователей👤👤👤")
}

func WriteUsers(users []models.User) error {
	data, err := json.MarshalIndent(users, "", "	")
	if err != nil {
		log.Printf("Ошибка в маршале пользователей!☢️\n%s", err.Error())
		return err
	}

	if err := ioutil.WriteFile("Data/users.json", data, 0); err != nil {
		log.Printf("Ошибка в записи данных пользователей в файл!☢️\n%s", err.Error())
		return err
	}
	return nil
}

func GetUsers() ([]models.User, error) {
	file, err := ioutil.ReadFile("Data/users.json")
	if err != nil {
		log.Printf("Произошла ошибка в чтении файла users.json\n%s", err.Error())
		return nil, err
	}

	var users []models.User

	if err = json.Unmarshal(file, &users); err != nil {
		log.Printf("Не удалось анмаршануть данные!\n%s", err.Error())
		return nil, err
	}

	return users, nil
}

func GetPoems() ([]models.Poem, error) {
	file, err := ioutil.ReadFile("Data/out.json")
	if err != nil {
		log.Printf("Не удалось прочитать файл!\n%s", err.Error())
		return nil, err
	}

	var poems []models.Poem

	if err = json.Unmarshal(file, &poems); err != nil {
		log.Printf("Не удалось анмаршануть данные!\n%s", err.Error())
		return nil, err
	}

	return poems, nil
}

func MustToken() string {
	token := flag.String(
		"bot-token",
		"",
		"telegram bot token for the application to work",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("No token enterned")
	}

	return *token
}
