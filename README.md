# TIC-TAC-TOE

Реализована классическая игра крестики-нолики, в которую можно играть одному или в компании друга, увы, пока только на localhost :D 

Ниже опишу контракт по API, для удобства взаимодействия с сервисом. 

## Контракт API
### 1. Регистрация новых пользователей
### POST /users
JSON в теле запроса: 
```json
{
  "login": "anna",
  "password": "12345"
}
```
#### Ограничения
- логин и пароль не могут быть пустыми
- логин уникален
### Ответ:
#### Успешно
Status code:
- 201 Created

JSON, представляющий созданного пользователя, в теле ответа:
```json
{
    "id": 1,
    "uuid": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "login": "anna",
    "createAt": "2026-10-27 15:58:49.690819 +0000 UTC"
}
```

#### Ошибка
Status code:
- 400 Bad Request
- 409 Conflict
- 500 Internal Server Error

JSON, описание ошибки

### 2. Авторизация существующих пользователей

### POST /sestions
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
#### Ограничения
- пользователь ранее зарегистрирован

### Ответ:
#### Успешно
Status code:
- 200

JSON, UUID пользователя в теле ответа:
```json
{
    "UUID": "182aca9f-e260-40fc-bfdf-654725db84cd"
}
```

#### Ошибка
Status code:
- 400 Bad Request - не передан логин и пароль
- 401 Unauthorized - не верный логин или пароль
- 500 Internal Server Error

JSON, описание ошибки

### 3. Создание новой игры с пользователем или компьютером

### POST /games
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
JSON в теле запроса:
```json
{
"withComputer": true/false
}
```
#### Ограничения
- пользователь ранее зарегистрирован
### Ответ:
#### Успешно
Status code:
- 201 Created

JSON, представляющий созданный ресурс, в теле ответа:
```json
{
    "id": 10,
    "uuid": "a7187fa4-3bb9-4153-87df-4e1fb889e491",
    "board": [
        [0,0,0],  
        [0,0,0],
        [0,0,0]
    ],
    "with_computer": true,
    "game_state": "next move user 182aca9f-e260-40fc-bfdf-654725db84cd",
    "uuid_player_first": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "uuid_player_second": "",
    "update_at": "2026-10-28 15:12:24.865521 +0000 UTC"
}
```
#### Ошибка
Status code:
- 400 Bad Request
- 500 Internal Server Error

JSON, описание ошибки

### 4. Получение информации о доступных текущих играх

### GET /games?gameState=waitUsers
Query параметры (?gameState=waitUsers)

в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
#### Ограничения
- пользователь ранее зарегистрирован

### Ответ:
#### Успешно
Status code:
- 200 OK

JSON, представляющий все найденные игры, в теле ответа:
```json
[
    {
        "id": 4,
        "uuid": "2b8040c3-14e0-415d-ab2d-61fcfea7e4ed",
        "board": [
            [0,0,0],  
            [0,0,0],
            [0,0,0]
        ],
        "with_computer": false,
        "game_state": "wait users",
        "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
        "uuid_player_second": "",
        "update_at": "2026-10-27 16:09:22.369805 +0000 UTC"
    },
    {
        "id": 5,
        "uuid": "47811e40-5971-4947-bb6f-993c1a42b199",
        "board": [
            [0,0,0],  
            [0,0,0],
            [0,0,0]
        ],
        "with_computer": false,
        "game_state": "wait users",
        "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
        "uuid_player_second": "",
        "update_at": "2026-10-27 16:11:46.006371 +0000 UTC"
    },
    {
        "id": 7,
        "uuid": "e33ac5c1-34fa-4348-8d55-4afa7e45c3bf",
        "board": [
            [0,0,0],  
            [0,0,0],
            [0,0,0]
        ],
        "with_computer": false,
        "game_state": "wait users",
        "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
        "uuid_player_second": "",
        "update_at": "2026-10-28 10:26:58.958947 +0000 UTC"
    }
]
```

#### Ошибка
Status code:
- 400 Bad Request
- 500 Internal Server Error

JSON, описание ошибки

### 5. Присоединение к игре

### PATCH /games/{uuid}/users
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
#### Ограничения
- пользователь ранее зарегистрирован
- игра с таким uuid ранее уже зарегистрирована
- у данной игры нет второга игрока и игра зарегистрирована c параметром "withComputer": false
### Ответ:
#### Успешно
Status code:
- 200 OK

JSON, представляющий новую изменнёную игру, в теле ответа:
```json
{
    "id": 7,
    "uuid": "e33ac5c1-34fa-4348-8d55-4afa7e45c3bf",
    "board": [
        [0,0,0],  
        [0,0,0],
        [0,0,0]
    ],
    "with_computer": false,
    "game_state": "next move user e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_second": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "update_at": "2026-10-28 15:44:08.189064 +0000 UTC"
}
```
#### Ошибка
Status code:
- 400 Bad Request (ErrInvalidUUID)
- 404 Not Found (ErrGameNotFound)
- 409 Conflict (ErrAddUser)
- 500 Internal Server Error (ErrContextCanceled)

JSON, описание ошибки

### 6. Обновление текущей игры (выполнение хода)

### PATCH /games/{uuid}
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
JSON в теле запроса:
```json
	{
		"board": [
			    [0,2,1],
			    [0,0,0],
			    [0,0,0]
			]
	}
```
#### Ограничения
- пользователь зарегистрирован в данной игре
- в данный момент ожидается ход этого пользователя
- первый пользователь ставит только 1
- второй пользователь ставит только 2
### Ответ:
#### Успешно
Status code:
- 200 OK

JSON, представляющий новую изменнёную игру, в теле ответа:
```json
{
    "id": 7,
    "uuid": "e33ac5c1-34fa-4348-8d55-4afa7e45c3bf",
    "board": [
    "board": [
        [0,2,1],  
        [0,0,0],
        [0,0,0]
    ],
    "with_computer": false,
    "game_state": "next move user e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_second": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "update_at": "2026-10-28 15:49:10.952893774 +0700 +07 m=+12.065406485"
}
```

#### Ошибка
Status code:
- 400 Bad Request (поле нарушает правила или изменены прошлые ходы, )
- 404 Not Found (игра не найдена)
- 409 Conflict (ход аппонента или игра ранее уже была завершена, недопустимое значение клетки)
- 500 Internal Server Error

JSON, описание ошибки

### 7. Получение информации о текущей игре

### GET /games/{uuid}
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
#### Ограничения
- пользователь зарегистрирован в данной игре
### Ответ:
#### Успешно
Status code:
- 200 OK

JSON, представляющий игру, в теле ответа:
```json
{
    "id": 7,
    "uuid": "e33ac5c1-34fa-4348-8d55-4afa7e45c3bf",
    "board": [
    "board": [
        [0,2,1],  
        [0,0,0],
        [0,0,0]
    ],
    "with_computer": false,
    "game_state": "next move user e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_first": "e4432ee8-26ba-43c4-b4f9-88e547798007",
    "uuid_player_second": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "update_at": "2026-10-28 15:49:10.952893774 +0700 +07 m=+12.065406485"
}
```

#### Ошибка
Status code:
- 400 Bad Request
- 404 Not Found
- 500 Internal Server Error

JSON, описание ошибки

### 8. Получение информации о пользователе по UUID

### GET /users/{uuid}
в Header логин и пароль в виде base64(login:password)
```
Key: Authorization Value: Basic YW5uYToxMjM0NQ==
```
#### Ограничения
- можно получить информацию о любом зарегестрированном пользователе
### Ответ:
#### Успешно
Status code:
- 200 OK

JSON, представляющий информацию о пользователе, в теле ответа:
```json
{
    "id": 1,
    "uuid": "182aca9f-e260-40fc-bfdf-654725db84cd",
    "login": "anna",
    "createAt": "2026-10-27 15:58:49.690819 +0000 UTC"
}
```

#### Ошибка
Status code:
- 400 Bad Request
- 404 Not Found
- 500 Internal Server Error

JSON, описание ошибки
