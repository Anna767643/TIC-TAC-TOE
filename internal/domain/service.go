package domain

import (
	"context"
	uuid "github.com/google/uuid"
)

// GameService - бизнес-логика игры.
type GameService interface {
	// CreateNewGame - создаёт новую игру с компьютером или другим игроком.
	// Возвращает информацию о новой игре и ошибку
	CreateNewGame(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (Game, error)

	// GetAvailableGames - находит и возвращает все доступные игры
	// для присоединения игрока с определённым UUID, или ошибку
	GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

	// JoinUser - присоединяет пользователя к игре
	// В случае успеха возвращает измененную игру
	// Ошибки: игра не найдена, к данной игре нельзя добавить пользователя,
	// context.Canceled, или иная неописанная ошибка
	JoinUser(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (Game, error)

	// NextMove - выполняет ход игроком
	// Если следующий ход компьютера, то он выполняется с помощью алгоритма минимакс
	// Возвращает игру с обновлённым полем.
	// Ошибки:
	// ErrGameFinished - игра была ранее окончена
	// ErrContextCanceled - игра прерывается
	// ErrGameNotFound - игра не найдена
	// ErrInvalidBoard - поле нарушает правила или изменены прошлые ходы
	// ErrInvalidCell - недопустимое значение клетки
	// ErrOpponentsMove - ход аппонента
	// ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
	NextMove(ctx context.Context, board Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error)

	// GetGame - получает информацию об игре
	// Ошибки:
	// ErrContextCanceled - игра прерывается
	// ErrGameNotFound - игра не найдена
	// ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
	GetGame(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error)
}

type UserService interface {
	// Register - регистрирует нового пользователя
	// возвращает факт успешной регистрации
	// или ошибку
	Register(ctx context.Context, signUpRequest SignUpRequest) (User, error)

	// LogIn - авторизирует пользователя
	// возвращает UUID пользователя
	// или ошибку
	LogIn(ctx context.Context, signUpRequest SignUpRequest) (uuid.UUID, error)

	// GetUserInfo - получение информации о пользователе по UUID пользователя
	// Возвращает информацию о пользователе
	// Ошибки:
	// ErrContextCanceled
	// ErrUserNotFound
	GetUserInfo(ctx context.Context, uuidUser uuid.UUID) (User, error)
}
