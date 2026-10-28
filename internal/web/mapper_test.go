package web

// import (
// 	"reflect"
// 	"src/internal/domain"
// 	"testing"

// 	uuid "github.com/google/uuid"
// )

// func TestMapperToDomainGame(t *testing.T) {
// 	uuidDomain := uuid.New()
// 	Id := uuidDomain.String()
// 	expectedDomainGame := domain.Game{
// 		Board: domain.Board{
// 			{0, 0, 2},
// 			{1, 0, 0},
// 			{0, 1, 0},
// 		},
// 		UUID: uuidDomain,
// 	}

// 	board := [3][3]int{
// 		{0, 0, 2},
// 		{1, 0, 0},
// 		{0, 1, 0},
// 	}

// 	gameDTO := gameDTO{Id, board}

// 	domainGame, err := toDomain(gameDTO)
// 	if err != nil {
// 		t.Errorf("toDomain() get error, want not error %v", err)
// 		return
// 	}
// 	if !reflect.DeepEqual(domainGame, expectedDomainGame) {
// 		t.Errorf("toDomain() got %v, %v, want %v, %v", domainGame.UUID, domainGame.Board, expectedDomainGame.UUID, expectedDomainGame.Board)
// 	}
// }

// func TestMapperToDTO(t *testing.T) {
// 	uuidDomain := uuid.New()
// 	Id := uuidDomain.String()
// 	DomainGame := domain.Game{
// 		Board: domain.Board{
// 			{0, 0, 2},
// 			{1, 0, 0},
// 			{0, 1, 0},
// 		},
// 		UUID: uuidDomain,
// 	}

// 	board := [3][3]int{
// 		{0, 0, 2},
// 		{1, 0, 0},
// 		{0, 1, 0},
// 	}

// 	expectedGameDTO := gameDTO{Id, board}

// 	gameDTO := toDTO(DomainGame)

// 	if !reflect.DeepEqual(gameDTO, expectedGameDTO) {
// 		t.Errorf("toDomain() got %v, %v, want %v, %v", gameDTO.ID, gameDTO.Board, expectedGameDTO.ID, expectedGameDTO.Board)
// 	}
// }
