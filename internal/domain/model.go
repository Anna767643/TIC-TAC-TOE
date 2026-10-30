package domain

import (
	"time"

	uuid "github.com/google/uuid"
)

// Значение клетки поля
const (
	X     = 1 // Ход первого игрока
	O     = 2 // Ход второго игрока
	Empty = 0 // Пустая клетка
)

// Size - размер стороны поля
const Size = 3

// Board - игровое поле, целочисленная матрица 3×3.
// Индексация Board[row][col], значения Empty / X / O.
type Board [Size][Size]int

// Поля Game.GameState
type GameState string

const (
	StateWaitUsers GameState = "wait users"
	StateNextMove  GameState = "next move"
	StateDraw      GameState = "draw"
	StateWin       GameState = "win"
)

// Game - текущая игра: идентификатор и состояние поля.
type Game struct {
	ID           int
	UUID         uuid.UUID
	Board        Board
	WithComputer bool
	State        GameState
	NextPlayer   *uuid.UUID // кто ходит, когда State == StateNextMove
	Winner       *uuid.UUID // nil при победе компьютера

	UUIDPlayerFirst  uuid.UUID
	UUIDPlayerSecond *uuid.UUID

	UpdatedAt time.Time
}

// SignUpRequest - модель для регистрации и авторизации
type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// User - модель пользователя
type User struct {
	ID           int
	UUID         uuid.UUID
	Login        string
	PasswordHash string
	CreateAt     time.Time
}
