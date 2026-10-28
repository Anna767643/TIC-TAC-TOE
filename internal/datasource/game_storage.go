package datasource

import (
	"context"
	"fmt"
	"src/internal/domain"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameStorage interface {
	Close()
	GetPool() *pgxpool.Pool
	insertIntoGames(ctx context.Context, game gameEntity) error
	insertIntoUsers(ctx context.Context, user userEntity) error
	selectAvaibleGamesFromGames(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error)
	selectFromGames(ctx context.Context, uuid uuid.UUID) (gameEntity, error)
	selectFromUsersWithLogin(ctx context.Context, login string) (userEntity, error)
	selectFromUsersWithUUID(ctx context.Context, uuidUser uuid.UUID) (userEntity, error)
	updateGames(ctx context.Context, game gameEntity) error
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

func (db *DBPool) Close() {
	db.pool.Close()
}

func (db *DBPool) GetPool() *pgxpool.Pool {
	return db.pool
}

// insertTableGames -добавляет новую игру в таблице games
func (db *DBPool) insertIntoGames(ctx context.Context, game gameEntity) error {
	sqlQuery := `
	INSERT INTO games(uuid, board, with_computer, game_state, uuid_player_first, update_at)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := db.pool.Exec(ctx, sqlQuery, game.UUID, game.Board, game.WithComputer, game.GameState, game.UUIDPlayerFirst, game.UpdatedAt)
	return err
}

// updateTableGames - обновляет информацию об игре в таблице games
func (db *DBPool) updateGames(ctx context.Context, game gameEntity) error {
	sqlQuery := `
	UPDATE games
	SET board = $1, game_state = $2, uuid_player_second = $3, update_at = $4
	WHERE uuid = $5;
	`

	_, err := db.pool.Exec(ctx, sqlQuery, game.Board, game.GameState, game.UUIDPlayerSecond, game.UpdatedAt, game.UUID)
	return err
}

// selectAvaibleGamesFromGames - получение среза игр из таблицы games
// доступных для присоединения игроку с uuid = uuidUser
func (db *DBPool) selectAvaibleGamesFromGames(ctx context.Context, uuidUser uuid.UUID) ([]gameEntity, error) {
	sqlQuery := `
	SELECT id, uuid, board, with_computer, game_state, uuid_player_first, uuid_player_second, update_at from games
	WHERE uuid_player_first != $1 and game_state = 'wait users' and with_computer = false;
	`

	rows, err := db.pool.Query(ctx, sqlQuery, uuidUser)
	if err != nil {
		return []gameEntity{}, err
	}
	defer rows.Close()
	games := make([]gameEntity, 0)

	for rows.Next() {
		var game gameEntity
		if err := rows.Scan(&game.ID, &game.UUID, &game.Board, &game.WithComputer, &game.GameState, &game.UUIDPlayerFirst, &game.UUIDPlayerSecond, &game.UpdatedAt); err != nil {
			return []gameEntity{}, err
		}
		games = append(games, game)
	}

	return games, nil
}

// selectFromGames - получение информации об игре по UUID игры из таблицы games
func (db *DBPool) selectFromGames(ctx context.Context, uuid uuid.UUID) (gameEntity, error) {
	sqlQuery := `
	SELECT id, uuid, board, with_computer, game_state, uuid_player_first, uuid_player_second, update_at from games
	WHERE uuid = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, uuid)
	var game gameEntity

	if err := row.Scan(&game.ID, &game.UUID, &game.Board, &game.WithComputer, &game.GameState, &game.UUIDPlayerFirst, &game.UUIDPlayerSecond, &game.UpdatedAt); err != nil {
		return gameEntity{}, err
	}
	return game, nil
}

// ---- USERS ----
// insertPlayer - добаление нового играка в БД
func (db *DBPool) insertIntoUsers(ctx context.Context, user userEntity) error {
	_, err := db.selectFromUsersWithLogin(ctx, user.Login)
	if err == nil {
		return domain.ErrLoginIncorrect
	}

	sqlQuery := `
	INSERT INTO users(uuid, login, password_hash, create_at)
	VALUES ($1, $2, $3, $4);
	`

	_, err = db.pool.Exec(ctx, sqlQuery, user.UUID, user.Login, user.PasswordHash, user.CreateAt)
	return err
}

// selectUser - получение информации об игроке по логину
func (db *DBPool) selectFromUsersWithLogin(ctx context.Context, login string) (userEntity, error) {
	sqlQuery := `
	SELECT id, uuid, login, password_hash, create_at from users
	WHERE login = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, login)
	var user userEntity

	if err := row.Scan(&user.ID, &user.UUID, &user.Login, &user.PasswordHash, &user.CreateAt); err != nil {
		return userEntity{}, err
	}
	return user, nil
}

func (db *DBPool) selectFromUsersWithUUID(ctx context.Context, uuidUser uuid.UUID) (userEntity, error) {
	sqlQuery := `
	SELECT id, uuid, login, password_hash, create_at from users
	WHERE uuid = $1;
	`

	row := db.pool.QueryRow(ctx, sqlQuery, uuidUser)
	var user userEntity

	if err := row.Scan(&user.ID, &user.UUID, &user.Login, &user.PasswordHash, &user.CreateAt); err != nil {
		return userEntity{}, err
	}
	return user, nil
}
