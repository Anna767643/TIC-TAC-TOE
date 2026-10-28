package domain

import (
	"errors"
)

var (
	ErrGameNotFound       = errors.New("game not found")                                                     // репозиторий не нашёл игру
	ErrAddUser            = errors.New("this game already has two users or the uuid matches the first user") // нельзя добавить этого пользователя к игре, или у игры уже два пользователя
	ErrInvalidBoard       = errors.New("invalid board")                                                      // поле нарушает правила или изменены прошлые ходы
	ErrGameFinished       = errors.New("game is already finished")                                           // игра окончена
	ErrInvalidCell        = errors.New("cell value out of range")                                            // недопустимое значение клетки
	ErrInvalidUUID        = errors.New("invalid UUID")                                                       // невалидный UUID
	ErrContextCanceled    = errors.New("context canceled")                                                   // контекст закрыт
	ErrOpponentsMove      = errors.New("wait for your opponent's move")                                      // ход оппонента
	ErrDoesntHaveSuchUser = errors.New("this game doesn't have such a user")                                 // в этой игре нет такого пользователя
	ErrWaitOpponent       = errors.New("your opponent has not joined the game yet")                          // второй игрок еще не присоеденился к игре

	ErrLoginOrPasswordIncorrect = errors.New("the login or password is incorrect")    // неверный логин или пароль
	ErrLoginIncorrect           = errors.New("a user with this login already exists") // пользователь с таким логином уже существует
	ErrUserNotFound             = errors.New("user not found")                        // пользователь не найден (но пользователь данную ошибку никогда не увидит)

)
