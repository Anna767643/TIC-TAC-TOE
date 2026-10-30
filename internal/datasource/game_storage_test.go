package datasource

import (
	"context"
	"crypto/rand"
	"strings"

	// "errors"
	"src/internal/domain"
	"testing"
	"time"

	uuid "github.com/google/uuid"
	// "github.com/jackc/pgx/v5"
)

// ------------- GAMES -------------

func TestInsertIntoGames(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	validUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	err = db.insertIntoUsers(ctxSetup, validUser)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	testCases := []struct {
		name          string
		game          gameEntity
		expectedError string
	}{
		{
			name: "OK",
			game: gameEntity{
				UUID:         uuid.New(),
				Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
				WithComputer: true,
				State:        domain.StateNextMove,
				NextPlayer:   &validUser.UUID,

				UUIDPlayerFirst: validUser.UUID,
				UpdatedAt:       time.Now(),
			},
			expectedError: "",
		},
		{
			name: "Error to insert",
			game: gameEntity{
				UUID:         uuid.Nil,
				Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
				WithComputer: true,
				State:        domain.StateNextMove,
				NextPlayer:   &validUser.UUID,

				UUIDPlayerFirst: uuid.Nil,
				UpdatedAt:       time.Now(),
			},
			expectedError: "repository: insert game",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := db.insertIntoGames(ctx, tc.game)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				newGame, err := db.selectFromGames(ctx, tc.game.UUID)
				if err != nil {
					t.Fatalf("Failed to retrieve game after insert: %v", err)
				}

				if newGame.UUID != tc.game.UUID {
					t.Errorf("UUID: got %s, want %s", newGame.UUID, tc.game.UUID)
				}
				if newGame.Board != tc.game.Board {
					t.Errorf("Board: got %s, want %s", newGame.Board, tc.game.Board)
				}
				if newGame.UUIDPlayerFirst != tc.game.UUIDPlayerFirst {
					t.Errorf("UUIDPlayerFirst: got %s, want %s", newGame.UUIDPlayerFirst, tc.game.UUIDPlayerFirst)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}

func TestUpdateGames(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	validUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	validGame := gameEntity{
		UUID:         uuid.New(),
		Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
		WithComputer: true,
		State:        domain.StateNextMove,
		NextPlayer:   &validUser.UUID,

		UUIDPlayerFirst: validUser.UUID,
		UpdatedAt:       time.Now(),
	}

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	err = db.insertIntoUsers(ctxSetup, validUser)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	err = db.insertIntoGames(ctxSetup, validGame)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	testCases := []struct {
		name          string
		game          gameEntity
		expectedError string
	}{
		{
			name: "OK",
			game: gameEntity{
				UUID:       validGame.UUID,
				Board:      "0, 1, 0, 0, 0, 0, 0, 0, 0",
				State:      domain.StateNextMove,
				NextPlayer: &validUser.UUID,

				UpdatedAt: time.Now(),
			},
			expectedError: "",
		},
		{
			name: "Error to update",
			game: gameEntity{
				UUID:  validGame.UUID,
				State: "teeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeest",
			},
			expectedError: "repository: update game",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := db.updateGames(ctx, tc.game)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				newGame, err := db.selectFromGames(ctx, tc.game.UUID)
				if err != nil {
					t.Fatalf("Failed to retrieve game after insert: %v", err)
				}

				if newGame.Board != tc.game.Board {
					t.Errorf("Board: got %s, want %s", newGame.Board, tc.game.Board)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}

func TestSelectAvaibleGamesFromGames(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	validUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	validGame := gameEntity{
		UUID:         uuid.New(),
		Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
		WithComputer: false,
		State:        domain.StateNextMove,
		NextPlayer:   &validUser.UUID,

		UUIDPlayerFirst: validUser.UUID,
		UpdatedAt:       time.Now(),
	}

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	err = db.insertIntoUsers(ctxSetup, validUser)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	err = db.insertIntoGames(ctxSetup, validGame)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	testCases := []struct {
		name          string
		uuidUser      uuid.UUID
		expectedError string
		canceledCtx   bool
	}{
		{
			name:          "OK",
			uuidUser:      validUser.UUID,
			expectedError: "",
		},
		{
			name:          "Game not found",
			uuidUser:      uuid.New(),
			expectedError: "repository: select game for user",
			canceledCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if tc.canceledCtx {
				cancel()
			} else {
				defer cancel()
			}

			_, err := db.selectAvaibleGamesFromGames(ctx, tc.uuidUser)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}

func TestSelectFromGames(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	validUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	validGame := gameEntity{
		UUID:         uuid.New(),
		Board:        "0, 0, 0, 0, 0, 0, 0, 0, 0",
		WithComputer: false,
		State:        domain.StateNextMove,
		NextPlayer:   &validUser.UUID,

		UUIDPlayerFirst: validUser.UUID,
		UpdatedAt:       time.Now(),
	}

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	err = db.insertIntoUsers(ctxSetup, validUser)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	err = db.insertIntoGames(ctxSetup, validGame)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	testCases := []struct {
		name          string
		uuidGame      uuid.UUID
		expectedGame  gameEntity
		expectedError string
		canceledCtx   bool
	}{
		{
			name:          "OK",
			uuidGame:      validGame.UUID,
			expectedGame:  validGame,
			expectedError: "",
		},
		{
			name:          "Game not found",
			uuidGame:      uuid.New(),
			expectedError: domain.ErrGameNotFound.Error(),
		},
		{
			name:          "Game not found",
			uuidGame:      uuid.New(),
			expectedError: "epository: select game",
			canceledCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if tc.canceledCtx {
				cancel()
			} else {
				defer cancel()
			}

			game, err := db.selectFromGames(ctx, tc.uuidGame)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				if game.Board != tc.expectedGame.Board {
					t.Errorf("Board: got %s, want %s", game.Board, tc.expectedGame.Board)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}

// ------------- USERS -------------

func TestInsertIntoUsers(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	insertUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	err = db.insertIntoUsers(ctxSetup, insertUser)
	if err != nil {
		t.Errorf("Failed to insert to the database before testing: %v", err)
	}

	testCases := []struct {
		name          string
		user          userEntity
		expectedError string
		canceledCtx   bool
	}{
		{
			name: "OK",
			user: userEntity{
				UUID:         uuid.New(),
				Login:        rand.Text(),
				PasswordHash: rand.Text(),
				CreateAt:     time.Now(),
			},
			expectedError: "",
		},
		{
			name: "a user with this login already exists",
			user: userEntity{
				UUID:         uuid.New(),
				Login:        insertUser.Login,
				PasswordHash: rand.Text(),
				CreateAt:     time.Now(),
			},
			expectedError: domain.ErrLoginIncorrect.Error(),
		},
		{
			name: "error",
			user: userEntity{
				UUID:         uuid.New(),
				Login:        rand.Text(),
				PasswordHash: rand.Text(),
				CreateAt:     time.Now(),
			},
			expectedError: "repository: insert user with",
			canceledCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if tc.canceledCtx {
				cancel()
			} else {
				defer cancel()
			}

			err := db.insertIntoUsers(ctx, tc.user)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				newUser, err := db.selectFromUsersWithLogin(ctx, tc.user.Login)
				if err != nil {
					t.Fatalf("selectFromUsersWithLogin got: %v, but want nil", err)
				}

				if newUser.UUID != tc.user.UUID {
					t.Errorf("UUID wrong: get %s, want %s", newUser.UUID, tc.user.UUID)
				}
				if newUser.Login != tc.user.Login {
					t.Errorf("Login wrong: get %s, want %s", newUser.Login, tc.user.Login)
				}
				if newUser.PasswordHash != tc.user.PasswordHash {
					t.Errorf("PasswordHash wrong: get %s, want %s", newUser.PasswordHash, tc.user.PasswordHash)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}

func TestSelectFromUsersWithLogin(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	insertUser := userEntity{
		UUID:         uuid.New(),
		Login:        rand.Text(),
		PasswordHash: rand.Text(),
		CreateAt:     time.Now(),
	}

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSetup()

	err = db.insertIntoUsers(ctxSetup, insertUser)
	if err != nil {
		t.Fatalf("Error to insert: %v", err)
	}

	testCases := []struct {
		name          string
		loginUser     string
		expectedUser  userEntity
		expectedError string
		canceledCtx   bool
	}{
		{
			name:          "OK",
			loginUser:     insertUser.Login,
			expectedUser:  insertUser,
			expectedError: "",
		},
		{
			name:          "Wrong uuid game",
			loginUser:     rand.Text(),
			expectedError: domain.ErrUserNotFound.Error(),
		},
		{
			name:          "context canceled",
			loginUser:     rand.Text(),
			expectedError: "repository: select user with",
			canceledCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if tc.canceledCtx {
				cancel()
			} else {
				defer cancel()
			}

			user, err := db.selectFromUsersWithLogin(ctx, tc.loginUser)

			if tc.expectedError == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}

				if user.UUID != tc.expectedUser.UUID {
					t.Errorf("UUID wrong: get %s, want %s", user.UUID, tc.expectedUser.UUID)
				}
				if user.Login != tc.expectedUser.Login {
					t.Errorf("Login wrong: get %s, want %s", user.Login, tc.expectedUser.Login)
				}
				if user.PasswordHash != tc.expectedUser.PasswordHash {
					t.Errorf("PasswordHash wrong: get %s, want %s", user.PasswordHash, tc.expectedUser.PasswordHash)
				}

			} else {
				if err == nil {
					t.Fatalf("Expected error containing %q, but got nil", tc.expectedError)
				}

				if !strings.Contains(strings.ToLower(err.Error()), tc.expectedError) {
					t.Errorf("Expected error to contain %q, but got: %v", tc.expectedError, err)
				}
			}

		})
	}
}
