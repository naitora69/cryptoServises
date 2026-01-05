package models

type Proposals struct {
	ID       string   `json:"id"`       // Уникальный идентификатор пропозиции
	Title    string   `json:"title"`    // Текст заголовка
	Author   string   `json:"author"`   // Автор (адрес валидного кошелька)
	Created  int64    `json:"created"`  // Время создания записи
	Start    int64    `json:"start"`    // Время начала голосования
	End      int64    `json:"end"`      // Время окончание голосования
	Snapshot int64    `json:"snapshot"` // Номер блока, на котором базируется голосование
	State    string   `json:"state"`    // Статус (active, closed, pending)
	Choices  []string `json:"choices"`  // Варианты для голосования
	Space    Space    `json:"space"`    // Информация по токену
}

type Space struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
