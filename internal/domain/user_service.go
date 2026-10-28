package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register - регистрирует нового пользователя
// возвращает факт успешной регистрации
// или ошибку
func (u *userService) Register(ctx context.Context, signUpRequest SignUpRequest) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, ErrContextCanceled
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), 10)
	if err != nil {
		return User{}, err
	}

	user := User{
		UUID:         uuid.New(),
		Login:        signUpRequest.Login,
		PasswordHash: string(passwordHash),
		CreateAt:     time.Now(),
	}

	if err := u.repo.Save(ctx, user); err != nil {
		return User{}, err
	}

	user, err = u.repo.Get(ctx, signUpRequest.Login, uuid.UUID{})
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = ""

	return user, nil
}

// LogIn - авторизирует пользователя
// возвращает UUIDпользователя
// или ошибку
func (u *userService) LogIn(ctx context.Context, signUpRequest SignUpRequest) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", ErrContextCanceled
	}

	user, err := u.repo.Get(ctx, signUpRequest.Login, uuid.UUID{})
	if err != nil {
		return "", ErrLoginOrPasswordIncorrect
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signUpRequest.Password)); err != nil {
		return "", ErrLoginOrPasswordIncorrect
	}

	return user.UUID.String(), nil
}

func (u *userService) GetUserInfo(ctx context.Context, uuidUser uuid.UUID) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, ErrContextCanceled
	}

	user, err := u.repo.Get(ctx, "", uuidUser)
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}
