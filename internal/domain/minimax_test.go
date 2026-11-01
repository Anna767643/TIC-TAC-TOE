package domain

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	uuid "github.com/google/uuid"
)

func TestSearchAllOptions(t *testing.T) {

	tests := []struct {
		testName string
		board    Board
		expected []pos
	}{
		{
			testName: "две пустые клетки",
			board: Board{
				{0, X, X},
				{X, X, O},
				{O, O, 0},
			},
			expected: []pos{{0, 0}, {2, 2}},
		},
		{
			testName: "пять пустых клеток",
			board: Board{
				{0, 0, X},
				{X, 0, 0},
				{O, O, 0},
			},
			expected: []pos{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}},
		},
	}

	var gs gameService
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()
			actual := gs.searchAllOptions(tt.board)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("earchAllOptions() got: %v, want: %v.", actual, tt.expected)
			}
		})
	}

}

func TestMinimax(t *testing.T) {
	tests := []struct {
		testName      string
		board         Board
		isMax         bool
		expectedScore int
		expectedErr   error
	}{
		{
			testName: "O выйграли",
			board: Board{
				{0, X, X},
				{X, X, O},
				{O, O, O},
			},
			isMax:         true,
			expectedScore: 10,
		},
		{
			testName: "X выйграли",
			board: Board{
				{0, X, X},
				{O, X, O},
				{X, O, O},
			},
			isMax:         false,
			expectedScore: -10,
		},
		{
			testName: "ничья",
			board: Board{
				{O, X, O},
				{O, X, O},
				{X, O, X},
			},
			isMax:         false,
			expectedScore: 0,
		},
		{
			testName: "O выйграли почти выиграли",
			board: Board{
				{O, O, 0},
				{X, O, X},
				{X, 0, 0},
			},
			isMax:         true,
			expectedScore: 9,
		},
		{
			testName: "X выйграли почти выиграли",
			board: Board{
				{0, O, 0},
				{X, X, 0},
				{X, O, 0},
			},
			isMax:         true,
			expectedScore: -8,
		},
		{
			testName: "Скоро ничья",
			board: Board{
				{O, X, 0},
				{X, O, 0},
				{X, O, 0},
			},
			isMax:         false,
			expectedScore: 0,
		},
	}

	var gs gameService
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			score, err := gs.minimax(ctx, tt.board, 0, tt.isMax)
			if err != nil {
				t.Errorf("Result was incorrect, got: %s, want: nil.", err)
			}
			if tt.expectedScore != score {
				t.Errorf("Result was incorrect, got: %d, want: %d.", score, tt.expectedScore)
			}
		})
	}
}

func TestNextMoveWithComputer(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()

	tests := []struct {
		testName     string
		setupMock    func(*mockGameRepository)
		game         Game
		expectedGame Game
		expectedErr  string
		cancelCtx    bool
	}{
		{
			testName: "ErrContextCanceled",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, X, O},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{},
			expectedErr:  ErrContextCanceled.Error(),
			cancelCtx:    true,
		},

		{
			testName: "Errror to update",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return errors.New("repository: update game")
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, X, O},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{},
			expectedErr:  "repository: update game",
		},
		{
			testName: "Game over, win computer",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{0, X, X},
					{0, O, O},
					{0, 0, X},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, X, X},
					{O, O, O},
					{0, 0, X},
				},
				WithComputer:     true,
				State:            StateWin,
				NextPlayer:       &uuidUser,
				Winner:           nil,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedErr: "",
		},
		{
			testName: "Game over, draw",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{X, 0, X},
					{X, X, O},
					{O, X, O},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{X, O, X},
					{X, X, O},
					{O, X, O},
				},
				WithComputer:     true,
				State:            StateDraw,
				NextPlayer:       &uuidUser,
				Winner:           nil,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedErr: "",
		},
		{
			testName: "Nex move",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, X, O},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, X, O},
					{O, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				Winner:           nil,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedErr: "",
		},
		{
			testName: "Nex move 2",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{0, X, O},
					{0, X, 0},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, X, O},
					{0, X, 0},
					{0, O, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				Winner:           nil,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedErr: "",
		},
		{
			testName: "start game",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
			},
			game: Game{
				UUID: uuidGame,
				Board: Board{
					{X, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{X, 0, 0},
					{0, O, 0},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            StateNextMove,
				NextPlayer:       &uuidUser,
				Winner:           nil,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
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

			gs := &gameService{
				repo: mock,
			}

			game, err := gs.nextMoveWithComputer(ctx, tc.game)

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

func TestNextMove(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	secondUser := uuid.New()

	tests := []struct {
		testName     string
		uuidUser     uuid.UUID
		setupMock    func(*mockGameRepository)
		board        Board
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
			testName: "Errror to select",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{}, errors.New("repository: select game")
				}

			},
			expectedGame: Game{},
			expectedErr:  "repository: select game",
		},
		{
			testName: "Game not found",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{}, ErrGameNotFound
				}

			},
			expectedGame: Game{},
			expectedErr:  ErrGameNotFound.Error(),
		},
		{
			testName: "Wait users",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:     false,
						State:            StateWaitUsers,
						NextPlayer:       nil,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
					}, nil
				}
			},
			expectedErr: ErrWaitOpponent.Error(),
		},
		{
			testName: "Game over, draw",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{X, O, X},
							{X, X, O},
							{O, X, O},
						},
						WithComputer:     true,
						State:            StateDraw,
						NextPlayer:       &uuidUser,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
					}, nil
				}
			},
			expectedErr: ErrGameFinished.Error(),
		},
		{
			testName: "Game over, win",
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, O, X},
							{X, X, O},
							{X, 0, O},
						},
						WithComputer:     true,
						State:            StateWin,
						NextPlayer:       &uuidUser,
						Winner:           &uuidUser,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
					}, nil
				}
			},
			expectedErr: ErrGameFinished.Error(),
		},
		{
			testName: "Wrong uuid",
			uuidUser: uuid.New(),
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, X},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:     true,
						State:            StateNextMove,
						NextPlayer:       &uuidUser,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
					}, nil
				}
			},
			expectedErr: ErrDoesntHaveSuchUser.Error(),
		},
		{
			testName: "Opponent move",
			uuidUser: uuidUser,
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, X},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:     false,
						State:            StateNextMove,
						NextPlayer:       &(secondUser),
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: &(secondUser),
					}, nil
				}
			},
			expectedErr: ErrOpponentsMove.Error(),
		},
		{
			testName: "Wrong board",
			uuidUser: uuidUser,
			board: Board{
				{0, O, X},
				{0, O, 0},
				{0, 0, 0},
			},
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, X},
							{0, O, 0},
							{0, 0, 0},
						},
						WithComputer:     false,
						State:            StateNextMove,
						NextPlayer:       &(uuidUser),
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: &(secondUser),
					}, nil
				}
			},
			expectedErr: ErrInvalidBoard.Error(),
		},
		{
			testName: "Valid board",
			uuidUser: uuidUser,
			board: Board{
				{0, 0, X},
				{0, O, 0},
				{X, O, 0},
			},
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, X},
							{0, O, 0},
							{X, 0, 0},
						},
						WithComputer:     false,
						State:            StateNextMove,
						NextPlayer:       &(uuidUser),
						UUIDPlayerFirst:  secondUser,
						UUIDPlayerSecond: &(uuidUser),
					}, nil
				}
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, O, 0},
					{X, O, 0},
				},
				WithComputer:     false,
				State:            StateNextMove,
				NextPlayer:       &(uuidUser),
				UUIDPlayerFirst:  secondUser,
				UUIDPlayerSecond: &(uuidUser),
			},
			expectedErr: "",
		},
		{
			testName: "Game over",
			uuidUser: uuidUser,
			board: Board{
				{0, 0, X},
				{0, X, O},
				{X, O, 0},
			},
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, 0},
							{0, X, O},
							{X, O, 0},
						},
						WithComputer:     false,
						State:            StateNextMove,
						NextPlayer:       &(uuidUser),
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: &(secondUser),
					}, nil
				}
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, X},
					{0, X, O},
					{X, O, 0},
				},
				WithComputer:     false,
				State:            StateWin,
				Winner:           &(uuidUser),
				NextPlayer:       &(uuidUser),
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: &(secondUser),
			},
			expectedErr: "",
		},
		{
			testName: "Game with computer",
			uuidUser: uuidUser,
			board: Board{
				{0, 0, 0},
				{0, X, O},
				{X, O, 0},
			},
			setupMock: func(m *mockGameRepository) {
				m.SaveFunc = func(ctx context.Context, game Game) error {
					return nil
				}
				m.GetFunc = func(ctx context.Context, uuid uuid.UUID) (Game, error) {
					return Game{
						UUID: uuidGame,
						Board: Board{
							{0, 0, 0},
							{0, 0, O},
							{X, O, 0},
						},
						WithComputer:    true,
						State:           StateNextMove,
						NextPlayer:      &(uuidUser),
						UUIDPlayerFirst: uuidUser,
					}, nil
				}
			},
			expectedGame: Game{
				UUID: uuidGame,
				Board: Board{
					{0, 0, O},
					{0, X, O},
					{X, O, 0},
				},
				WithComputer:    true,
				State:           StateNextMove,
				NextPlayer:      &(uuidUser),
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
			game, err := gs.NextMove(ctx, tc.board, uuidGame, tc.uuidUser)

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
