package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*
Register - регистрирует нового пользователя
возвращает факт успешной регистрации
Ошибки:
ErrContextCanceled
"domain: faled to generate password hash %w"
ErrLoginIncorrect
"repository: insert user with %s: %w"
"repository: select user with %s: %w"
"domain: faled to get user %w"
*/
func (u *userService) Register(ctx context.Context, signUpRequest SignUpRequest) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, ErrContextCanceled
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), 10)
	if err != nil {
		return User{}, fmt.Errorf("domain: faled to generate password hash %w", err)
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
		if err == ErrUserNotFound {
			return User{}, fmt.Errorf("domain: faled to get user %w", err)

		}
		return User{}, err
	}

	user.PasswordHash = ""

	return user, nil
}

/*
LogIn - авторизирует пользователя
возвращает UUIDпользователя
Ошибки:
ErrContextCanceled
ErrLoginOrPasswordIncorrect
"repository: select user with %s: %w"
*/
func (u *userService) LogIn(ctx context.Context, signUpRequest SignUpRequest) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.UUID{}, ErrContextCanceled
	}

	user, err := u.repo.Get(ctx, signUpRequest.Login, uuid.UUID{})
	if err != nil {
		if err == ErrUserNotFound {
			return uuid.UUID{}, ErrLoginOrPasswordIncorrect
		}
		return uuid.UUID{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signUpRequest.Password)); err != nil {
		return uuid.UUID{}, ErrLoginOrPasswordIncorrect
	}

	return user.UUID, nil
}

/*
GetUserInfo - получение информации о пользователе по UUID пользователя
Возвращает информацию о пользователе
Ошибки:
ErrContextCanceled
ErrUserNotFound
"repository: select user with %s: %w"
*/
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
