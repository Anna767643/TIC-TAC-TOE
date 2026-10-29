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
		"repository: failed to insert game: error"
		ErrGameNotFound
		"repository: failed to map game entity to domain: error"
	*/
	SaveNewGame(ctx context.Context, game Game) (Game, error)

	/*
		Save - сохраняет игру
		Ошибки:
		ErrContextCanceled
		"repository: failed to update game: error"
	*/
	Save(ctx context.Context, game Game) error

	/*
		GetAvaibleGames - получение среза доступных игр для присоединения игрока с определенным uuid
		Возвращает срез игр
		Ошибки:
		ErrContextCanceled
		ErrGameNotFound
		"repository: failed to map game entity to domain: error"
	*/
	GetAvaibleGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

	/*
		Get - получает информацию об игре по UUID игры
		Возвращает игру
		Ошибки:
		ErrContextCanceled
		ErrGameNotFound
		"repository: failed to map game entity to domain: error"
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
	*/
	Save(ctx context.Context, user User) error

	/*
		Get - получает данные о пользователе по логину или по UUID
		Возвращает
		Ошибки:
		ErrContextCanceled
		ErrUserNotFound
	*/
	Get(ctx context.Context, login string, uuidUser uuid.UUID) (User, error)
}
