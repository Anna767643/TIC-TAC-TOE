package datasource

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateTableUsers -создает таблицу пользователей, если её еще нет
// Возвращает ошибку, если они возникли при выполнении запроса
func CreateTableUsers(ctx context.Context, pool *pgxpool.Pool) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS users (
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid UUID NOT NULL UNIQUE,
	login VARCHAR(255) NOT NULL UNIQUE,
	password_hash VARCHAR(255) NOT NULL,
	create_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := pool.Exec(ctx, sqlQuery)

	return err
}

// CreateTableGames -создает таблицу игр, если её еще нет
// Возвращает ошибку, если они возникли при выполнении запроса
func CreateTableGames(ctx context.Context, pool *pgxpool.Pool) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS games (
	id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	uuid UUID NOT NULL UNIQUE,
	board VARCHAR(255) NOT NULL,
	with_computer BOOLEAN NOT NULL DEFAULT FALSE,
	state VARCHAR(16) NOT NULL,
	next_player UUID,
	winner UUID,

	uuid_player_first UUID NOT NULL,
	uuid_player_second UUID,

	update_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (uuid_player_first) REFERENCES users (uuid),
	FOREIGN KEY (uuid_player_second) REFERENCES users (uuid)
	);
	`

	_, err := pool.Exec(ctx, sqlQuery)

	return err
}
