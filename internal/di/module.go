package di

import (
	"src/internal/datasource"
	"src/internal/domain"
	"src/internal/web"

	"go.uber.org/fx"
)

// Module - полный граф зависимостей приложения.
var Module = fx.Options(
	fx.Provide(
		newConfig, // сервер с портоми данными для подключения к БД
		func(cfg *Config) string {
			return cfg.DBDSN
		},
		datasource.NewDBPool,         // Создание poll подключения к БД
		datasource.NewGameRepository, // domain.GameRepository
		datasource.NewUserRepository,
		domain.NewGameService, // domain.GameService
		domain.NewUserService,
		web.NewGameHandler, // *web.GameHandler
		web.NewUserHandler,
		web.NewUserAuthenticator,
		newMux,    // *http.ServeMux с зарегистрированными маршрутами
		newServer, // *http.Server с таймаутами

	),
	fx.Invoke(runServer),
)
