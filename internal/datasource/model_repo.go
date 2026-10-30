package datasource

import (
	"src/internal/domain"
	"time"

	uuid "github.com/google/uuid"
)

// gameEntity представляет запись в таблице games
type gameEntity struct {
	ID           int              `db:"id"`
	UUID         uuid.UUID        `db:"uuid"`
	Board        string           `db:"board"`
	WithComputer bool             `db:"with_computer"`
	GameState    string           `db:"game_state"`
	State        domain.GameState `db:"state"`
	NextPlayer   *uuid.UUID       `db:"next_player"`
	Winner       *uuid.UUID       `db:"winner"`

	UUIDPlayerFirst  uuid.UUID  `db:"uuid_player_first"`  // X
	UUIDPlayerSecond *uuid.UUID `db:"uuid_player_second"` // O

	UpdatedAt time.Time `db:"update_at"`
}

// userEntity представляет запись в таблице players
type userEntity struct {
	ID           int       `db:"id"`
	UUID         uuid.UUID `db:"uuid"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	CreateAt     time.Time `db:"create_at"`
}
