package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"src/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
)

// mockUserService - заглушка UserService
type mockUserService struct {
	RegisterFunc    func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error)
	LogInFunc       func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error)
	GetUserInfoFunc func(ctx context.Context, uuidUser uuid.UUID) (domain.User, error)
}

func (m *mockUserService) Register(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
	return m.RegisterFunc(ctx, signUpRequest)
}

func (m *mockUserService) LogIn(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
	return m.LogInFunc(ctx, signUpRequest)
}

func (m *mockUserService) GetUserInfo(ctx context.Context, uuidUser uuid.UUID) (domain.User, error) {
	return m.GetUserInfoFunc(ctx, uuidUser)
}

func TestHandleUsersCreate(t *testing.T) {

	type usersCreateRequest struct {
		Login    string
		Password string
	}

	testCases := []struct {
		name           string
		requestBody    usersCreateRequest
		setupmocks     func(*mockUserService)
		expectedCode   int
		expectedInBody string
	}{
		{
			name: "OK",
			requestBody: usersCreateRequest{
				Login:    "anna",
				Password: "12345",
			},
			setupmocks: func(us *mockUserService) {
				us.RegisterFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
					return domain.User{
						ID:       2,
						UUID:     uuid.New(),
						Login:    "anna",
						CreateAt: time.Now(),
					}, nil
				}
			},
			expectedCode:   http.StatusCreated,
			expectedInBody: "",
		},
		{
			name:        "Wrong request body",
			requestBody: usersCreateRequest{},
			setupmocks: func(us *mockUserService) {
				us.RegisterFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
					return domain.User{}, nil
				}
			},
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `"wrong body, login or password not be empty"`,
		},
		{
			name:        "Empty json",
			requestBody: usersCreateRequest{},
			setupmocks: func(us *mockUserService) {
				us.RegisterFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
					return domain.User{}, nil
				}
			},
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `"wrong body, need login and password in json, error: EOF"`,
		},
		{
			name: "Already exists",
			requestBody: usersCreateRequest{
				Login:    "anna",
				Password: "12345",
			},
			setupmocks: func(us *mockUserService) {
				us.RegisterFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
					return domain.User{}, domain.ErrLoginIncorrect
				}
			},
			expectedCode:   http.StatusConflict,
			expectedInBody: `"a user with this login already exists"`,
		},
		{
			name: "Context canceled",
			requestBody: usersCreateRequest{
				Login:    "anna",
				Password: "12345",
			},
			setupmocks: func(us *mockUserService) {
				us.RegisterFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (domain.User, error) {
					return domain.User{}, domain.ErrContextCanceled
				}
			},
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: `"context canceled"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Настраиваем моки
			mockUS := &mockUserService{}
			tc.setupmocks(mockUS)

			// Хэндлер с нашими моками
			handler := NewUserHandler(mockUS)

			// Создаем запрос
			var bodyBytes []byte
			if tc.name != "Empty json" {
				bodyBytes, _ = json.Marshal(tc.requestBody)
			}
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Перехватываем ответ
			rec := httptest.NewRecorder()

			handler.handleUsersCreate(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}

func TestHandleSessionCreate(t *testing.T) {

	testCases := []struct {
		name           string
		login          string
		password       string
		setupmocks     func(*mockUserService)
		expectedCode   int
		expectedInBody string
	}{
		{
			name:     "OK",
			login:    "anna",
			password: "12345",
			setupmocks: func(us *mockUserService) {
				us.LogInFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
					return "303c12d2-618f-45ec-94b5-260e68e156d0", nil
				}
			},
			expectedCode:   http.StatusOK,
			expectedInBody: `{"UUID":"303c12d2-618f-45ec-94b5-260e68e156d0"}`,
		},
		{
			name:     "Wrong password",
			login:    "anna",
			password: "1245",
			setupmocks: func(us *mockUserService) {
				us.LogInFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
					return "", domain.ErrLoginOrPasswordIncorrect
				}
			},
			expectedCode:   http.StatusUnauthorized,
			expectedInBody: `"the login or password is incorrect"`,
		},
		{
			name:     "Wrong login",
			login:    "ana",
			password: "12345",
			setupmocks: func(us *mockUserService) {
				us.LogInFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
					return "", domain.ErrLoginOrPasswordIncorrect
				}
			},
			expectedCode:   http.StatusUnauthorized,
			expectedInBody: `"the login or password is incorrect"`,
		},
		{
			name:     "Empty auth",
			login:    "anna",
			password: "",
			setupmocks: func(us *mockUserService) {
				us.LogInFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
					return "", domain.ErrLoginOrPasswordIncorrect
				}
			},
			expectedCode:   http.StatusBadRequest,
			expectedInBody: `"need authorization"`,
		},
		{
			name:     "Context canceled",
			login:    "anna",
			password: "12345",
			setupmocks: func(us *mockUserService) {
				us.LogInFunc = func(ctx context.Context, signUpRequest domain.SignUpRequest) (string, error) {
					return "", domain.ErrContextCanceled
				}
			},
			expectedCode:   http.StatusInternalServerError,
			expectedInBody: `"context canceled"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Настраиваем моки
			mockUS := &mockUserService{}
			tc.setupmocks(mockUS)

			// Хэндлер с нашими моками
			handler := NewUserHandler(mockUS)

			// Создаем запрос
			req := httptest.NewRequest(http.MethodPost, "/sessions", nil)
			if tc.login != "" && tc.password != "" {
				req.SetBasicAuth(tc.login, tc.password)
			}

			// Перехватываем ответ
			rec := httptest.NewRecorder()

			handler.handleSessionCreate(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Wait HTTP status %d, get %d. Body request: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedInBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedInBody)) {
				t.Errorf("Wait body response %q, get: %s", tc.expectedInBody, rec.Body.String())
			}
		})
	}
}
