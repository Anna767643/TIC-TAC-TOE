package domain

// Validate - проверяет не изменены ли предыдущие ходы и новое значение равно только playerCellValue
// ErrInvalidCell - значения вне диапозона
// ErrInvalidBoard - значение не playerCellValue, больше одного изменения или поле не совпадает
func (s *gameService) validate(incoming, stored Board, playerCellValue int) error {
	var count int

	for i := range stored {
		for j := range stored[i] {
			if incoming[i][j] > O || incoming[i][j] < Empty {
				return ErrInvalidCell
			}
			if stored[i][j] == Empty && incoming[i][j] != Empty {
				if incoming[i][j] != playerCellValue {
					return ErrInvalidBoard
				}
				count++
			}
			if stored[i][j] != Empty && stored[i][j] != incoming[i][j] {
				return ErrInvalidBoard
			}
		}
	}

	if count != 1 {
		return ErrInvalidBoard
	}

	return nil
}
