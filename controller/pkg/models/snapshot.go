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
	ID         string       `json:"id" db:"space_id"`
	Name       string       `json:"name"`
	About      string       `json:"about"`
	Network    string       `json:"network"`
	Symbol     string       `json:"symbol"`
	Created    int64        `json:"created"`
	Strategies []Strategies `json:"strategies" `
	Admins     []string     `json:"admins"`
	Members    []string     `json:"members"`
	Filters    Filters      `json:"filters"`
}

type Strategies struct {
	Name string `json:"name"`
	//Params []struct{} `json:"params"`
}

type Filters struct {
	MinScore    int  `json:"min_score"`
	OnlyMembers bool `json:"only_members"`
}

type Votes struct {
	ID      string  `json:"id"`
	Voter   string  `json:"voter"`
	Vp      float32 `json:"vp"`
	VpState string  `json:"vp_state"`
	Created int64   `json:"created"`
	Choice  int64   `json:"choice"`
}
