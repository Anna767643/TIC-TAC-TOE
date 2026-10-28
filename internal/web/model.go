package web

// BoardDTO - поле в JSON: массив 3×3 целых чисел (0 пусто, 1 X, 2 O).
type boardDTO [3][3]int

type reqBoard struct {
	BoardDTO boardDTO `json:"board"`
}

// newGame - запрос при создании новой игры
type reqNewGame struct {
	WithComputer *bool `json:"with_computer"`
}

// GameDTO - модель игры
type gameDTO struct {
	ID           int      `json:"id"`
	UUID         string   `json:"uuid"`
	Board        boardDTO `json:"board"`
	WithComputer bool     `json:"with_computer"`
	GameState    string   `json:"game_state"`

	UUIDPlayerFirst  string  `json:"uuid_player_first"`  // X
	UUIDPlayerSecond *string `json:"uuid_player_second"` // O

	UpdatedAt string `json:"update_at"`
}

// User - модель пользователя
type userDTO struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Login    string `json:"login"`
	CreateAt string `json:"createAt"`
}
