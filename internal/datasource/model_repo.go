package datasource

import (
	uuid "github.com/google/uuid"
	"time"
)

// gameEntity представляет запись в таблице games
type gameEntity struct {
	ID           int       `db:id`
	UUID         uuid.UUID `db:uuid`
	Board        string    `db:board`
	WithComputer bool      `db:with_computer`
	GameState    string    `db:game_state`

	UUIDPlayerFirst  uuid.UUID  `db:id_player_first`  // X
	UUIDPlayerSecond *uuid.UUID `db:id_player_second` // O

	UpdatedAt time.Time `db:update_at`
}

// userEntity представляет запись в таблице players
type userEntity struct {
	ID           int       `db:id`
	UUID         uuid.UUID `db:uuid`
	Login        string    `db:login`
	PasswordHash string    `db:password`
	CreateAt     time.Time `db:create_at`
}
