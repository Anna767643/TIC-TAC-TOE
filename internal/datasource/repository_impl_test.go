package datasource

import (
	"context"
	"errors"
	"src/internal/domain"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	return m.selectFromUsersWithUUIDFunc(ctx, uuidUser)
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
		expectedGame domain.Game
		cancelCtx    bool
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
						State:            domain.StateWaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			expectedGame: domain.Game{
				UUID: uuidGame,
				Board: domain.Board{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:     true,
				State:            domain.StateWaitUsers,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
				UpdatedAt:        updatedAt,
			},
			expectedErr: nil,
		},
		{
			testName: "Error to insert",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoGamesFunc = func(ctx context.Context, game gameEntity) error {
					return errors.New("some error")
				}
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{
						UUID:             uuidGame,
						Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
						WithComputer:     true,
						State:            domain.StateWaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  errors.New("some error"),
		},
		{
			testName: "ErrGameNotFound",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{}, domain.ErrGameNotFound
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  domain.ErrGameNotFound,
		},
		{
			testName: "Error to select",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{}, errors.New("some error")
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  errors.New("some error"),
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.insertIntoGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{
						UUID:             uuidGame,
						Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
						WithComputer:     true,
						State:            domain.StateWaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			gr := &gameRepository{
				gs: mock,
			}

			game, err := gr.SaveNewGame(ctx, domain.Game{})

			if tc.expectedErr == domain.ErrContextCanceled || tc.expectedErr == domain.ErrGameNotFound {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("gr.SaveNewGame() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
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

func TestSaveGame(t *testing.T) {
	testCases := []struct {
		testName    string
		setupMock   func(*mockGameStorage)
		cancelCtx   bool
		expectedErr error
	}{
		{
			testName: "OK",
			setupMock: func(m *mockGameStorage) {
				m.updateGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
			},
			expectedErr: nil,
		},
		{
			testName: "Error to update",
			setupMock: func(m *mockGameStorage) {
				m.updateGamesFunc = func(ctx context.Context, game gameEntity) error {
					return errors.New("some error")
				}
			},
			expectedErr: errors.New("some error"),
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.updateGamesFunc = func(ctx context.Context, game gameEntity) error {
					return nil
				}
			},
			expectedErr: domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			gr := &gameRepository{
				gs: mock,
			}
			err := gr.Save(ctx, domain.Game{})

			if tc.expectedErr == domain.ErrContextCanceled {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("gr.Save() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("gr.Save() get error %v, but want %v", err, tc.expectedErr)
				return
			}

		})
	}
}

func TestGetAvailableGames(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	updatedAt := time.Now()

	testCases := []struct {
		testName      string
		setupMock     func(*mockGameStorage)
		expectedGames []domain.Game
		cancelCtx     bool
		expectedErr   error
	}{
		{
			testName: "OK",
			setupMock: func(m *mockGameStorage) {
				m.selectAvaibleGamesFromGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
					return []gameEntity{
						{
							UUID:             uuidGame,
							Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
							WithComputer:     false,
							State:            domain.StateWaitUsers,
							UUIDPlayerFirst:  uuidUser,
							UUIDPlayerSecond: nil,
							UpdatedAt:        updatedAt,
						},
					}, nil
				}
			},
			expectedGames: []domain.Game{
				{
					UUID: uuidGame,
					Board: domain.Board{
						{0, 0, 0},
						{0, 0, 0},
						{0, 0, 0},
					},
					WithComputer:     false,
					State:            domain.StateWaitUsers,
					UUIDPlayerFirst:  uuidUser,
					UUIDPlayerSecond: nil,
					UpdatedAt:        updatedAt,
				},
			},
			expectedErr: nil,
		},
		{
			testName: "Error select",
			setupMock: func(m *mockGameStorage) {
				m.selectAvaibleGamesFromGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
					return []gameEntity{}, errors.New("some error")
				}
			},
			expectedGames: []domain.Game{},
			expectedErr:   errors.New("some error"),
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.selectAvaibleGamesFromGamesFunc = func(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
					return []gameEntity{
						{
							UUID:             uuidGame,
							Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
							WithComputer:     false,
							State:            domain.StateWaitUsers,
							UUIDPlayerFirst:  uuidUser,
							UUIDPlayerSecond: nil,
							UpdatedAt:        updatedAt,
						},
					}, nil
				}
			},
			expectedGames: []domain.Game{
				{
					UUID: uuidGame,
					Board: domain.Board{
						{0, 0, 0},
						{0, 0, 0},
						{0, 0, 0},
					},
					WithComputer:     false,
					State:            domain.StateWaitUsers,
					UUIDPlayerFirst:  uuidUser,
					UUIDPlayerSecond: nil,
					UpdatedAt:        updatedAt,
				},
			},
			expectedErr: domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			gr := &gameRepository{
				gs: mock,
			}

			games, err := gr.GetAvailableGames(ctx, uuidUser)

			if tc.expectedErr == domain.ErrContextCanceled {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("gr.GetAvailableGames() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("gr.GetAvailableGames() get error %v, but want %v", err, tc.expectedErr)
				return
			}

			if tc.expectedErr == nil {
				for i := range games {
					if games[i] != tc.expectedGames[i] {
						t.Errorf("Wrong result: get %v, want %v", games, tc.expectedGames)
					}
				}
			}
		})
	}
}

func TestGetGame(t *testing.T) {
	uuidGame := uuid.New()
	uuidUser := uuid.New()
	updatedAt := time.Now()

	testCases := []struct {
		testName     string
		setupMock    func(*mockGameStorage)
		expectedGame domain.Game
		cancelCtx    bool
		expectedErr  error
	}{
		{
			testName: "OK",
			setupMock: func(m *mockGameStorage) {
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{
						UUID:             uuid,
						Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
						WithComputer:     false,
						State:            domain.StateWaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			expectedGame: domain.Game{
				UUID: uuidGame,
				Board: domain.Board{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:     false,
				State:            domain.StateWaitUsers,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
				UpdatedAt:        updatedAt,
			},
			expectedErr: nil,
		},
		{
			testName: "ErrGameNotFound",
			setupMock: func(m *mockGameStorage) {
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{}, domain.ErrGameNotFound
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  domain.ErrGameNotFound,
		},
		{
			testName: "Error select",
			setupMock: func(m *mockGameStorage) {
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{}, errors.New("some error")
				}
			},
			expectedGame: domain.Game{},
			expectedErr:  errors.New("some error"),
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.selectFromGamesFunc = func(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
					return gameEntity{
						UUID:             uuid,
						Board:            "0, 0, 0, 0, 0, 0, 0, 0, 0",
						WithComputer:     false,
						State:            domain.StateWaitUsers,
						UUIDPlayerFirst:  uuidUser,
						UUIDPlayerSecond: nil,
						UpdatedAt:        updatedAt,
					}, nil
				}
			},
			expectedGame: domain.Game{
				UUID: uuidGame,
				Board: domain.Board{
					{0, 0, 0},
					{0, 0, 0},
					{0, 0, 0},
				},
				WithComputer:     false,
				State:            domain.StateWaitUsers,
				NextPlayer:       &uuidUser,
				UUIDPlayerFirst:  uuidUser,
				UUIDPlayerSecond: nil,
				UpdatedAt:        updatedAt,
			},
			expectedErr: domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			gr := &gameRepository{
				gs: mock,
			}

			game, err := gr.Get(ctx, uuidGame)

			if tc.expectedErr == domain.ErrContextCanceled || tc.expectedErr == domain.ErrGameNotFound {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("gr.Get() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("gr.Gets() get error %v, but want %v", err, tc.expectedErr)
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

func TestSaveUser(t *testing.T) {
	uuidUser := uuid.New()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)
	createAt := time.Now()

	testCases := []struct {
		testName    string
		setupMock   func(*mockGameStorage)
		user        domain.User
		cancelCtx   bool
		expectedErr error
	}{
		{
			testName: "OK",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoUsersFunc = func(ctx context.Context, user userEntity) error {
					return nil
				}
			},
			user: domain.User{
				UUID:         uuidUser,
				Login:        "anna",
				PasswordHash: string(passwordHash),
				CreateAt:     createAt,
			},
			expectedErr: nil,
		},
		{
			testName: "Error insert",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoUsersFunc = func(ctx context.Context, user userEntity) error {
					return errors.New("some error")
				}
			},
			user:        domain.User{},
			expectedErr: errors.New("some error"),
		},
		{
			testName: "ErrLoginIncorrect",
			setupMock: func(m *mockGameStorage) {
				m.insertIntoUsersFunc = func(ctx context.Context, user userEntity) error {
					return domain.ErrLoginIncorrect
				}
			},
			user:        domain.User{},
			expectedErr: domain.ErrLoginIncorrect,
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.insertIntoUsersFunc = func(ctx context.Context, user userEntity) error {
					return nil
				}
			},
			user: domain.User{
				UUID:         uuidUser,
				Login:        "anna",
				PasswordHash: string(passwordHash),
				CreateAt:     createAt,
			},
			expectedErr: domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			ur := &userRepository{
				gs: mock,
			}

			err := ur.Save(ctx, tc.user)

			if tc.expectedErr == domain.ErrContextCanceled || tc.expectedErr == domain.ErrLoginIncorrect {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("ur.Save() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("ur.Save() get error %v, but want %v", err, tc.expectedErr)
				return
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	uuidUser := uuid.New()
	loginAnna := "anna"
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)
	createAt := time.Now()

	testCases := []struct {
		testName     string
		setupMock    func(*mockGameStorage)
		login        string
		uuid         uuid.UUID
		expectedUser domain.User
		cancelCtx    bool
		expectedErr  error
	}{
		{
			testName: "OK, login",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{
						ID:           1,
						UUID:         uuidUser,
						Login:        login,
						PasswordHash: string(passwordHash),
						CreateAt:     createAt,
					}, nil
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{}, nil
				}
			},
			login: loginAnna,
			expectedUser: domain.User{
				ID:           1,
				UUID:         uuidUser,
				Login:        loginAnna,
				PasswordHash: string(passwordHash),
				CreateAt:     createAt,
			},
			expectedErr: nil,
		},
		{
			testName: "OK, uuid",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, nil
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{
						ID:           1,
						UUID:         uuidUser,
						Login:        loginAnna,
						PasswordHash: string(passwordHash),
						CreateAt:     createAt,
					}, nil
				}
			},
			uuid: uuidUser,
			expectedUser: domain.User{
				ID:           1,
				UUID:         uuidUser,
				Login:        loginAnna,
				PasswordHash: string(passwordHash),
				CreateAt:     createAt,
			},
			expectedErr: nil,
		},
		{
			testName: "ErrUserNotFound with login",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, domain.ErrUserNotFound
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{}, nil
				}
			},
			login:        loginAnna,
			expectedUser: domain.User{},
			expectedErr:  domain.ErrUserNotFound,
		},
		{
			testName: "Error select with login",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, errors.New("some error")
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{}, nil
				}
			},
			login:        loginAnna,
			expectedUser: domain.User{},
			expectedErr:  errors.New("some error"),
		},
		{
			testName: "ErrUserNotFound with uuid",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, nil
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{}, domain.ErrUserNotFound
				}
			},
			uuid:         uuidUser,
			expectedUser: domain.User{},
			expectedErr:  domain.ErrUserNotFound,
		},
		{
			testName: "Error select with uuid",
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, nil
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{}, errors.New("some error")
				}
			},
			uuid:         uuidUser,
			expectedUser: domain.User{},
			expectedErr:  errors.New("some error"),
		},
		{
			testName:  "ErrContextCanceled",
			cancelCtx: true,
			setupMock: func(m *mockGameStorage) {
				m.selectFromUsersWithLoginFunc = func(ctx context.Context, login string) (userEntity, error) {
					return userEntity{}, nil
				}
				m.selectFromUsersWithUUIDFunc = func(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
					return userEntity{
						ID:           1,
						UUID:         uuidUser,
						Login:        loginAnna,
						PasswordHash: string(passwordHash),
						CreateAt:     createAt,
					}, nil
				}
			},
			uuid: uuidUser,
			expectedUser: domain.User{
				ID:           1,
				UUID:         uuidUser,
				Login:        loginAnna,
				PasswordHash: string(passwordHash),
				CreateAt:     createAt,
			},
			expectedErr: domain.ErrContextCanceled,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if tc.cancelCtx {
				cancel()
			} else {
				defer cancel()
			}

			mock := &mockGameStorage{}
			tc.setupMock(mock)

			ur := &userRepository{
				gs: mock,
			}

			user, err := ur.Get(ctx, tc.login, tc.uuid)

			if tc.expectedErr == domain.ErrContextCanceled || tc.expectedErr == domain.ErrUserNotFound {
				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("ur.Get() get error %v, but want %v", err, tc.expectedErr)
					return
				}
			} else if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("ur.Get() get error %v, but want %v", err, tc.expectedErr)
				return
			}

			if tc.expectedErr == nil {
				if user != tc.expectedUser {
					t.Errorf("Wrong result: get %v, want %v", user, tc.expectedUser)
				}
			}
		})
	}
}
