package domain

import (
	"context"
	"time"

	uuid "github.com/google/uuid"
)

// CreateNewGame - создаёт новую игру с компьютером или другим игроком.
// Возвращает информацию о новой игре и ошибку
func (s *gameService) CreateNewGame(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	gameState := NextMove + uuidUser.String()
	if !withComputer {
		gameState = WaitUsers
	}

	game := Game{
		UUID: uuid.New(),
		Board: Board{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
		WithComputer:    withComputer,
		GameState:       gameState,
		UUIDPlayerFirst: uuidUser,

		UpdatedAt: time.Now(),
	}

	newGame, err := s.repo.SaveNewGame(ctx, game)
	if err != nil {
		return Game{}, err
	}
	return newGame, nil
}

// GetAvaibleGames - находит и возвращает все доступные игры
// для присоединения игрока с определённым UUID, или ошибку
func (s *gameService) GetAvaibleGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error) {
	if err := ctx.Err(); err != nil {
		return []Game{}, ErrContextCanceled
	}

	games, err := s.repo.GetAvaibleGames(ctx, uuidUser)
	if err != nil {
		return []Game{}, err
	}

	return games, err
}

// JoinUser - присоединяет пользователя к игре
// В случае успеха возвращает измененную игру
// Ошибки: игра не найдена, к данной игре нельзя добавить пользователя,
// context.Canceled, или иная неописанная ошибка
func (s *gameService) JoinUser(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	game, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	// ErrAddUser
	if game.GameState != WaitUsers ||
		game.UUIDPlayerSecond != nil ||
		game.WithComputer != false ||
		game.UUIDPlayerFirst == uuidUser {
		return Game{}, ErrAddUser
	}

	game.GameState = NextMove + game.UUIDPlayerFirst.String()
	game.UUIDPlayerSecond = &uuidUser
	game.UpdatedAt = time.Now()

	// не обновляется если второй игрок с таким же ID
	if err := s.repo.Save(ctx, game); err != nil {
		return Game{}, ErrAddUser
	}

	newGame, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	return newGame, nil
}

// GetGame - получает информацию об игре
// Ошибки:
// ErrContextCanceled - игра прерывается
// ErrGameNotFound - игра не найдена
// ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
func (s *gameService) GetGame(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	game, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	if game.UUIDPlayerFirst != uuidUser && *game.UUIDPlayerSecond != uuidUser {
		return Game{}, ErrDoesntHaveSuchUser
	}

	return game, nil
}
