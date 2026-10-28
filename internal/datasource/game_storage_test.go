package datasource

// import (
// 	"context"
// 	"crypto/rand"
// 	"errors"
// 	"src/internal/domain"
// 	"testing"
// 	"time"

// 	uuid "github.com/google/uuid"
// 	"github.com/jackc/pgx/v5"
// )

// // ---- GAMES ----
// // func TestPutGame(t *testing.T) {
// // 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// // 	if err != nil {
// // 		t.Fatalf("Failed to connect to the database: %v", err)
// // 	}
// // 	defer db.pool.Close()

// // 	testCase := []struct {
// // 		name          string
// // 		game          gameEntity
// // 		expectedError error
// // 	}{
// // 		{
// // 			name: "OK",
// // 			game: gameEntity{
// // 				UUID:           uuid.New().String(),
// // 				Board:          "0, 0, 0, 0, 0, 0, 0, 0, 0",
// // 				WithComputer:   true,
// // 				GameState:      domain.WaitUsers,
// // 				IDPlayerFirst:  1,
// // 				IDPlayerSecond: nil,
// // 				UpdatedAt:      time.Now(),
// // 			},
// // 			expectedError: nil,
// // 		},
// // 		{
// // 			name: "Wrong id player",
// // 			game: gameEntity{
// // 				UUID:           uuid.New().String(),
// // 				Board:          "0, 0, 0, 0, 0, 0, 0, 0, 0",
// // 				WithComputer:   true,
// // 				GameState:      domain.WaitUsers,
// // 				IDPlayerFirst:  1000000,
// // 				IDPlayerSecond: nil,
// // 				UpdatedAt:      time.Now(),
// // 			},
// // 			expectedError: errors.New("some error"),
// // 		},
// // 	}

// // 	for _, tc := range testCase {
// // 		t.Run(tc.name, func(t *testing.T) {
// // 			ctx, cancel := context.WithCancel(context.Background())
// // 			defer cancel()

// // 			err = db.putGame(ctx, tc.game)

// // 			if (err != nil) != (tc.expectedError != nil) {
// // 				t.Errorf("db.putGame() get error %v, but want %v", err, tc.expectedError)
// // 				return
// // 			}

// // 			if tc.expectedError == nil {
// // 				newGame, err := db.selectFromGames(ctx, tc.game.UUID)
// // 				if err != nil {
// // 					t.Fatalf("Don't get game after insert/update : %v", err)
// // 				}

// // 				if newGame.UUID != tc.game.UUID {
// // 					t.Errorf("UUID wrong: get %s, want %s", newGame.UUID, tc.game.UUID)
// // 				}
// // 				if newGame.Board != tc.game.Board {
// // 					t.Errorf("Board wrong: get %s, want %s", newGame.Board, tc.game.Board)
// // 				}
// // 			}

// // 		})
// // 	}
// // }

// func TestInsertIntoGames(t *testing.T) {
// 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.pool.Close()

// 	insertUser := userEntity{
// 		UUID:         uuid.New().String(),
// 		Login:        rand.Text(),
// 		PasswordHash: rand.Text(),
// 		CreateAt:     time.Now(),
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	err = db.insertIntoUsers(ctx, insertUser)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	testCases := []struct {
// 		name          string
// 		game          gameEntity
// 		userUUID      string
// 		expectedError error
// 	}{
// 		{
// 			name: "OK",
// 			game: gameEntity{
// 				UUID:         uuid.New().String(),
// 				Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
// 				WithComputer: true,
// 				GameState:    domain.WaitUsers,
// 				UpdatedAt:    time.Now(),
// 			},
// 			userUUID:      insertUser.UUID,
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Wrong uuid player",
// 			game: gameEntity{
// 				UUID:         uuid.New().String(),
// 				Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
// 				WithComputer: true,
// 				GameState:    domain.WaitUsers,
// 				UpdatedAt:    time.Now(),
// 			},
// 			userUUID:      uuid.New().String(),
// 			expectedError: errors.New("Wrong uuid player"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			defer cancel()

// 			err := db.insertIntoGames(ctx, tc.game, tc.userUUID)

// 			if (err != nil) != (tc.expectedError != nil) {
// 				t.Errorf("db.insertIntoGames() get error %v, but want %v", err, tc.expectedError)
// 				return
// 			}

// 			if tc.expectedError == nil {
// 				newGame, err := db.selectFromGames(ctx, tc.game.UUID)
// 				if err != nil {
// 					t.Fatalf("Don't get game after insert/update : %v", err)
// 				}

// 				if newGame.UUID != tc.game.UUID {
// 					t.Errorf("UUID wrong: get %s, want %s", newGame.UUID, tc.game.UUID)
// 				}
// 				if newGame.Board != tc.game.Board {
// 					t.Errorf("Board wrong: get %s, want %s", newGame.Board, tc.game.Board)
// 				}
// 			}

// 		})
// 	}
// }

// func TestUpdateGames(t *testing.T) {
// 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.pool.Close()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	insertUser := userEntity{
// 		UUID:         uuid.New().String(),
// 		Login:        rand.Text(),
// 		PasswordHash: rand.Text(),
// 		CreateAt:     time.Now(),
// 	}

// 	err = db.insertIntoUsers(ctx, insertUser)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	insertGame := gameEntity{
// 		UUID:         uuid.New().String(),
// 		Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
// 		WithComputer: true,
// 		GameState:    domain.WaitUsers,
// 		UpdatedAt:    time.Now(),
// 	}

// 	err = db.insertIntoGames(ctx, insertGame, insertUser.UUID)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	id_user_second := 10000

// 	testCases := []struct {
// 		name          string
// 		game          gameEntity
// 		expectedError error
// 	}{
// 		{
// 			name: "OK",
// 			game: gameEntity{
// 				UUID:           insertGame.UUID,
// 				Board:          "0, 0, 1, 0, 0, 0, 0, 0, 0",
// 				GameState:      domain.WaitUsers,
// 				IDPlayerSecond: nil,
// 				UpdatedAt:      time.Now(),
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "Wrong id player",
// 			game: gameEntity{
// 				UUID:           insertGame.UUID,
// 				Board:          "0, 0, 1, 0, 0, 0, 0, 0, 0",
// 				GameState:      domain.WaitUsers,
// 				IDPlayerSecond: &id_user_second,
// 				UpdatedAt:      time.Now(),
// 			},
// 			expectedError: errors.New("some error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := db.updateGames(ctx, tc.game)

// 			if (err != nil) != (tc.expectedError != nil) {
// 				t.Errorf("db.updateGames() get error %v, but want %v", err, tc.expectedError)
// 				return
// 			}

// 			if tc.expectedError == nil {
// 				newGame, err := db.selectFromGames(ctx, tc.game.UUID)
// 				if err != nil {
// 					t.Fatalf("Don't get game after insert/update : %v", err)
// 				}

// 				if newGame.UUID != tc.game.UUID {
// 					t.Errorf("UUID wrong: get %s, want %s", newGame.UUID, tc.game.UUID)
// 				}
// 				if newGame.Board != tc.game.Board {
// 					t.Errorf("Board wrong: get %s, want %s", newGame.Board, tc.game.Board)
// 				}
// 			}

// 		})
// 	}
// }

// func TestSelectFromGames(t *testing.T) {
// 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.pool.Close()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	insertUser := userEntity{
// 		UUID:         uuid.New().String(),
// 		Login:        rand.Text(),
// 		PasswordHash: rand.Text(),
// 		CreateAt:     time.Now(),
// 	}

// 	err = db.insertIntoUsers(ctx, insertUser)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	insertGame := gameEntity{
// 		UUID:           uuid.New().String(),
// 		Board:          "0, 0, 0, 0, 0, 0, 0, 0, 0",
// 		WithComputer:   true,
// 		GameState:      domain.WaitUsers,
// 		IDPlayerFirst:  1,
// 		IDPlayerSecond: nil,
// 		UpdatedAt:      time.Now(),
// 	}

// 	err = db.insertIntoGames(ctx, insertGame, insertUser.UUID)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	testCases := []struct {
// 		name          string
// 		uuidGame      string
// 		expectedGame  gameEntity
// 		expectedError error
// 	}{
// 		{
// 			name:          "OK",
// 			uuidGame:      insertGame.UUID,
// 			expectedGame:  insertGame,
// 			expectedError: nil,
// 		},
// 		{
// 			name:          "Wrong uuid game",
// 			uuidGame:      uuid.New().String(),
// 			expectedError: pgx.ErrNoRows,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			game, err := db.selectFromGames(ctx, tc.uuidGame)

// 			if !errors.Is(err, tc.expectedError) {
// 				t.Errorf("db.selectFromGames() get error %v, but want %v", err, tc.expectedError)
// 				return
// 			}

// 			if tc.expectedError == nil {

// 				if game.UUID != tc.expectedGame.UUID {
// 					t.Errorf("UUID wrong: get %s, want %s", game.UUID, tc.expectedGame.UUID)
// 				}
// 				if game.Board != tc.expectedGame.Board {
// 					t.Errorf("Board wrong: get %s, want %s", game.Board, tc.expectedGame.Board)
// 				}
// 			}

// 		})
// 	}
// }

// // ---- USERS ----
// func TestInsertIntoUsers(t *testing.T) {
// 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.pool.Close()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	insertUser := userEntity{
// 		UUID:         uuid.New().String(),
// 		Login:        rand.Text(),
// 		PasswordHash: rand.Text(),
// 		CreateAt:     time.Now(),
// 	}

// 	err = db.insertIntoUsers(ctx, insertUser)
// 	if err != nil {
// 		t.Errorf("Failed to insert to the database before testing: %v", err)
// 	}

// 	testCases := []struct {
// 		name          string
// 		user          userEntity
// 		expectedError error
// 	}{
// 		{
// 			name: "OK",
// 			user: userEntity{
// 				UUID:         uuid.New().String(),
// 				Login:        rand.Text(),
// 				PasswordHash: rand.Text(),
// 				CreateAt:     time.Now(),
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name: "a user with this login already exists",
// 			user: userEntity{
// 				UUID:         uuid.New().String(),
// 				Login:        insertUser.Login,
// 				PasswordHash: rand.Text(),
// 				CreateAt:     time.Now(),
// 			},
// 			expectedError: errors.New("some error"),
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := db.insertIntoUsers(ctx, tc.user)

// 			if (err != nil) != (tc.expectedError != nil) {
// 				t.Errorf("db.insertIntoUsers() get error %v, but want %v", err, tc.expectedError)
// 				return
// 			}

// 			if tc.expectedError == nil {
// 				newUser, err := db.selectFromUsers(ctx, tc.user.Login)
// 				if err != nil {
// 					t.Fatalf("Don't get user after insert : %v", err)
// 				}

// 				if newUser.UUID != tc.user.UUID {
// 					t.Errorf("UUID wrong: get %s, want %s", newUser.UUID, tc.user.UUID)
// 				}
// 				if newUser.Login != tc.user.Login {
// 					t.Errorf("Login wrong: get %s, want %s", newUser.Login, tc.user.Login)
// 				}
// 				if newUser.PasswordHash != tc.user.PasswordHash {
// 					t.Errorf("PasswordHash wrong: get %s, want %s", newUser.PasswordHash, tc.user.PasswordHash)
// 				}
// 			}

// 		})
// 	}
// }

// func TestSelectFromUsers(t *testing.T) {
// 	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
// 	if err != nil {
// 		t.Fatalf("Failed to connect to the database: %v", err)
// 	}
// 	defer db.pool.Close()

// 	insertUser := userEntity{
// 		UUID:         uuid.New().String(),
// 		Login:        rand.Text(),
// 		PasswordHash: rand.Text(),
// 		CreateAt:     time.Now(),
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	err = db.insertIntoUsers(ctx, insertUser)
// 	if err != nil {
// 		t.Fatalf("Error to insert: %v", err)
// 	}

// 	testCases := []struct {
// 		name          string
// 		loginUser     string
// 		expectedUser  userEntity
// 		expectedError error
// 	}{
// 		{
// 			name:          "OK",
// 			loginUser:     insertUser.Login,
// 			expectedUser:  insertUser,
// 			expectedError: nil,
// 		},
// 		{
// 			name:          "Wrong uuid game",
// 			loginUser:     rand.Text(),
// 			expectedError: pgx.ErrNoRows,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			user, err := db.selectFromUsers(ctx, tc.loginUser)

// 			if !errors.Is(err, tc.expectedError) {
// 				t.Errorf("db.selectFromUsers() get error %v, but want %v", err, tc.expectedError)
// 				return
// 			}

// 			if tc.expectedError == nil {

// 				if user.UUID != tc.expectedUser.UUID {
// 					t.Errorf("UUID wrong: get %s, want %s", user.UUID, tc.expectedUser.UUID)
// 				}
// 				if user.Login != tc.expectedUser.Login {
// 					t.Errorf("Login wrong: get %s, want %s", user.Login, tc.expectedUser.Login)
// 				}
// 				if user.PasswordHash != tc.expectedUser.PasswordHash {
// 					t.Errorf("PasswordHash wrong: get %s, want %s", user.PasswordHash, tc.expectedUser.PasswordHash)
// 				}
// 			}

// 		})
// 	}
// }
