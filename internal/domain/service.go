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

	// GetAvaibleGames - находит и возвращает все доступные игры
	// для присоединения игрока с определённым UUID, или ошибку
	GetAvaibleGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

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

	// Validate проверяет, что incoming - допустимое продолжение stored:
	// все непустые клетки stored не изменились, добавлена ровно одна клетка
	// и это ход пользователя (X), значения клеток в диапазоне.
	// Возвращает ErrInvalidBoard или ErrInvalidCell с описанием причины.
	// Validate(incoming Board, stored Board, playerCellValue int) error

	// IsFinished сообщает, окончена ли игра, и кто победил:
	// winner = X, O или Empty при ничьей / незавершённой игре.
	// IsFinished(board Board) (bool, int)
}

type UserService interface {
	// Register - регистрирует нового пользователя
	// возвращает факт успешной регистрации
	// или ошибку
	Register(ctx context.Context, signUpRequest SignUpRequest) (User, error)

	// LogIn - авторизирует пользователя
	// возвращает UUIDпользователя
	// или ошибку
	LogIn(ctx context.Context, signUpRequest SignUpRequest) (string, error)

	// GetUserInfo - получение информации о пользователе по UUID пользователя
	// Возвращает информацию о пользователе
	// Ошибки:
	// ErrContextCanceled
	// ErrUserNotFound
	GetUserInfo(ctx context.Context, uuidUser uuid.UUID) (User, error)
}
