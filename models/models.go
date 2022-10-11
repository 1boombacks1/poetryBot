package models

const (
	Start = iota
	SearchAuthor
	SearchTitle
)

type Poem struct {
	Author string `json:"Автор"`
	Title  string `json:"Название"`
	Text   string `json:"Текст"`
}

type User struct {
	ChatID    int64  `json:"chat_id"`
	Username  string `json:"username"`
	ReadPoems []int  `json:"read_poems"`
	State     int    `json:"state"`
}
