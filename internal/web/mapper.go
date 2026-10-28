package web

import (
	uuid "github.com/google/uuid"
	"src/internal/domain"
)

// toDomain преобразует GameDTO в domain.Game
// func toGameDomain(dto gameDTO) (domain.Game, error) {
// 	UUID, err := uuid.Parse(dto.UUID)
// 	if err != nil {
// 		return domain.Game{}, domain.ErrInvalidUUID
// 	}

// 	return domain.Game{
// 		ID: ,
// 		UUID:  UUID,
// 		Board: domain.Board(dto.Board),

// 	}, nil
// }

// toUUID - парсит UUID
// Возвращает UUID или ErrInvalidUUID
func toUUID(str string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(str)
	if err != nil {
		return uuid, domain.ErrInvalidUUID
	}

	return uuid, nil
}

// toDTO - преобразует domain.Game в GameDTO
func toGameDTO(game domain.Game) gameDTO {
	var uuidPlayerSecond string
	if game.UUIDPlayerSecond != nil {
		uuidPlayerSecond = (*game.UUIDPlayerSecond).String()
	}
	return gameDTO{
		ID:               game.ID,
		UUID:             game.UUID.String(),
		Board:            boardDTO(game.Board),
		WithComputer:     game.WithComputer,
		GameState:        game.GameState,
		UUIDPlayerFirst:  game.UUIDPlayerFirst.String(),
		UUIDPlayerSecond: &uuidPlayerSecond,
		UpdatedAt:        game.UpdatedAt.String(),
	}
}

// toUserDT - парсит domain.User в userDTO
func toUserDTO(user domain.User) userDTO {
	return userDTO{
		ID:       user.ID,
		UUID:     user.UUID.String(),
		Login:    user.Login,
		CreateAt: user.CreateAt.String(),
	}
}
