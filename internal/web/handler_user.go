package web

import (
	"encoding/json"
	"net/http"
	"src/internal/domain"
)

// NewUserHandler - конструктор для di
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// UserHandler - HTTP-обработчик
type UserHandler struct {
	userService domain.UserService
}

/*
HandleCreateUser - Хэндел для регистрация нового пользователя

Method POST /users
JSON в теле запроса: login и passwort

Ответ:
Успешно
Status code:
- 201 Created
JSON, представляющий созданного пользователя, в теле ответа

Ошибка
Status code:
- 400 Bad Request
- 409 Conflict
- 500 Internal Server Error
JSON, описание ошибки
*/
func (u *UserHandler) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var signUpRequest domain.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpRequest); err != nil {
		writeError(w, http.StatusBadRequest, "wrong body, need login and password in json, error: "+err.Error())
		return
	}

	if signUpRequest.Login == "" || signUpRequest.Password == "" {
		writeError(w, http.StatusBadRequest, "wrong body, login or password not be empty")
		return
	}

	// В зависимости от ошибки - разные ответы, либо Conflict, либо Internal Server Error
	user, err := u.userService.Register(ctx, signUpRequest)
	if err == domain.ErrLoginIncorrect {
		writeError(w, http.StatusConflict, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userDTO := toUserDTO(user)
	writeJSON(w, http.StatusCreated, userDTO)
}

/*
handleSessionCreate - создание сессии (отправка UUID авторизованного пользователя)

Method POST /sestions
в Header логин и пароль в виде base64(login:password)

Ответ:
Успешно
Status code:
- 200
JSON, UUID пользователя в теле ответа

Ошибка
Status code:
- 400 Bad Request - не передан логин и пароль
- 401 Unauthorized - не верный логин или пароль
- 500 Internal Server Error
JSON, описание ошибки
*/
func (u *UserHandler) handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	login, password, ok := r.BasicAuth()
	if !ok {
		writeError(w, http.StatusBadRequest, "need authorization")
		return
	}

	signUpRequest := domain.SignUpRequest{
		Login:    login,
		Password: password,
	}

	uuid, err := u.userService.LogIn(ctx, signUpRequest)
	if err == domain.ErrLoginOrPasswordIncorrect {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type response struct {
		UUID string
	}

	res := response{UUID: uuid}
	writeJSON(w, http.StatusOK, res)
}

/*
handleGetUserInfo - получение информации о пользователе по UUID
Method GET /users/{uuid}

Ответ:
Успешно
Status code:
- 200 OK
JSON, представляющий информацию о пользователе, в теле ответа

Ошибка
Status code:
- 400 Bad Request
- 404 Not Found
- 500 Internal Server Error
JSON, описание ошибки
*/
func (u *UserHandler) handleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.PathValue("uuid")

	uuidUser, err := toUUID(userId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := u.userService.GetUserInfo(ctx, uuidUser)

	if err == domain.ErrUserNotFound {
		writeError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userDTO := toUserDTO(user)
	writeJSON(w, http.StatusCreated, userDTO)

}
