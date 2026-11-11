package datasource

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestCreateTableUsers(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	testCases := []struct {
		name          string
		expectedError string
		canceledCtx   bool
	}{
		{
			name:          "OK",
			expectedError: "",
		},
		{
			name:          "Context Canceled",
			expectedError: "context canceled",
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

			err := CreateTableUsers(ctx, db.pool)

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

func TestCreateTableGames(t *testing.T) {
	db, err := NewDBPool("postgres://postgres:12345@localhost:5432/tic_tac_toe")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.pool.Close()

	ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
	err = CreateTableUsers(ctxSetup, db.pool)
	cancelSetup()
	if err != nil {
		t.Fatalf("Failed to create prerequisite users table: %v", err)
	}

	testCases := []struct {
		name          string
		expectedError string
		canceledCtx   bool
	}{
		{
			name:          "OK",
			expectedError: "",
		},
		{
			name:          "Context Canceled",
			expectedError: "context canceled",
			canceledCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if tc.canceledCtx {
				cancel() // Немедленно отменяем контекст
			} else {
				defer cancel()
			}

			err := CreateTableGames(ctx, db.pool)

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
