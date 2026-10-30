package datasource

import (
	"context"
	"fmt"
	"src/internal/domain"

	uuid "github.com/google/uuid"
)

// gameRepository - реализация domain.GameRepository поверх gameStorage
type gameRepository struct {
	gs gameStorage
}

// NewGameRepository - конструктор для di
func NewGameRepository(db *DBPool) domain.GameRepository {
	return &gameRepository{gs: db}
}

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
func (r *gameRepository) SaveNewGame(ctx context.Context, game domain.Game) (domain.Game, error) {
	if err := ctx.Err(); err != nil {
		return domain.Game{}, domain.ErrContextCanceled
	}

	gameEntity := toEntityGame(game)

	err := r.gs.insertIntoGames(ctx, gameEntity)
	if err != nil {
		return domain.Game{}, err
	}

	gameEntity, err = r.gs.selectFromGames(ctx, game.UUID)
	if err != nil {
		return domain.Game{}, err
	}

	newGame, err := toDomainGame(gameEntity)
	if err != nil {
		return domain.Game{}, fmt.Errorf("repository: failed to map game entity to domain: %w", err)
	}

	return newGame, err
}

/*
Save - сохраняет игру
Ошибки:
ErrContextCanceled
"repository: update game %s: %w"
*/
func (r *gameRepository) Save(ctx context.Context, game domain.Game) error {
	if err := ctx.Err(); err != nil {
		return domain.ErrContextCanceled
	}

	gameEntity := toEntityGame(game)

	if err := r.gs.updateGames(ctx, gameEntity); err != nil {
		return err
	}

	return nil
}

/*
GetAvailableGames - получение среза доступных игр для присоединения игрока с определенным uuid
Возвращает срез игр
Ошибки:
ErrContextCanceled
"repository: select game for user %s: %w"
"repository: scan game for user %s: %w"
*/
func (r *gameRepository) GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error) {
	if err := ctx.Err(); err != nil {
		return []domain.Game{}, domain.ErrContextCanceled
	}

	entityGames, err := r.gs.selectAvaibleGamesFromGames(ctx, uuidUser)
	if err != nil {
		return []domain.Game{}, err
	}

	domainGames := make([]domain.Game, 0)
	for i := range entityGames {
		domainGame, err := toDomainGame(entityGames[i])
		if err != nil {
			return []domain.Game{}, fmt.Errorf("repository: failed to map game entity to domain: %w", err)
		}
		domainGames = append(domainGames, domainGame)
	}

	return domainGames, nil
}

/*
Get - получает информацию об игре по UUID игры
Возвращает игру
Ошибки:
ErrContextCanceled
ErrGameNotFound
"repository: select game %s: %w"
*/
func (r *gameRepository) Get(ctx context.Context, uuid uuid.UUID) (domain.Game, error) {
	if err := ctx.Err(); err != nil {
		return domain.Game{}, domain.ErrContextCanceled
	}

	var (
		entity gameEntity
		err    error
	)

	entity, err = r.gs.selectFromGames(ctx, uuid)
	if err != nil {
		return domain.Game{}, err
	}

	v, err := toDomainGame(entity)
	if err != nil {
		return domain.Game{}, fmt.Errorf("repository: failed to map game entity to domain: %w", err)
	}

	return v, nil
}

// gameRepository - реализация domain.GameRepository поверх gameStorage
type userRepository struct {
	gs gameStorage
}

// NewUserRepository - конструктор для di
func NewUserRepository(db *DBPool) domain.UserRepository {
	return &userRepository{gs: db}
}

/*
Save - сохраняет нового пользователя
Ошибки:
ErrContextCanceled
ErrLoginIncorrect
"repository: insert user with %s: %w"
*/
func (u *userRepository) Save(ctx context.Context, user domain.User) error {
	if err := ctx.Err(); err != nil {
		return domain.ErrContextCanceled
	}

	userEntity := toUserEntity(user)

	err := u.gs.insertIntoUsers(ctx, userEntity)
	if err != nil {
		return err
	}

	return nil
}

/*
Get - получает данные о пользователе по логину или по UUID
Возвращает
Ошибки:
ErrContextCanceled
ErrUserNotFound
"repository: select user with %s: %w"
*/
func (u *userRepository) Get(ctx context.Context, login string, uuidUser uuid.UUID) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, domain.ErrContextCanceled
	}

	var (
		userEntity userEntity
		err        error
	)

	if login != "" {
		userEntity, err = u.gs.selectFromUsersWithLogin(ctx, login)
		if err != nil {
			return domain.User{}, err
		}
	} else {
		userEntity, err = u.gs.selectFromUsersWithUUID(ctx, uuidUser)
		if err != nil {
			return domain.User{}, err
		}
	}

	domainUser := toDomainUser(userEntity)

	return domainUser, nil
}
