package web

import (
	"context"
	"net/http"
	"src/internal/domain"
)

// contextKey - тип для ключей контекста
type contextKey string

// Константа-ключ
const (
	userUUIDKey contextKey = "userUUID"
)

// UserAuthenticator - защита от неовторизованных запросов
type UserAuthenticator struct {
	userService domain.UserService
}

// NewUserAuthenticator - конструктор для di
func NewUserAuthenticator(userService domain.UserService) *UserAuthenticator {
	return &UserAuthenticator{userService: userService}
}

// Validate - это Middleware. Он принимает следующий хендлер и возвращает новый
func (auth *UserAuthenticator) Validate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		login, password, ok := r.BasicAuth()
		if !ok {
			writeError(w, http.StatusBadRequest, "Need authorization")
			return
		}

		signUpRequest := domain.SignUpRequest{
			Login:    login,
			Password: password,
		}

		uuid, err := auth.userService.LogIn(ctx, signUpRequest)
		if err == domain.ErrLoginOrPasswordIncorrect {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		} else if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Складываем наше значение в контекст
		newCtx := context.WithValue(ctx, userUUIDKey, uuid)

		next(w, r.WithContext(newCtx))
	}
}
