package domain

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	uuid "github.com/google/uuid"
)

type mockGameRepository struct {
	SaveNewGameFunc       func(ctx context.Context, game Game) (Game, error)
	SaveFunc              func(ctx context.Context, game Game) error
	GetAvailableGamesFunc func(ctx context.Context, uuidUser uuid.UUID) ([]Game, error)
	GetFunc               func(ctx context.Context, uuid uuid.UUID) (Game, error)
}

func (m *mockGameRepository) SaveNewGame(ctx context.Context, game Game) (Game, error) {
	return m.SaveNewGameFunc(ctx, game)
}

func (m *mockGameRepository) Save(ctx context.Context, game Game) error {
	return m.SaveFunc(ctx, game)
}

func (m *mockGameRepository) GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]Game, error) {
	return m.GetAvailableGamesFunc(ctx, uuidUser)
}

func (m *mockGameRepository) Get(ctx context.Context, uuid uuid.UUID) (Game, error) {
	return m.GetFunc(ctx, uuid)
}

func TestCreateNewGame(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	updatedAt := time.Now()

	tests := []struct {
		testName     string
		setupMock    func(*mockGameRepository)
		expectedGame Game
		expectedErr  string
		cancelCtx    bool
		withComputer bool
	}{
		{
			testName: "ErrContextCanceled",
			setupMock: func(m *mockGameRepository) {
			},
			expectedGame: Game{},
			expectedErr:  ErrContextCanceled.Error(),
			cancelCtx:    true,
		},
		{
			testName: "Errror save new game",
			setupMock: func(m *mockGameRepository) {
				m.SaveNewGameFunc = func(ctx context.Context, game Game) (Game, error) {
					return Game{}, errors.New("repository: insert game %s: %w")
				}
			},
			expectedGame: Game{},
			expectedErr:  "repository: insert game %s: %w",
		},
		{
			testName: "OK",
			setupMock: func(m *mockGameRepository) {
				m.SaveNewGameFunc = func(ctx context.Context, game Game) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:    true,
						State:           StateNextMove,
						NextPlayer:      &(uuidUser),
						UUIDPlayerFirst: uuidUser,
						UpdatedAt:       updatedAt,
					}, nil
				}
			},
			withComputer: true,
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:    true,
				State:           StateNextMove,
				NextPlayer:      &(uuidUser),
				UUIDPlayerFirst: uuidUser,
				UpdatedAt:       updatedAt,
			},
			expectedErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())

			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

			mock := &mockGameRepository{}
			tc.setupMock(mock)

			gs := NewGameService(mock)
			game, err := gs.CreateNewGame(ctx, tc.withComputer, uuidUser)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				} else {
					if game.Board != tc.expectedGame.Board {
						t.Errorf("Wrong result: get %v, want %v", game.Board, tc.expectedGame.Board)
					}
					if game.State != tc.expectedGame.State {
						t.Errorf("Wrong result: get %s, want %s", game.State, tc.expectedGame.State)
					}
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedErr)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedErr) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedErr, err)
				} else {
					if game.Board != tc.expectedGame.Board {
						t.Errorf("Wrong result: get %v, want %v", game.Board, tc.expectedGame.Board)
					}

					if game.State != tc.expectedGame.State {
						t.Errorf("Wrong result: get %s, want %s", game.State, tc.expectedGame.State)
					}
				}

			}
		})
	}
}

func TestGetAvailableGames(t *testing.T) {
	uuidUser := uuid.New()

	tests := []struct {
		testName      string
		setupMock     func(*mockGameRepository)
		expectedGames []Game
		expectedErr   string
		cancelCtx     bool
		withComputer  bool
	}{
		{
			testName: "ErrContextCanceled",
			setupMock: func(m *mockGameRepository) {
			},
			expectedGames: []Game{},
			expectedErr:   ErrContextCanceled.Error(),
			cancelCtx:     true,
		},
		{
			testName: "Errror get games",
			setupMock: func(m *mockGameRepository) {
				m.GetAvailableGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]Game, error) {
					return []Game{}, errors.New("repository: select game for user %s: %w")
				}
			},
			expectedGames: []Game{},
			expectedErr:   "repository: select game for user %s: %w",
		},
		{
			testName: "OK",
			setupMock: func(m *mockGameRepository) {
				m.GetAvailableGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]Game, error) {
					return []Game{
						{
							Board: Board{
								{0, 0, 0},
								{0, 0, 0},
								{0, 0, 0},
							},
							WithComputer:    false,
							State:           StateNextMove,
							NextPlayer:      &(uuidUser),
							UUIDPlayerFirst: uuidUser,
						},
					}, nil
				}
			},
			withComputer: true,
			expectedGames: []Game{
				{
					Board: Board{
						{0, 0, 0},
						{0, 0, 0},
						{0, 0, 0},
					},
					WithComputer:    false,
					State:           StateNextMove,
					NextPlayer:      &(uuidUser),
					UUIDPlayerFirst: uuidUser,
				},
			},
			expectedErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())

			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

			mock := &mockGameRepository{}
			tc.setupMock(mock)

			gs := NewGameService(mock)
			games, err := gs.GetAvailableGames(ctx, uuidUser)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				} else {
					if games[0].Board != tc.expectedGames[0].Board {
						t.Errorf("Wrong result: get %v, want %v", games[0].Board, tc.expectedGames[0].Board)
					}
					if games[0].State != tc.expectedGames[0].State {
						t.Errorf("Wrong result: get %s, want %s", games[0].State, tc.expectedGames[0].State)
					}
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedErr)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedErr) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedErr, err)
				}
			}
		})
	}
}

func TestGetGame(t *testing.T) {
	uuidUser := uuid.New()
	uuidGame := uuid.New()

	secondUser := uuid.New()

	tests := []struct {
		testName     string
		setupMock    func(*mockGameRepository)
		expectedGame Game
		expectedErr  string
		cancelCtx    bool
	}{
		{
			testName: "ErrContextCanceled",
			setupMock: func(m *mockGameRepository) {
			},
			expectedGame: Game{},
			expectedErr:  ErrContextCanceled.Error(),
			cancelCtx:    true,
		},
		{
			testName: "ErrGameNotFound",
			setupMock: func(m *mockGameRepository) {
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{}, ErrGameNotFound
				}
			},
			expectedGame: Game{},
			expectedErr:  ErrGameNotFound.Error(),
		},
		{
			testName: "Error get game",
			setupMock: func(m *mockGameRepository) {
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{}, errors.New("repository: select game")
				}
			},
			expectedGame: Game{},
			expectedErr:  "repository: select game",
		},
		{
			testName: "ErrDoesntHaveSuchUser",
			setupMock: func(m *mockGameRepository) {
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUIDPlayerFirst:  uuidGame,
						UUIDPlayerSecond: &(secondUser),
					}, nil
				}
			},
			expectedGame: Game{},
			expectedErr:  ErrDoesntHaveSuchUser.Error(),
		},
		{
			testName: "OK",
			setupMock: func(m *mockGameRepository) {
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						Board: Board{
							{0, 0, 0},
							{0, 1, 0},
							{0, 0, 0},
						},
						UUIDPlayerFirst: uuidUser,
					}, nil
				}
			},
			expectedGame: Game{
				Board: Board{
					{0, 0, 0},
					{0, 1, 0},
					{0, 0, 0},
				},
				UUIDPlayerFirst: uuidUser,
			},
			expectedErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())

			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

			mock := &mockGameRepository{}
			tc.setupMock(mock)

			gs := NewGameService(mock)
			game, err := gs.GetGame(ctx, uuidGame, uuidUser)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				} else {
					if game.Board != tc.expectedGame.Board {
						t.Errorf("Wrong result: get %v, want %v", game.Board, tc.expectedGame.Board)
					}
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedErr)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedErr) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedErr, err)
				}
			}
		})
	}
}
