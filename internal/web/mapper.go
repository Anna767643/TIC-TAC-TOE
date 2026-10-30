package web

import (
	uuid "github.com/google/uuid"
	"src/internal/domain"
	"time"
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
	return gameDTO{
		ID:           game.ID,
		UUID:         game.UUID.String(),
		Board:        boardDTO(game.Board),
		WithComputer: game.WithComputer,
		State:        game.State,
		NextPlayer:   uuidPtrToString(game.NextPlayer),
		Winner:       uuidPtrToString(game.Winner),

		UUIDPlayerFirst:  game.UUIDPlayerFirst.String(),
		UUIDPlayerSecond: uuidPtrToString(game.UUIDPlayerSecond),
		UpdatedAt:        game.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// uuidPtrToString - переводит *uuid.UUID в *string
func uuidPtrToString(val *uuid.UUID) *string {
	if val == nil {
		return nil
	}
	str := val.String()
	return &str
}

// toUserDT - парсит domain.User в userDTO
func toUserDTO(user domain.User) userDTO {
	return userDTO{
		ID:       user.ID,
		UUID:     user.UUID.String(),
		Login:    user.Login,
		CreateAt: user.CreateAt.UTC().Format(time.RFC3339),
	}
}
