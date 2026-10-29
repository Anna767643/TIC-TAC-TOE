package di

import (
	"src/internal/datasource"
	"src/internal/domain"
	"src/internal/web"

	"go.uber.org/fx"
)

// Module - полный граф зависимостей приложения
var Module = fx.Options(
	fx.Provide(
		newConfig,
		func(cfg *Config) string {
			return cfg.DBDSN
		},
		datasource.NewDBPool,
		datasource.NewGameRepository,
		datasource.NewUserRepository,
		domain.NewGameService,
		domain.NewUserService,
		web.NewGameHandler,
		web.NewUserHandler,
		web.NewUserAuthenticator,
		newMux,
		newServer,
	),
	fx.Invoke(runServer),
)
