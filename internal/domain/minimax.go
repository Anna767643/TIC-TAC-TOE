package domain

import (
	"context"
	"strings"
	"time"

	uuid "github.com/google/uuid"
)

// pos - расположение на игровом поле
type pos struct {
	i int
	j int
}

// NextMove - выполняет ход игроком
// Если следующий ход компьютера, то он выполняется с помощью алгоритма минимакс
// Возвращает игру с обновлённым полем.
// Ошибки:
// ErrContextCanceled - игра прерывается
// ErrGameFinished - игра была ранее окончена
// ErrWaitOpponent - второй игрок еще не присоеденился к игре
// ErrGameNotFound - игра не найдена
// ErrDoesntHaveSuchUser - в этой игре нет такого пользователя
// ErrOpponentsMove - ход аппонента
// ErrInvalidBoard - поле нарушает правила или изменены прошлые ходы
// ErrInvalidCell - недопустимое значение клетки
func (s *gameService) NextMove(ctx context.Context, board Board, uuidGame uuid.UUID, uuidUser uuid.UUID) (Game, error) {
	if err := ctx.Err(); err != nil {
		return Game{}, ErrContextCanceled
	}

	// Сперва получим инфу об игре вместе с игроком
	game, err := s.repo.Get(ctx, uuidGame)
	if err != nil {
		return Game{}, err
	}

	// Проверяем что ход нашего игрока
	userNextMove := strings.TrimPrefix(game.GameState, NextMove)
	if game.GameState == WaitUsers {
		return Game{}, ErrWaitOpponent
	} else if game.GameState == Draw || strings.HasPrefix(game.GameState, Win) {
		return Game{}, ErrGameFinished
	}

	// Сверяем UUID игрока
	if uuidUser.String() != userNextMove && game.UUIDPlayerSecond != nil {
		if game.UUIDPlayerFirst != uuidUser && *game.UUIDPlayerSecond != uuidUser {
			return Game{}, ErrDoesntHaveSuchUser
		} else {
			return Game{}, ErrOpponentsMove
		}
	}

	var playerCellValue int
	var playerNextMove uuid.UUID

	// Устанавливаем чем он должен ходить
	if game.UUIDPlayerFirst == uuidUser {
		playerCellValue = X
		if game.UUIDPlayerSecond != nil {
			playerNextMove = *game.UUIDPlayerSecond
		}
	} else {
		playerCellValue = O
		playerNextMove = game.UUIDPlayerFirst
	}

	// Проверяем что ход допустимый
	if err := s.validate(board, game.Board, playerCellValue); err != nil {
		return Game{}, err
	}

	game.Board = board

	// Проверяем не окончена ли игра после этого хода
	if finished, _ := s.isFinished(game.Board); finished {
		game.GameState = Win + uuidUser.String()
		game.UpdatedAt = time.Now()
		if err := s.repo.Save(ctx, game); err != nil {
			return Game{}, err
		}
		return game, nil
	}

	// Если игра с компьютером
	if game.WithComputer {
		return s.nextMoveWithCComputer(ctx, game)
	}

	// Записываем, что следующий ход за оппонентом
	game.GameState = NextMove + playerNextMove.String()
	game.UpdatedAt = time.Now()
	if err := s.repo.Save(ctx, game); err != nil {
		return Game{}, err
	}

	return game, nil
}

func (s *gameService) nextMoveWithCComputer(ctx context.Context, game Game) (Game, error) {
	bestScore := -1000
	var bestPosition pos

	// Перебираем все возможные ходы для O (компьютер)
	for _, posotion := range s.searchAllOptions(game.Board) {
		newBoard := game.Board
		newBoard[posotion.i][posotion.j] = O

		//Уходим в рекурсию для подсчета выгодности данного хода
		score, err := s.minimax(ctx, newBoard, 0, false)
		if err != nil {
			return game, ErrContextCanceled
		}

		// Если счет лучший, то сохраняем
		if score > bestScore {
			bestScore = score
			bestPosition = posotion
		}
	}

	// Устанавливаем позицию в копию и возвращаем
	newGame := game
	newGame.Board[bestPosition.i][bestPosition.j] = O
	newGame.UpdatedAt = time.Now()

	if finished, _ := s.isFinished(newGame.Board); finished {
		newGame.GameState = Win + "computer"
	}

	if err := s.repo.Save(ctx, newGame); err != nil {
		return Game{}, err
	}

	return newGame, nil
}

// minimax возвращает оценку позиции для O: чем выше, тем лучше.
// Флаг isMax показывает чей ход, компьютера или игрока
func (s *gameService) minimax(ctx context.Context, b Board, depth int, isMax bool) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, ErrContextCanceled
	}

	// В зависимости от выиграшевого игрока уcтанавливаем счёт с поправкой на глубину
	if finished, winner := s.isFinished(b); finished {
		switch winner {
		case O:
			return 10 - depth, nil
		case X:
			return depth - 10, nil
		default:
			return 0, nil
		}
	}

	//Уставливаем счет и игрока в зависимости от isMax
	bestScore := -1000
	player := O
	if !isMax {
		bestScore = 1000
		player = X
	}

	// Перебираем все позиции и ищем самую лучшую
	for _, posotion := range s.searchAllOptions(b) {
		newBoard := b
		newBoard[posotion.i][posotion.j] = player

		score, err := s.minimax(ctx, newBoard, depth+1, !isMax)
		if err != nil {
			return 0, ErrContextCanceled
		}

		if isMax && score > bestScore {
			bestScore = score
		} else if !isMax && score < bestScore {
			bestScore = score
		}
	}

	// Возвращаем счет для определенной позиции
	return bestScore, nil
}

// searchAllOptions ищет все доступные пустые клетки на доске.
func (s *gameService) searchAllOptions(board Board) []pos {
	options := make([]pos, 0, 9)
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 0 {
				options = append(options, pos{i, j})
			}
		}
	}
	return options
}
