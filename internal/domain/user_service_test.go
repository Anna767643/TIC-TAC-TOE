package domain

import (
	"context"
	"errors"
	"strings"
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
	id := 1
	userUuid := uuid.New()

	testCases := []struct {
		testName      string
		setupMock     func(*mockUserRepository)
		signUpRequest SignUpRequest
		expectedUser  User
		expectedErr   string
		cancelCtx     bool
	}{
		{
			testName: "OK",
			setupMock: func(repo *mockUserRepository) {
				repo.SaveFunc = func(ctx context.Context, user User) error {
					return nil
				}
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{
						ID:    id,
						UUID:  userUuid,
						Login: login,
					}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedUser: User{
				ID:    id,
				UUID:  userUuid,
				Login: "Anna",
			},
			expectedErr: "",
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
			expectedErr: ErrLoginIncorrect.Error(),
		},
		{
			testName: "domain: faled to get user",
			setupMock: func(repo *mockUserRepository) {
				repo.SaveFunc = func(ctx context.Context, user User) error {
					return nil
				}
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, ErrUserNotFound
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedErr: "domain: faled to get user",
		},
		{
			testName: "repository: select user with",
			setupMock: func(repo *mockUserRepository) {
				repo.SaveFunc = func(ctx context.Context, user User) error {
					return nil
				}
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, errors.New("repository: select user with")
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedErr: "repository: select user with",
		},
		{
			testName: "ErrContextCanceled",
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
			expectedErr: ErrContextCanceled.Error(),
			cancelCtx:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

			mock := &mockUserRepository{}

			tc.setupMock(mock)

			u := NewUserService(mock)

			user, err := u.Register(ctx, tc.signUpRequest)

			if tc.expectedErr == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				} else {
					if user.ID != tc.expectedUser.ID {
						t.Errorf("Wrong result: get %v, want %v", user.ID, tc.expectedUser.ID)
					}
					if user.Login != tc.expectedUser.Login {
						t.Errorf("Wrong result: get %s, want %s", user.Login, tc.expectedUser.Login)
					}
					if user.UUID != tc.expectedUser.UUID {
						t.Errorf("Wrong result: get %s, want %s", user.UUID, tc.expectedUser.UUID)
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

func TestLogIn(t *testing.T) {
	userId := uuid.New()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)

	testCases := []struct {
		testName       string
		setupMock      func(*mockUserRepository)
		signUpRequest  SignUpRequest
		expectedResult uuid.UUID
		expectedErr    error
		cancelCtx      bool
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
		{
			testName: "ErrContextCanceled",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, nil
				}
			},
			signUpRequest: SignUpRequest{
				Login:    "Anna",
				Password: "12345",
			},
			expectedResult: uuid.UUID{},
			expectedErr:    ErrContextCanceled,
			cancelCtx:      true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

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

func TestGetUserInfo(t *testing.T) {
	userId := uuid.New()
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("12345"), 10)
	loginUser := "Anna"

	testCases := []struct {
		testName     string
		setupMock    func(*mockUserRepository)
		userUUID     uuid.UUID
		expectedUser User
		expectedErr  error
		cancelCtx    bool
	}{
		{
			testName: "OK",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{
						UUID:         userId,
						Login:        loginUser,
						PasswordHash: string(passwordHash),
					}, nil
				}
			},
			userUUID: userId,
			expectedUser: User{
				UUID:         userId,
				Login:        loginUser,
				PasswordHash: string(passwordHash),
			},
			expectedErr: nil,
		},
		{
			testName: "ErrUserNotFound",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, ErrUserNotFound
				}
			},
			userUUID:    uuid.New(),
			expectedErr: ErrUserNotFound,
		},
		{
			testName: "ErrContextCanceled",
			setupMock: func(repo *mockUserRepository) {
				repo.GetFunc = func(ctx context.Context, login string, uuidUser uuid.UUID) (User, error) {
					return User{}, nil
				}
			},
			userUUID:    userId,
			expectedErr: ErrContextCanceled,
			cancelCtx:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Parallel()

			if !tc.cancelCtx {
				defer cancel()
			} else {
				cancel()
			}

			mock := &mockUserRepository{}

			tc.setupMock(mock)

			u := NewUserService(mock)

			user, err := u.GetUserInfo(ctx, tc.userUUID)

			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Result was incorrect, got: %s, want: %s", err, tc.expectedErr)
			}

			if tc.expectedErr == nil {
				if user.Login != tc.expectedUser.Login {
					t.Errorf("Wrong result: get %s, want %s", user.Login, tc.expectedUser.Login)
				}
				if user.UUID != tc.expectedUser.UUID {
					t.Errorf("Wrong result: get %s, want %s", user.UUID, tc.expectedUser.UUID)
				}
			}

		})
	}
}
