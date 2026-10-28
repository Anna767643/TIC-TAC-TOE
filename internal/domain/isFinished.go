package domain

// IsFinished сообщает, окончена ли игра, и кто победил:
// player = X, O или Empty при ничьей / незавершённой игре.
func (s *gameService) isFinished(board Board) (bool, int) {
	var finished bool
	var player int

	if finished = s.isWin(board, X); finished {
		player = X
	} else if finished = s.isWin(board, O); finished {
		player = O
	} else if finished = s.isFilledBoard(board); finished {
		player = Empty
	}

	return finished, player
}

// isWin проверяет выйграл ли определенный игрок
func (s *gameService) isWin(board Board, winner int) bool {
	var (
		count int
		win   bool
	)

	//Проходим по горизонталям
	for idx := range board {
		for _, val := range board[idx] {
			if val == winner {
				count++
			}
			if count == Size {
				win = true
				return win
			}
		}
		count = 0
	}

	//Проходим по вертикалям
	for i := 0; i < Size; i++ {
		for idx := range board {
			if board[idx][i] == winner {
				count++
			}
			if count == Size {
				win = true
				return win
			}
		}
		count = 0
	}

	//Крест-накрест
	if (board[0][0] == winner && board[1][1] == winner && board[2][2] == winner) ||
		(board[0][2] == winner && board[1][1] == winner && board[2][0] == winner) {
		win = true
		return win
	}

	return win
}

// isFilledBoard проверяет заполнено ли поле полностью (в случаях ничьи бывает)
func (s *gameService) isFilledBoard(board Board) bool {
	var (
		count int
		over  bool
	)

	for idx := range board {
		for _, val := range board[idx] {
			if val == 0 {
				count++
			}
		}
	}

	if count == 0 {
		over = true
	}

	return over
}
