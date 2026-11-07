package domain

import (
	"context"
	uuid "github.com/google/uuid"
)

// GameService - бизнес-логика игры.
type GameService interface {
	/*
		CreateNewGame - создаёт новую игру с компьютером или другим игроком
		Возвращает информацию о новой игре
		Ошибки:
		ErrContextCanceled
		"repository: insert game %s: %w"
		"repository: failed select game after save: %w"
		"repository: select game %s: %w"
		"repository: failed to map game entity to domain: error"
	*/
	CreateNewGame(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (Game, error)

	/*
		GetAvailableGames - находит и возвращает все доступные игры
		для присоединения игрока с определённым UUID, или ошибку
		Ошибки:
		ErrContextCanceled
		"repository: select game for user %s: %w"
		"repository: scan game for user %s: %w"
	*/
	GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

	/*
		JoinUser - присоединяет пользователя к игре
		В случае успеха возвращает измененную игру
		Ошибки:
		ErrContextCanceled
		"repository: select game %s: %w"
		ErrAddUser
		"repository: update user in game %s: %w"
	*/
	JoinUser(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (Game, error)

	/*
		NextMove - выполняет ход игроком
		Если следующий ход компьютера, то он выполняется с помощью алгоритма минимакс
		Возвращает игру с обновлённым полем.
		Ошибки:
		ErrContextCanceled - игра прерывается
		ErrGameFinished - игра была ранее окончена
		ErrWaitOpponent - второй игрок еще не присоеденился к игре
		ErrGameNotFound - игра не найдена
		ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
		ErrOpponentsMove - ход аппонента
		ErrInvalidBoard - поле нарушает правила или изменены прошлые ходы
		ErrInvalidCell - недопустимое значение клетки
	*/
	NextMove(ctx context.Context, board Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error)

	/*
		GetGame - получает информацию об игре
		Ошибки:
		ErrContextCanceled - игра прерывается
		ErrGameNotFound - игра не найдена
		"repository: select game %s: %w"
		ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
	*/
	GetGame(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error)
}

type UserService interface {
	/*
		Register - регистрирует нового пользователя
		возвращает факт успешной регистрации
		Ошибки:
		ErrContextCanceled
		"domain: faled to generate password hash %w"
		ErrLoginIncorrect
		"repository: insert user with %s: %w"
		"repository: select user with %s: %w"
		"domain: faled to get user %w"
	*/
	Register(ctx context.Context, signUpRequest SignUpRequest) (User, error)

	/*
		LogIn - авторизирует пользователя
		возвращает UUIDпользователя
		Ошибки:
		ErrContextCanceled
		ErrLoginOrPasswordIncorrect
		"repository: select user with %s: %w"
	*/
	LogIn(ctx context.Context, signUpRequest SignUpRequest) (uuid.UUID, error)

	/*
		GetUserInfo - получение информации о пользователе по UUID пользователя
		Возвращает информацию о пользователе
		Ошибки:
		ErrContextCanceled
		ErrUserNotFound
		"repository: select user with %s: %w"
	*/
	GetUserInfo(ctx context.Context, uuidUser uuid.UUID) (User, error)
}
