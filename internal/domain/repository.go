package domain

import (
	"context"
	uuid "github.com/google/uuid"
)

// GameRepository - хранение текущих игр
type GameRepository interface {
	/*
		SaveNewGame - сохранение новой игры
		Возвращает новую игру
		Ошибки:
		ErrContextCanceled
		"repository: insert game %s: %w"
		ErrGameNotFound
		"repository: select game %s: %w"
		"repository: failed to map game entity to domain: error"
	*/
	SaveNewGame(ctx context.Context, game Game) (Game, error)

	/*
		Save - сохраняет игру
		Ошибки:
		ErrContextCanceled
		"repository: update game %s: %w"
	*/
	Save(ctx context.Context, game Game) error

	/*
		GetAvailableGames - получение среза доступных игр для присоединения игрока с определенным uuid
		Возвращает срез игр
		Ошибки:
		ErrContextCanceled
		"repository: select game for user %s: %w"
		"repository: scan game for user %s: %w"
	*/
	GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

	/*
		Get - получает информацию об игре по UUID игры
		Возвращает игру
		Ошибки:
		ErrContextCanceled
		ErrGameNotFound
		"repository: select game %s: %w"
	*/
	Get(ctx context.Context, uuid uuid.UUID) (Game, error)
}

// UserRepository - хранение информации об игроках
type UserRepository interface {
	/*
		Save - сохраняет нового пользователя
		Ошибки:
		ErrContextCanceled
		ErrLoginIncorrect
		"repository: insert user with %s: %w"
	*/
	Save(ctx context.Context, user User) error

	/*
		Get - получает данные о пользователе по логину или по UUID
		Возвращает
		Ошибки:
		ErrContextCanceled
		ErrUserNotFound
		"repository: select user with %s: %w"
	*/
	Get(ctx context.Context, login string, uuidUser uuid.UUID) (User, error)
}
