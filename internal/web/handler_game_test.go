package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"src/internal/domain"
	"testing"
	"time"

	uuid "github.com/google/uuid"
)

// mockGameService - заглушка GameService,
type mockGameService struct {
	CreateNewGameFunc     func(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error)
	GetAvailableGamesFunc func(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error)
	JoinUserFunc          func(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error)
	NextMoveFunc          func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error)
	GetGameFunc           func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error)
}

func (m *mockGameService) CreateNewGame(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error) {
	return m.CreateNewGameFunc(ctx, withComputer, uuidUser)
}

func (m *mockGameService) GetAvailableGames(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error) {
	return m.GetAvailableGamesFunc(ctx, uuidUser)
}

func (m *mockGameService) JoinUser(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error) {
	return m.JoinUserFunc(ctx, uuidUser, uuidGame)
}

func (m *mockGameService) NextMove(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
	return m.NextMoveFunc(ctx, board, uuidGame, uuidUser)
}

func (m *mockGameService) GetGame(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
	return m.GetGameFunc(ctx, uuidGame, uuidUser)
}

func TestHandleCreateGame(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	createAt := time.Now().UTC()

	testCases := []struct {
		name           string
		setupmocks     func(*mockGameService)
		requestBody    any
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			setupmocks: func(m *mockGameService) {
				m.CreateNewGameFunc = func(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{
						ID:   144,
						UUID: uuidGame,
						Board: domain.Board{
							{0, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:    withComputer,
						State:           domain.StateNextMove,
						NextPlayer:      &uuidUser,
						UUIDPlayerFirst: uuidUser,
						UpdatedAt:       createAt,
					}, nil
				}
			},
			requestBody: struct {
				WithComputer bool `json:"with_computer"`
			}{WithComputer: true},
			expectedCode:   http.StatusCreated,
			expectedInBody: fmt.Sprintf(`{"id":144,"uuid":"%s","board":[[0,0,0],[0,0,0],[0,0,0]],"with_computer":true,"state":"next move","next_player":"%s","winner":null,"uuid_player_first":"%s","uuid_player_second":null,"update_at":"%s"}`, uuidGame.String(), uuidUser.String(), uuidUser.String(), createAt.Format(time.RFC3339)),
		},
		{
			name: "Empty json",
			setupmocks: func(m *mockGameService) {
				m.CreateNewGameFunc = func(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			requestBody:    nil,
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"wrong body, need param with_computer, error: EOF"}`,
		},
		{
			name: "don't have with_computer",
			setupmocks: func(m *mockGameService) {
				m.CreateNewGameFunc = func(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			requestBody:    struct{}{},
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"missing required field: with_computer"}`,
		},
		{
			name: "ErrContextCanceled",
			setupmocks: func(m *mockGameService) {
				m.CreateNewGameFunc = func(ctx context.Context, withComputer bool, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrContextCanceled
				}
			},
			requestBody: struct {
				WithComputer bool `json:"with_computer"`
			}{WithComputer: true},
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: `{"error":"context canceled"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGS := &mockGameService{}
			tc.setupmocks(mockGS)

			handler := NewGameHandler(mockGS)

			var bodyBytes []byte

			if tc.name != "Empty json" {
				bodyBytes, _ = json.Marshal(tc.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/games", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			newCtx := context.WithValue(context.Background(), userUUIDKey, uuidUser)
			req = req.WithContext(newCtx)

			rec := httptest.NewRecorder()

			handler.handleCreateGame(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}

func TestHandleGetAvailableGames(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	createAt := time.Now().UTC()

	testCases := []struct {
		name           string
		setupmocks     func(*mockGameService)
		queryParams    string
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			setupmocks: func(m *mockGameService) {
				m.GetAvailableGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error) {
					return []domain.Game{
						{
							ID:   144,
							UUID: uuidGame,
							Board: domain.Board{
								{0, 0, 0},
								{0, 0, 0},
								{0, 0, 0},
							},
							WithComputer:    false,
							State:           domain.StateWaitUsers,
							NextPlayer:      nil,
							UUIDPlayerFirst: uuidUser,
							UpdatedAt:       createAt,
						},
					}, nil
				}
			},
			queryParams:    "?gameState=waitUsers",
			expectedCode:   http.StatusOK,
			expectedInBody: fmt.Sprintf(`"uuid":"%s"`, uuidGame.String()),
		},
		{
			name: "Wrong query param",
			setupmocks: func(m *mockGameService) {
				m.GetAvailableGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error) {
					return nil, nil
				}
			},
			queryParams:    "?gameState=active",
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"need query param gameState"}`,
		},
		{
			name: "Service error",
			setupmocks: func(m *mockGameService) {
				m.GetAvailableGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]domain.Game, error) {
					return nil, domain.ErrContextCanceled
				}
			},
			queryParams:    "?gameState=waitUsers",
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: `{"error":"context canceled"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGS := &mockGameService{}
			tc.setupmocks(mockGS)

			handler := NewGameHandler(mockGS)

			req := httptest.NewRequest(http.MethodGet, "/games"+tc.queryParams, nil)
			newCtx := context.WithValue(context.Background(), userUUIDKey, uuidUser)
			req = req.WithContext(newCtx)

			rec := httptest.NewRecorder()

			handler.handleGetAvailableGames(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}

func TestHandleJoinUser(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	createAt := time.Now().UTC()

	testCases := []struct {
		name           string
		setupmocks     func(*mockGameService)
		gameID         string
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			setupmocks: func(m *mockGameService) {
				m.JoinUserFunc = func(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error) {
					return domain.Game{
						ID:   144,
						UUID: uuidGame,
						Board: domain.Board{
							{0, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:    false,
						State:           domain.StateNextMove,
						NextPlayer:      &uuidUser,
						UUIDPlayerFirst: uuidUser,
						UpdatedAt:       createAt,
					}, nil
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusOK,
			expectedInBody: fmt.Sprintf(`"uuid":"%s"`, uuidGame.String()),
		},
		{
			name: "Invalid UUID",
			setupmocks: func(m *mockGameService) {
				m.JoinUserFunc = func(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			gameID:         "invalid-uuid",
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"`,
		},
		{
			name: "ErrGameNotFound",
			setupmocks: func(m *mockGameService) {
				m.JoinUserFunc = func(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrGameNotFound
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusNotFound,
			expectedInBody: `{"error":"error to add user: `,
		},
		{
			name: "ErrAddUser",
			setupmocks: func(m *mockGameService) {
				m.JoinUserFunc = func(ctx context.Context, uuidUser uuid.UUID, uuidGame uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrAddUser
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusConflict,
			expectedInBody: `{"error":"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGS := &mockGameService{}
			tc.setupmocks(mockGS)

			handler := NewGameHandler(mockGS)

			req := httptest.NewRequest(http.MethodPatch, "/games/"+tc.gameID+"/users", nil)
			newCtx := context.WithValue(context.Background(), userUUIDKey, uuidUser)
			req = req.WithContext(newCtx)

			req.SetPathValue("uuid", tc.gameID)

			rec := httptest.NewRecorder()

			handler.handleJoinUser(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}

func TestHandlePatchCurrentGame(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	createAt := time.Now().UTC()

	testCases := []struct {
		name           string
		setupmocks     func(*mockGameService)
		gameID         string
		requestBody    any
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{
						ID:   144,
						UUID: uuidGame,
						Board: domain.Board{
							{1, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:    false,
						State:           domain.StateNextMove,
						NextPlayer:      &uuidUser,
						UUIDPlayerFirst: uuidUser,
						UpdatedAt:       createAt,
					}, nil
				}
			},
			gameID: uuidGame.String(),
			requestBody: map[string]any{
				"board": [][]int{
					{1, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
			},
			expectedCode:   http.StatusOK,
			expectedInBody: fmt.Sprintf(`"uuid":"%s"`, uuidGame.String()),
		},
		{
			name: "Invalid UUID",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			gameID:         "invalid-uuid",
			requestBody:    nil,
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"`,
		},
		{
			name: "Empty json",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			gameID:         uuidGame.String(),
			requestBody:    nil,
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"wrong body EOF"}`,
		},
		{
			name: "ErrGameNotFound",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrGameNotFound
				}
			},
			gameID: uuidGame.String(),
			requestBody: map[string]any{
				"board": [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			},
			expectedCode:   http.StatusNotFound,
			expectedInBody: `{"error":"`,
		},
		{
			name: "ErrInvalidBoard",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrInvalidBoard
				}
			},
			gameID: uuidGame.String(),
			requestBody: map[string]any{
				"board": [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			},
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"wrong body `,
		},
		{
			name: "ErrOpponentsMove",
			setupmocks: func(m *mockGameService) {
				m.NextMoveFunc = func(ctx context.Context, board domain.Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrOpponentsMove
				}
			},
			gameID: uuidGame.String(),
			requestBody: map[string]any{
				"board": [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			},
			expectedCode:   http.StatusConflict,
			expectedInBody: `{"error":"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGS := &mockGameService{}
			tc.setupmocks(mockGS)

			handler := NewGameHandler(mockGS)

			var bodyBytes []byte
			if tc.requestBody != nil {
				bodyBytes, _ = json.Marshal(tc.requestBody)
			}

			req := httptest.NewRequest(http.MethodPatch, "/games/"+tc.gameID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			newCtx := context.WithValue(context.Background(), userUUIDKey, uuidUser)
			req = req.WithContext(newCtx)

			req.SetPathValue("uuid", tc.gameID)

			rec := httptest.NewRecorder()

			handler.handlePatchCurrentGame(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}

func TestHandleGetGameInfo(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	createAt := time.Now().UTC()

	testCases := []struct {
		name           string
		setupmocks     func(*mockGameService)
		gameID         string
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			setupmocks: func(m *mockGameService) {
				m.GetGameFunc = func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{
						ID:   144,
						UUID: uuidGame,
						Board: domain.Board{
							{0, 0, 0},
							{0, 0, 0},
							{0, 0, 0},
						},
						WithComputer:    false,
						State:           domain.StateNextMove,
						NextPlayer:      &uuidUser,
						UUIDPlayerFirst: uuidUser,
						UpdatedAt:       createAt,
					}, nil
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusOK,
			expectedInBody: fmt.Sprintf(`"uuid":"%s"`, uuidGame.String()),
		},
		{
			name: "Invalid UUID",
			setupmocks: func(m *mockGameService) {
				m.GetGameFunc = func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, nil
				}
			},
			gameID:         "invalid-uuid",
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"`,
		},
		{
			name: "ErrGameNotFound",
			setupmocks: func(m *mockGameService) {
				m.GetGameFunc = func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrGameNotFound
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusNotFound,
			expectedInBody: `{"error":"`,
		},
		{
			name: "ErrDoesntHaveSuchUser",
			setupmocks: func(m *mockGameService) {
				m.GetGameFunc = func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrDoesntHaveSuchUser
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `{"error":"wrong body `,
		},
		{
			name: "Service Error",
			setupmocks: func(m *mockGameService) {
				m.GetGameFunc = func(ctx context.Context, uuidGame uuid.UUID, uuidUser uuid.UUID) (domain.Game, error) {
					return domain.Game{}, domain.ErrContextCanceled
				}
			},
			gameID:         uuidGame.String(),
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: `{"error":"context canceled"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockGS := &mockGameService{}
			tc.setupmocks(mockGS)

			handler := NewGameHandler(mockGS)

			req := httptest.NewRequest(http.MethodGet, "/games/"+tc.gameID, nil)
			newCtx := context.WithValue(context.Background(), userUUIDKey, uuidUser)
			req = req.WithContext(newCtx)

			req.SetPathValue("uuid", tc.gameID)

			rec := httptest.NewRecorder()

			handler.handleGetGameInfo(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}
