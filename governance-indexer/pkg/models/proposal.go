package models

type Proposals struct {
	ID       string   `json:"id,omitempty"`       // Уникальный идентификатор пропозиции
	Title    string   `json:"title,omitempty"`    // Текст заголовка
	Author   string   `json:"author,omitempty"`   // Автор (адрес валидного кошелька)
	Created  int64    `json:"created,omitempty"`  // Время создания записи
	Start    int64    `json:"start,omitempty"`    // Время начала голосования
	End      int64    `json:"end,omitempty"`      // Время окончание голосования
	Snapshot int64    `json:"snapshot,omitempty"` // Номер блока, на котором базируется голосование
	State    string   `json:"state,omitempty"`    // Статус (active, closed, pending)
	Choices  []string `json:"choices,omitempty"`  // Варианты для голосования
	Space    Space    `json:"space,omitempty"`    // Информация по токену
}

type Space struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
