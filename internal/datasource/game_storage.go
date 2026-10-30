package datasource

import (
	"context"
	"errors"
	"fmt"
	"src/internal/domain"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameStorage interface {
	// Close - закрывает соединение
	Close()

	// GetPool - передает соединение
	GetPool() *pgxpool.Pool

	// insertIntoGames - добавляет новую игру в таблице games
	// "repository: insert game %s: %w"
	insertIntoGames(ctx context.Context, game gameEntity) error

	// updateGames - обновляет информацию об игре в таблице games
	// "repository: update game %s: %w"
	updateGames(ctx context.Context, game gameEntity) error

	// selectAvaibleGamesFromGames - получение среза игр из таблицы games
	// доступных для присоединения игроку с uuid = uuidUser
	// "repository: select game for user %s: %w"
	// "repository: scan game for user %s: %w"
	selectAvaibleGamesFromGames(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error)

	// selectFromGames - получение информации об игре по UUID игры из таблицы games
	// ErrGameNotFound
	// "repository: select game %s: %w"
	selectFromGames(ctx context.Context, uuid uuid.UUID) (gameEntity, error)

	// insertPlayer - добаление нового играка в БД
	// ErrLoginIncorrect
	// "repository: insert user with %s: %w"
	insertIntoUsers(ctx context.Context, user userEntity) error

	// selectFromUsersWithLogin - получение информации об игроке по логину
	// ErrUserNotFound
	// "repository: select user with %s: %w"
	selectFromUsersWithLogin(ctx context.Context, login string) (userEntity, error)

	// selectFromUsersWithUUID - получение информации об игроке по uuid
	// ErrUserNotFound
	// "repository: select user with %s: %w"
	selectFromUsersWithUUID(ctx context.Context, uuidUser uuid.UUID) (userEntity, error)
}

// DBPool - соединение с БД, потокобезопасен
// Реализация gameStorage
type DBPool struct {
	pool *pgxpool.Pool
}

// NewDBPool - конструктор для di
func NewDBPool(dsn string) (*DBPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пула: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("не удалось пингануть БД: %w", err)
	}

	return &DBPool{pool: pool}, nil
}

// Close - закрывает соединение
func (db *DBPool) Close() {
	db.pool.Close()
}

// GetPool - передает соединение
func (db *DBPool) GetPool() *pgxpool.Pool {
	return db.pool
}

// ------------- GAMES -------------

// insertIntoGames - добавляет новую игру в таблице games
// "repository: insert game %s: %w"
func (db *DBPool) insertIntoGames(ctx context.Context, game gameEntity) error {
	sqlQuery := `
	INSERT INTO games(uuid, board, with_computer, state, next_player, winner, uuid_player_first, update_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err := db.pool.Exec(ctx, sqlQuery, game.UUID, game.Board, game.WithComputer, game.State, game.NextPlayer, game.Winner, game.UUIDPlayerFirst, game.UpdatedAt)
	if err != nil {
		return fmt.Errorf("repository: insert game %s: %w", game.UUID, err)
	}

	return nil
}

// updateGames - обновляет информацию об игре в таблице games
// "repository: update game %s: %w"
func (db *DBPool) updateGames(ctx context.Context, game gameEntity) error {
	sqlQuery := `
	UPDATE games
	SET board = $1, state = $2, next_player = $3, winner = $4, uuid_player_second = $5, update_at = $6
	WHERE uuid = $7;
	`

	_, err := db.pool.Exec(ctx, sqlQuery, game.Board, game.State, game.NextPlayer, game.Winner, game.UUIDPlayerSecond, game.UpdatedAt, game.UUID)
	if err != nil {
		return fmt.Errorf("repository: update game %s: %w", game.UUID, err)
	}

	return nil
}

// selectAvaibleGamesFromGames - получение среза игр из таблицы games
// доступных для присоединения игроку с uuid = uuidUser
// "repository: select game for user %s: %w"
// "repository: scan game for user %s: %w"
func (db *DBPool) selectAvaibleGamesFromGames(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
	sqlQuery := `
	SELECT id, uuid, board, with_computer, state, next_player, winner, uuid_player_first, uuid_player_second, update_at from games
	WHERE uuid_player_first != $1 and state = 'wait users' and with_computer = false;
	`

	rows, err := db.pool.Query(ctx, sqlQuery, uuidUser)
	if err != nil {
		return []gameEntity{}, fmt.Errorf("repository: select game for user %s: %w", uuidUser, err)
	}
	defer rows.Close()
	games := make([]gameEntity, 0)

	for rows.Next() {
		var game gameEntity
		if err := rows.Scan(&game.ID, &game.UUID, &game.Board, &game.WithComputer, &game.State, &game.NextPlayer, &game.Winner, &game.UUIDPlayerFirst, &game.UUIDPlayerSecond, &game.UpdatedAt); err != nil {
			return []gameEntity{}, fmt.Errorf("repository: scan game for user %s: %w", uuidUser, err)
		}
		games = append(games, game)
	}

	return games, nil
}

// selectFromGames - получение информации об игре по UUID игры из таблицы games
// ErrGameNotFound
// "repository: select game %s: %w"
func (db *DBPool) selectFromGames(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
	sqlQuery := `
	SELECT id, uuid, board, with_computer, state, next_player, winner, uuid_player_first, uuid_player_second, update_at from games
	WHERE uuid = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, uuid)
	var game gameEntity

	if err := row.Scan(&game.ID, &game.UUID, &game.Board, &game.WithComputer, &game.State, &game.NextPlayer, &game.Winner, &game.UUIDPlayerFirst, &game.UUIDPlayerSecond, &game.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return gameEntity{}, domain.ErrGameNotFound
		}
		return gameEntity{}, fmt.Errorf("repository: select game %s: %w", uuid, err)
	}
	return game, nil
}

// ------------- USERS -------------

// insertPlayer - добаление нового играка в БД
// ErrLoginIncorrect
// "repository: insert user with %s: %w"
func (db *DBPool) insertIntoUsers(ctx context.Context, user userEntity) error {

	sqlQuery := `
	INSERT INTO users(uuid, login, password_hash, create_at)
	VALUES ($1, $2, $3, $4);
	`

	_, err := db.pool.Exec(ctx, sqlQuery, user.UUID, user.Login, user.PasswordHash, user.CreateAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation { // "23505"
		return domain.ErrLoginIncorrect
	}
	if err != nil {
		return fmt.Errorf("repository: insert user with %s: %w", user.Login, err)
	}

	return nil
}

// selectFromUsersWithLogin - получение информации об игроке по логину
// ErrUserNotFound
// "repository: select user with %s: %w"
func (db *DBPool) selectFromUsersWithLogin(ctx context.Context, login string) (userEntity, error) {
	sqlQuery := `
	SELECT id, uuid, login, password_hash, create_at from users
	WHERE login = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, login)
	var user userEntity

	if err := row.Scan(&user.ID, &user.UUID, &user.Login, &user.PasswordHash, &user.CreateAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userEntity{}, domain.ErrUserNotFound
		}
		return userEntity{}, fmt.Errorf("repository: select user with %s: %w", login, err)
	}

	return user, nil
}

// selectFromUsersWithUUID - получение информации об игроке по uuid
// ErrUserNotFound
// "repository: select user with %s: %w"
func (db *DBPool) selectFromUsersWithUUID(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
	sqlQuery := `
	SELECT id, uuid, login, password_hash, create_at from users
	WHERE uuid = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, uuidUser)
	var user userEntity

	if err := row.Scan(&user.ID, &user.UUID, &user.Login, &user.PasswordHash, &user.CreateAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userEntity{}, domain.ErrUserNotFound
		}
		return userEntity{}, fmt.Errorf("repository: select user with %s: %w", uuidUser, err)
	}
	return user, nil
}
