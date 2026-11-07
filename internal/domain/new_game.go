package domain

import (
	"context"
	"time"

	uuid "github.com/google/uuid"
)

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
func (s *gameService) CreateNewGame(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	game := Game{
		UUID: uuid.New(),
		Board: Board{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		WithComputer:    withComputer,
		State:           StateWaitUsers,
		UUIDPlayerFirst: uuidUser,

		UpdatedAt: time.Now(),
	}

	if withComputer {
		game.State = StateNextMove
		game.NextPlayer = &uuidUser
	}

	newGame, err := s.repo.SaveNewGame(ctx, game)
	if err != nil {
		return Game{}, err
	}
	return newGame, nil
}

/*
GetAvailableGames - находит и возвращает все доступные игры
для присоединения игрока с определённым UUID, или ошибку
Ошибки:
ErrContextCanceled
"repository: select game for user %s: %w"
"repository: scan game for user %s: %w"
*/
func (s *gameService) GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error) {
	if err := ctx.Err(); err != nil {
		return []Game{}, ErrContextCanceled
	}

	games, err := s.repo.GetAvailableGames(ctx, uuidUser)
	if err != nil {
		return []Game{}, err
	}

	return games, nil
}

/*
JoinUser - присоединяет пользователя к игре
В случае успеха возвращает измененную игру
Ошибки:
ErrContextCanceled
"repository: select game %s: %w"
ErrAddUser
"repository: update user in game %s: %w"
*/
func (s *gameService) JoinUser(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	game, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	game.State = StateNextMove
	game.NextPlayer = &game.UUIDPlayerFirst
	game.UUIDPlayerSecond = &uuidUser
	game.UpdatedAt = time.Now()

	if err := s.repo.JoinSecondUser(ctx, game); err != nil {
		return Game{}, err
	}

	newGame, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	return newGame, nil
}

/*
GetGame - получает информацию об игре
Ошибки:
ErrContextCanceled - игра прерывается
ErrGameNotFound - игра не найдена
"repository: select game %s: %w"
ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
*/
func (s *gameService) GetGame(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	game, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	if !game.HasPlayer(uuidUser) {
		return Game{}, ErrDoesntHaveSuchUser
	}

	return game, nil
}
