package datasource

import (
	"context"
	"src/internal/domain"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// mockGameStorage - заглушка DBPool
type mockGameStorage struct {
	CloseFunc                       func()
	GetPoolFunc                     func() *pgxpool.Pool
	insertIntoGamesFunc             func(ctx context.Context, game gameEntity) error
	insertIntoUsersFunc             func(ctx context.Context, user userEntity) error
	selectAvaibleGamesFromGamesFunc func(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error)
	selectFromGamesFunc             func(ctx context.Context, uuid uuid.UUID) (gameEntity, error)
	selectFromUsersWithLoginFunc    func(ctx context.Context, login string) (userEntity, error)
	selectFromUsersWithUUIDFunc     func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error)
	updateGamesFunc                 func(ctx context.Context, game gameEntity) error
}

func (m *mockGameStorage) Close() {}
func (m *mockGameStorage) GetPool() *pgxpool.Pool {
	return m.GetPoolFunc()
}

func (m *mockGameStorage) insertIntoGames(ctx context.Context, game gameEntity) error {
	return m.insertIntoGamesFunc(ctx, game)
}

func (m *mockGameStorage) insertIntoUsers(ctx context.Context, user userEntity) error {
	return m.insertIntoUsersFunc(ctx, user)
}

func (m *mockGameStorage) selectAvaibleGamesFromGames(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
	return m.selectAvaibleGamesFromGamesFunc(ctx, uuidUser)
}

func (m *mockGameStorage) selectFromGames(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
	return m.selectFromGamesFunc(ctx, uuid)
}

func (m *mockGameStorage) selectFromUsersWithLogin(ctx context.Context, login string) (userEntity, error) {
	return m.selectFromUsersWithLoginFunc(ctx, login)
}

func (m *mockGameStorage) selectFromUsersWithUUID(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
	return m.selectFromUsersWithUUID(ctx, uuidUser)
}

func (m *mockGameStorage) updateGames(ctx context.Context, game gameEntity) error {
	return m.updateGamesFunc(ctx, game)
}

func TestSaveNewGame(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	updatedAt := time.Now()
	testCases := []struct {
		testName     string
		setupMock    func(*mockGameStorage)
		game         domain.Game
		expectedGame domain.Game
		expectedErr  error
	}{
		{
			testName: "OK",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{
						UUID:             uuidGame,
						Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
						WithComputer:     true,
						GameState:        domain.WaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			game: domain.Game{},
			expectedGame: domain.Game{
				UUID: uuidGame,
				Board: domain.Board{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:     true,
				GameState:        domain.WaitUsers,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
				UpdatedAt:        updatedAt,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mock := &mockGameStorage{}
			tc.setupMock(mock)

			gr := &gameRepository{
				gs: mock,
			}
			game, err := gr.SaveNewGame(ctx, tc.game)

			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("gr.SaveNewGame() get error %v, but want %v", err, tc.expectedErr)
				return
			}

			if tc.expectedErr == nil {
				if game != tc.expectedGame {
					t.Errorf("Wrong result: get %v, want %v", game, tc.expectedGame)
				}
			}

		})
	}
}

/*
SaveNewGame - сохранение новой игры.
Возвращает новую игру
Ошибки:
ErrContextCanceled
"repository: failed to insert game: error"
"repository: failed to select game after insert: error"
"repository: failed to map game entity to domain: error"
*/
// func (r *gameRepository) SaveNewGame(ctx context.Context, game domain.Game) (domain.Game, error) {
// 	if err := ctx.Err(); err != nil {
// 		return domain.Game{}, domain.ErrContextCanceled
// 	}

// 	gameEntity := toEntityGame(game)

// 	err := r.db.insertIntoGames(ctx, gameEntity)
// 	if err != nil {
// 		return domain.Game{}, fmt.Errorf("repository: failed to insert game: %w", err)
// 	}

// 	gameEntity, err = r.db.selectFromGames(ctx, game.UUID)
// 	if err != nil {
// 		return domain.Game{}, fmt.Errorf("repository: failed to select game after insert: %w", err)
// 	}

// 	newGame, err := toDomainGame(gameEntity)
// 	if err != nil {
// 		return domain.Game{}, fmt.Errorf("repository: failed to map game entity to domain: %w", err)
// 	}

// 	return newGame, err
// }
