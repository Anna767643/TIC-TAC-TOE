package datasource

import (
	"src/internal/domain"
	"strconv"
	"strings"
)

// toEntity преобразует domain.Game в GameEntity
func toEntityGame(game domain.Game) gameEntity {
	return gameEntity{
		UUID:         game.UUID,
		Board:        toEntityBoard(game.Board),
		WithComputer: game.WithComputer,
		GameState:    game.GameState,

		UUIDPlayerFirst:  game.UUIDPlayerFirst,
		UUIDPlayerSecond: game.UUIDPlayerSecond,

		UpdatedAt: game.UpdatedAt,
	}
}

// toDomain преобразует GameEntity в domain.Game
// если ID в хранилище не парситься в UUID, то возвращается ошибка
func toDomainGame(entity gameEntity) (domain.Game, error) {

	dboard, err := toDomainBoard(entity.Board)
	if err != nil {
		return domain.Game{}, err
	}

	return domain.Game{
		ID:           entity.ID,
		UUID:         entity.UUID,
		Board:        dboard,
		WithComputer: entity.WithComputer,
		GameState:    entity.GameState,

		UUIDPlayerFirst:  entity.UUIDPlayerFirst,
		UUIDPlayerSecond: entity.UUIDPlayerSecond,

		UpdatedAt: entity.UpdatedAt,
	}, nil
}

// toDomainBoard преобразует строку в domain.Board
// если строка не парситься в int, то возвращается ошибка
func toDomainBoard(board string) (domain.Board, error) {
	parts := strings.Split(board, ", ")
	var dboard domain.Board
	n := 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			number, err := strconv.Atoi(parts[n])
			if err != nil {
				return domain.Board{}, err
			}
			dboard[i][j] = number
			n++
		}
	}

	return dboard, nil
}

// toEntityBoard преобразует domain.Board в строку
// если int не парситься в строку, то возвращается ошибка
func toEntityBoard(board domain.Board) string {
	strArray := make([]string, 9)
	i := 0
	for _, r := range board {
		for _, value := range r {
			strArray[i] = strconv.Itoa(value)
			i++
		}
	}
	result := strings.Join(strArray, ", ")
	return result
}

// toUserEntity преобразует domain.User в userEntity
func toUserEntity(user domain.User) userEntity {
	return userEntity{
		UUID:         user.UUID,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		CreateAt:     user.CreateAt,
	}
}

// toDomainUser преобразует userEntity в domain.User
// если ID в хранилище не парситься в UUID, то возвращается ошибка
func toDomainUser(entity userEntity) domain.User {
	return domain.User{
		ID:           entity.ID,
		UUID:         entity.UUID,
		Login:        entity.Login,
		PasswordHash: entity.PasswordHash,
		CreateAt:     entity.CreateAt,
	}
}
