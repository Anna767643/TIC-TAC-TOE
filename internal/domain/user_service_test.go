package domain

import (
	"context"
	"errors"
	"testing"

	uuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepository - заглушка для репозитория,
// что бы проверить работу функций userService
type mockUserRepository struct {
	SaveFunc func(ctx context.Context, user User) error
	GetFunc  func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error)
}

func (m *mockUserRepository) Save(ctx context.Context, user User) error {
	return m.SaveFunc(ctx, user)
}

func (m *mockUserRepository) Get(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
	return m.GetFunc(ctx, login, uuidUser)
}

func TestRegister(t *testing.T) {
	testCases := []struct {
		testName      string
		setupMock     func(*mockUserRepository)
		signUpRequest SignUpRequest
		expectedErr   error
		cancelCtx     bool
	}{
		{
			testName: "OK",
			setupMock: func(repo *mockUserRepository) {
				repo.SaveFunc = func(ctx context.Context, user User) error {
					return nil
				}
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedErr: nil,
		},
		{
			testName: "User with this login already exists",
			setupMock: func(repo *mockUserRepository) {
				repo.SaveFunc = func(ctx context.Context, user User) error {
					return ErrLoginIncorrect
				}
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedErr: ErrLoginIncorrect,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mock := &mockUserRepository{}

			tc.setupMock(mock)

			u := NewUserService(mock)

			_, err := u.Register(ctx, tc.signUpRequest)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Result was incorrect, got: %s, want: %s", err, tc.expectedErr)
			}
		})
	}
}

func TestLogIn(t *testing.T) {
	userId := uuid.New()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)

	testCases := []struct {
		testName       string
		setupMock      func(*mockUserRepository)
		signUpRequest  SignUpRequest
		expectedResult uuid.UUID
		expectedErr    error
	}{
		{
			testName: "OK",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{
						UUID:         userId,
						Login:        login,
						PasswordHash: string(passwordHash),
					}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedResult: userId,
			expectedErr:    nil,
		},
		{
			testName: "Password is incorrect",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{
						UUID:         userId,
						Login:        login,
						PasswordHash: string(passwordHash),
					}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "123",
			},
			expectedResult: uuid.UUID{},
			expectedErr:    ErrLoginOrPasswordIncorrect,
		},
		{
			testName: "Login is incorrect",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, ErrUserNotFound
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedResult: uuid.UUID{},
			expectedErr:    ErrLoginOrPasswordIncorrect,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			mock := &mockUserRepository{}

			tc.setupMock(mock)

			u := NewUserService(mock)

			result, err := u.LogIn(ctx, tc.signUpRequest)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Result was incorrect, got: %s, want: %s", err, tc.expectedErr)
			}

			if result != tc.expectedResult {
				t.Errorf("Result was incorrect, got: %s, want: %s", result, tc.expectedResult)
			}
		})
	}
}
