package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"src/internal/domain"

	uuid "github.com/google/uuid"
)

// GameHandler - HTTP-обработчик
type GameHandler struct {
	gameService domain.GameService
}

// NewGameHandler - конструктор для di
func NewGameHandler(gameService domain.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// ----------------!!!!!!!!!!!!------------------------------Подумать как удобно сделать,
// ведь зависимость теперь от GameHandler странная, потому что в равной степени еще зависит от UserAuthenticator и UserHandler
// Register - регистрирует маршруты
func (h *GameHandler) Register(mux *http.ServeMux, auth *UserAuthenticator, u *UserHandler) {
	mux.HandleFunc("POST /users", u.handleUsersCreate)
	mux.HandleFunc("POST /sessions", u.handleSessionCreate)
	mux.HandleFunc("GET /users/{uuid}", auth.Validate(u.handleGetUserInfo))

	mux.HandleFunc("POST /games", auth.Validate(h.handleCreateGame))
	mux.HandleFunc("GET /games", auth.Validate(h.handleGetAvailableGames))
	mux.HandleFunc("GET /games/{uuid}", auth.Validate(h.handleGetGameInfo))
	mux.HandleFunc("PATCH /games/{uuid}/users", auth.Validate(h.handleJoinUser))
	mux.HandleFunc("PATCH /games/{uuid}", auth.Validate(h.handlePatchCurrentGame))
}

/*
handleCreateGame - создание новой игры с пользователем или с компьютером

Method POST /games
JSON в теле запроса:
{
"withComputer": true/false
}

Ответ:
Успешно
Status code:
- 201 Created
JSON, представляющий созданный ресурс, в теле ответа

Ошибка
Status code:
- 400 Bad Request
- 500 Internal Server Error
JSON, описание ошибки
*/
func (h *GameHandler) handleCreateGame(w http.ResponseWriter, r *http.Request) {

	ctx, uuidUser, err := getUUIDUser(w, r)
	if err != nil {
		return
	}

	var req reqNewGame
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "wrong body, need param with_computer, error: "+err.Error())
		return
	}

	if req.WithComputer == nil {
		writeError(w, http.StatusBadRequest, "missing required field: with_computer")
		return
	}

	domainGame, err := h.gameService.CreateNewGame(ctx, *req.WithComputer, uuidUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newGame := toGameDTO(domainGame)

	writeJSON(w, http.StatusCreated, newGame)
}

/*
handleGetAvailableGames - получение информации о доступных играх
Method GET /games
Query параметры (?gameState=waitUsers)

Ответ:
Успешно
Status code:
- 200 OK
JSON, представляющий все найденные игры, в теле ответа

Ошибка
Status code:
- 400 Bad Request
- 500 Internal Server Error
JSON, описание ошибки
*/
func (h *GameHandler) handleGetAvailableGames(w http.ResponseWriter, r *http.Request) {
	ctx, uuidUser, err := getUUIDUser(w, r)
	if err != nil {
		return
	}

	queryParam := r.FormValue("gameState")
	if queryParam != "waitUsers" {
		writeError(w, http.StatusBadRequest, "need query param gameState")
		return
	}

	games, err := h.gameService.GetAvailableGames(ctx, uuidUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	gamesDTO := make([]gameDTO, 0)
	for i := range games {
		gameDTO := toGameDTO(games[i])
		gamesDTO = append(gamesDTO, gameDTO)
	}

	writeJSON(w, http.StatusOK, gamesDTO)
}

/*
handleJoinUser - присоединение пользователя к игре
Method PATCH /games/{uuid}/users

Ответ:
Успешно
Status code:
- 200 OK
JSON, представляющий новую изменнёную игру, в теле ответа

Ошибка
Status code:
- 400 Bad Request (ErrInvalidUUID)
- 404 Not Found (ErrGameNotFound)
- 409 Conflict (ErrAddUser)
- 500 Internal Server Error (ErrContextCanceled)
JSON, описание ошибки
*/
func (h *GameHandler) handleJoinUser(w http.ResponseWriter, r *http.Request) {
	ctx, uuidUser, err := getUUIDUser(w, r)
	if err != nil {
		return
	}

	gameId := r.PathValue("uuid")

	uuidGame, err := toUUID(gameId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	game, err := h.gameService.JoinUser(ctx, uuidUser, uuidGame)
	if err == domain.ErrGameNotFound {
		writeError(w, http.StatusNotFound, "error to add user: "+err.Error())
		return
	} else if err == domain.ErrAddUser {
		writeError(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	gameDTO := toGameDTO(game)
	writeJSON(w, http.StatusOK, gameDTO)
}

/*
handlePatchCurrentGame - обновление текущей игры
с учетом игры с пользователем или с компьютером
Method PATCH /games/{uuid}
JSON в теле запроса:

	{
		"board": [
			    [0,0,1],
			    [0,0,0],
			    [0,0,0]
			]
	}

Ответ:
Успешно
Status code:
- 200 OK
JSON, представляющий новую изменнёную игру, в теле ответа

Ошибка
Status code:
- 400 Bad Request (поле нарушает правила или изменены прошлые ходы, )
- 404 Not Found (игра не найдена)
- 409 Conflict (ход аппонента или игра ранее уже была завершена, недопустимое значение клетки)
- 500 Internal Server Error
JSON, описание ошибки
*/
func (h *GameHandler) handlePatchCurrentGame(w http.ResponseWriter, r *http.Request) {
	ctx, uuidUser, err := getUUIDUser(w, r)
	if err != nil {
		return
	}

	gameId := r.PathValue("uuid")

	uuidGame, err := toUUID(gameId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var rboard reqBoard

	if err := json.NewDecoder(r.Body).Decode(&rboard); err != nil {
		writeError(w, http.StatusBadRequest, "wrong body "+err.Error())
		return
	}

	game, err := h.gameService.NextMove(ctx, domain.Board(rboard.BoardDTO), uuidGame, uuidUser)
	if err == domain.ErrGameNotFound {

		writeError(w, http.StatusNotFound, err.Error())
		return
	} else if err == domain.ErrInvalidBoard || err == domain.ErrInvalidCell || err == domain.ErrDoesntHaveSuchUser {

		writeError(w, http.StatusBadRequest, "wrong body "+err.Error())
		return
	} else if err == domain.ErrOpponentsMove || err == domain.ErrGameFinished {

		writeError(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if game.State == domain.StateDraw {

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "GAME OVER",
			"winner":  "DRAW",
		})

	} else if game.State == domain.StateWin {

		writeJSON(w, http.StatusOK, map[string]string{
			"message": "GAME OVER",
			"winner":  game.Winner.String(),
		})

	} else {
		gameDTO := toGameDTO(game)
		writeJSON(w, http.StatusOK, gameDTO)
	}
}

/*
handleGetGameInfo - получение информации об игре
Method GET /games/{uuid}

Ответ:
Успешно
Status code:
- 200 OK
JSON, представляющий игру, в теле ответа

Ошибка
Status code:
- 400 Bad Request
- 404 Not Found
- 500 Internal Server Error
JSON, описание ошибки
*/
func (h *GameHandler) handleGetGameInfo(w http.ResponseWriter, r *http.Request) {
	ctx, uuidUser, err := getUUIDUser(w, r)
	if err != nil {
		return
	}

	gameId := r.PathValue("uuid")

	uuidGame, err := toUUID(gameId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	game, err := h.gameService.GetGame(ctx, uuidGame, uuidUser)

	if err == domain.ErrGameNotFound {

		writeError(w, http.StatusNotFound, err.Error())
		return
	} else if err == domain.ErrDoesntHaveSuchUser {

		writeError(w, http.StatusBadRequest, "wrong body "+err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	gameDTO := toGameDTO(game)
	writeJSON(w, http.StatusOK, gameDTO)
}

// getUUIDUser - получает user uuidиз контекста
func getUUIDUser(w http.ResponseWriter, r *http.Request) (context.Context, uuid.UUID, error) {
	ctx := r.Context()

	uuidUser, ok := ctx.Value(userUUIDKey).(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "need authorization")
		return ctx, uuid.UUID{}, errors.New("error get user uuid")
	}

	return ctx, uuidUser, nil
}
