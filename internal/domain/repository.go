package domain

import (
	"context"
	uuid "github.com/google/uuid"
)

// GameRepository - хранение текущих игр.
type GameRepository interface {
	/*
		SaveNewGame - сохранение новой игры.
		Возвращает новую игру
		Ошибки:
		ErrContextCanceled
		"repository: failed to insert game: error"
		"repository: failed to select game after insert: error"
		"repository: failed to map game entity to domain: error"
	*/
	SaveNewGame(ctx context.Context, game Game) (Game, error)

	// Save сохраняет игру, перезаписывая существующую с тем же ID.
	Save(ctx context.Context, game Game) error

	// Get срез игр доступных для присоединения игроку с определенным uuid или ошибку
	GetAvaibleGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)

	// Get возвращает игру по ID или ErrGameNotFound.
	Get(ctx context.Context, uuid uuid.UUID) (Game, error)
}

// UserRepository - хранение информации об игроках.
type UserRepository interface {
	// Save сохраняет нового пользователя
	// в случае если пользователь с таким логином уже существует - ErrLoginIncorrect
	Save(ctx context.Context, user User) error

	// Get возвращает данные о пользователе по логину или uuid пользователя
	// В случае если ничего не найдено - ErrUserNotFound
	Get(ctx context.Context, login string, uuidUser uuid.UUID) (User, error)
}
