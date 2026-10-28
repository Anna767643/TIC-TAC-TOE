package domain

// import (
// 	"context"
// 	// "errors"
// 	"reflect"
// 	"testing"
// )

// func TestSearchAllOptions(t *testing.T) {

// 	tests := []struct {
// 		testName string
// 		board    Board
// 		expected []pos
// 	}{
// 		{
// 			testName: "две пустые клетки",
// 			board: Board{
// 				{0, X, X},
// 				{X, X, O},
// 				{O, O, 0},
// 			},
// 			expected: []pos{{0, 0}, {2, 2}},
// 		},
// 		{
// 			testName: "пять пустых клеток",
// 			board: Board{
// 				{0, 0, X},
// 				{X, 0, 0},
// 				{O, O, 0},
// 			},
// 			expected: []pos{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}},
// 		},
// 	}

// 	var gs gameService
// 	for _, tt := range tests {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			actual := gs.searchAllOptions(tt.board)
// 			if !reflect.DeepEqual(actual, tt.expected) {
// 				t.Errorf("earchAllOptions() got: %v, want: %v.", actual, tt.expected)
// 			}
// 			tt.board.Print()
// 		})
// 	}

// }

// func TestMinimax(t *testing.T) {
// 	tests := []struct {
// 		testName      string
// 		board         Board
// 		isMax         bool
// 		expectedScore int
// 		expectedErr   error
// 	}{
// 		{
// 			testName: "O выйграли",
// 			board: Board{
// 				{0, X, X},
// 				{X, X, O},
// 				{O, O, O},
// 			},
// 			isMax:         true,
// 			expectedScore: 10,
// 		},
// 		{
// 			testName: "X выйграли",
// 			board: Board{
// 				{0, X, X},
// 				{O, X, O},
// 				{X, O, O},
// 			},
// 			isMax:         false,
// 			expectedScore: -10,
// 		},
// 		{
// 			testName: "ничья",
// 			board: Board{
// 				{O, X, O},
// 				{O, X, O},
// 				{X, O, X},
// 			},
// 			isMax:         false,
// 			expectedScore: 0,
// 		},
// 		{
// 			testName: "O выйграли почти выиграли",
// 			board: Board{
// 				{O, O, 0},
// 				{X, O, X},
// 				{X, 0, 0},
// 			},
// 			isMax:         true,
// 			expectedScore: 9,
// 		},
// 		{
// 			testName: "X выйграли почти выиграли",
// 			board: Board{
// 				{0, O, 0},
// 				{X, X, 0},
// 				{X, O, 0},
// 			},
// 			isMax:         true,
// 			expectedScore: -8,
// 		},
// 		{
// 			testName: "Скоро ничья",
// 			board: Board{
// 				{O, X, 0},
// 				{X, O, 0},
// 				{X, O, 0},
// 			},
// 			isMax:         false,
// 			expectedScore: 0,
// 		},
// 	}

// 	var gs gameService
// 	for _, tt := range tests {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			defer cancel()

// 			score, err := gs.minimax(ctx, tt.board, 0, tt.isMax)
// 			if err != nil {
// 				t.Errorf("Result was incorrect, got: %s, want: nil.", err)
// 			}
// 			if tt.expectedScore != score {
// 				t.Errorf("Result was incorrect, got: %d, want: %d.", score, tt.expectedScore)
// 			}
// 			tt.board.Print()
// 		})
// 	}
// }

// func TestNextMove(t *testing.T) {
// 	tests := []struct {
// 		testName      string
// 		board         Board
// 		expectedBoard Board
// 		expectedErr   error
// 		cancelCtx     bool
// 	}{
// 		{
// 			testName: "Закрываем контекст",
// 			board: Board{
// 				{0, 0, X},
// 				{0, X, O},
// 				{0, 0, 0},
// 			},
// 			expectedErr: ErrContextCanceled,
// 			cancelCtx:   true,
// 		},
// 		{
// 			testName: "Игра окончена",
// 			board: Board{
// 				{X, O, X},
// 				{X, X, O},
// 				{O, X, O},
// 			},
// 			expectedBoard: Board{
// 				{X, O, X},
// 				{X, X, O},
// 				{O, X, O},
// 			},
// 			expectedErr: ErrGameFinished,
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш",
// 			board: Board{
// 				{0, 0, X},
// 				{0, X, O},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{0, 0, X},
// 				{0, X, O},
// 				{O, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш2",
// 			board: Board{
// 				{O, 0, X},
// 				{X, X, 0},
// 				{O, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{O, 0, X},
// 				{X, X, O},
// 				{O, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш3",
// 			board: Board{
// 				{X, 0, X},
// 				{O, X, 0},
// 				{O, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{X, O, X},
// 				{O, X, 0},
// 				{O, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш4",
// 			board: Board{
// 				{0, 0, X},
// 				{O, X, 0},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{0, 0, X},
// 				{O, X, 0},
// 				{O, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш5",
// 			board: Board{
// 				{0, 0, O},
// 				{X, X, 0},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{0, 0, O},
// 				{X, X, O},
// 				{0, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш6",
// 			board: Board{
// 				{0, 0, O},
// 				{0, X, 0},
// 				{0, 0, X},
// 			},
// 			expectedBoard: Board{
// 				{O, 0, O},
// 				{0, X, 0},
// 				{0, 0, X},
// 			},
// 		},
// 		{
// 			testName: "Предотвращаем пройгрыш7",
// 			board: Board{
// 				{0, X, O},
// 				{0, X, 0},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{0, X, O},
// 				{0, X, 0},
// 				{0, O, 0},
// 			},
// 		},
// 		{
// 			testName: "Начало игры",
// 			board: Board{
// 				{0, X, 0},
// 				{0, 0, 0},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{O, X, 0},
// 				{0, 0, 0},
// 				{0, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Начало игры2",
// 			board: Board{
// 				{X, 0, 0},
// 				{0, 0, 0},
// 				{0, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{X, 0, 0},
// 				{0, O, 0},
// 				{0, 0, 0},
// 			},
// 		},
// 		{
// 			testName: "Боремся до конца",
// 			board: Board{
// 				{0, X, X},
// 				{0, X, O},
// 				{0, 0, O},
// 			},
// 			expectedBoard: Board{
// 				{O, X, X},
// 				{0, X, O},
// 				{0, 0, O},
// 			},
// 		},
// 		{
// 			testName: "Выигрываем",
// 			board: Board{
// 				{X, X, O},
// 				{X, 0, O},
// 				{O, 0, 0},
// 			},
// 			expectedBoard: Board{
// 				{X, X, O},
// 				{X, O, O},
// 				{O, 0, 0},
// 			},
// 			expectedErr: ErrGameFinished,
// 		},
// 	}

// 	var gs gameService
// 	for _, tt := range tests {
// 		t.Run(tt.testName, func(t *testing.T) {
// 			ctx, cancel := context.WithCancel(context.Background())

// 			if !tt.cancelCtx {
// 				defer cancel()
// 			} else {
// 				cancel()
// 			}

// 			game := Game{Board: tt.board}
// 			actualGame, err := gs.NextMove(ctx, game)

// 			if !errors.Is(err, tt.expectedErr) {
// 				t.Errorf("NextMove() error = %v, want %v", err, tt.expectedErr)
// 				return
// 			}

// 			if tt.expectedErr != ErrContextCanceled {
// 				if !reflect.DeepEqual(actualGame.Board, tt.expectedBoard) {
// 					t.Errorf("NextMove() board mismatch:\n got:\n %v\n want:\n %v", actualGame.Board, tt.expectedBoard)
// 				}
// 			}
// 			actualGame.Board.Print()
// 		})
// 	}
// }
