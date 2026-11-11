package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"src/internal/domain"
	"testing"

	uuid "github.com/google/uuid"
)

func TestUserAuthenticator_Validate(t *testing.T) {
	testUUID := uuid.New()

	testCases := []struct {
		name         string
		setupMock    func(*mockUserService)
		username     string
		password     string
		expectedCode int
		expectedBody string
		checkContext bool
	}{
		{
			name: "OK - успешная авторизация",
			setupMock: func(m *mockUserService) {
				m.LogInFunc = func(ctx context.Context, req domain.SignUpRequest) (uuid.UUID, error) {
					if req.Login != "testuser" || req.Password != "testpass" {
						t.Errorf("unexpected login/password in request: got %s/%s", req.Login, req.Password)
					}
					return testUUID, nil
				}
			},
			username:     "testuser",
			password:     "testpass",
			expectedCode: http.StatusOK,
			expectedBody: "success",
			checkContext: true,
		},
		{
			name: "Bad Request - отсутствие заголовка авторизации",
			setupMock: func(m *mockUserService) {
				m.LogInFunc = func(ctx context.Context, req domain.SignUpRequest) (uuid.UUID, error) {
					t.Error("LogInFunc should not be called without Basic Auth")
					return uuid.Nil, nil
				}
			},
			username:     "",
			password:     "",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Need authorization"}`,
			checkContext: false,
		},
		{
			name: "Unauthorized - неверный логин или пароль",
			setupMock: func(m *mockUserService) {
				m.LogInFunc = func(ctx context.Context, req domain.SignUpRequest) (uuid.UUID, error) {
					return uuid.Nil, domain.ErrLoginOrPasswordIncorrect
				}
			},
			username:     "wronguser",
			password:     "wrongpass",
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"`,
			checkContext: false,
		},
		{
			name: "Internal Server Error - ошибка сервиса",
			setupMock: func(m *mockUserService) {
				m.LogInFunc = func(ctx context.Context, req domain.SignUpRequest) (uuid.UUID, error) {
					return uuid.Nil, errors.New("database connection failed")
				}
			},
			username:     "testuser",
			password:     "testpass",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"database connection failed"}`,
			checkContext: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockUS := &mockUserService{}
			tc.setupMock(mockUS)

			auth := NewUserAuthenticator(mockUS)

			nextCalled := false

			nextHandler := func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true

				if tc.checkContext {
					val := r.Context().Value(userUUIDKey)
					if val == nil {
						t.Error("userUUIDKey not found in context")
						return
					}

					uuidFromCtx, ok := val.(uuid.UUID)
					if !ok {
						t.Errorf("expected uuid.UUID in context, got %T", val)
						return
					}

					if uuidFromCtx != testUUID {
						t.Errorf("expected UUID %s, got %s", testUUID.String(), uuidFromCtx.String())
					}
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}

			handler := auth.Validate(nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/protected-resource", nil)
			if tc.username != "" {
				req.SetBasicAuth(tc.username, tc.password)
			}

			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf("Expected HTTP status %d, got %d. Body: %s", tc.expectedCode, rec.Code, rec.Body.String())
			}

			if tc.expectedBody != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tc.expectedBody)) {
				t.Errorf("Expected body to contain %q, got %q", tc.expectedBody, rec.Body.String())
			}

			if tc.checkContext && !nextCalled {
				t.Error("Expected next handler to be called, but it wasn't")
			}
			if !tc.checkContext && nextCalled {
				t.Error("Expected next handler NOT to be called, but it was")
			}
		})
	}
}
