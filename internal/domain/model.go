package domain

import (
	"time"

	uuid "github.com/google/uuid"
)

// Значение клетки поля
const (
	X     = 1 // Ход игрока
	O     = 2 // Ход компьютера
	Empty = 0 // Пустая клетка
)

// Size - размер стороны поля
const Size = 3

// Board - игровое поле, целочисленная матрица 3×3.
// Индексация Board[row][col], значения Empty / X / O.
type Board [Size][Size]int

// Поля Game.GameState
var (
	WaitUsers = "wait users"
	NextMove  = "next move user " // + UUID user
	Draw      = "draw"
	Win       = "win user " // + UUID user
)

// Game - текущая игра: идентификатор и состояние поля.
type Game struct {
	ID           int
	UUID         uuid.UUID
	Board        Board
	WithComputer bool
	GameState    string

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
