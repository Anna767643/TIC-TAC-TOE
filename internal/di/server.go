package di

import (
	"context"
	"log"
	"net/http"
	"src/internal/datasource"
	"src/internal/web"
	"time"

	"go.uber.org/fx"
)

// Config - параметры сервера из переменных окружения (HTTP_ADDR, по умолчанию ":8080").
type Config struct {
	Addr  string
	DBDSN string
}

// newConfig - задает порт
func newConfig() *Config {
	return &Config{Addr: "localhost:8081",
		DBDSN: "postgres://postgres:12345@localhost:5432/tic_tac_toe"}

}

func newMux(handlerGame *web.GameHandler, handlerUser *web.UserHandler, auth *web.UserAuthenticator) *http.ServeMux {
	mux := http.NewServeMux()
	handlerGame.Register(mux, auth, handlerUser)
	return mux
}

func newServer(cfg *Config, mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// runServer подключает сервер к жизненному циклу fx:
// OnStart - net.Listen + Serve в горутине, OnStop - Shutdown(ctx).
func runServer(lc fx.Lifecycle, srv *http.Server, db *datasource.DBPool) {
	lc.Append(fx.Hook{
		// OnStart выполняется при запуске приложения
		OnStart: func(ctx context.Context) error {
			datasource.CreateTableUsers(ctx, db.GetPool())
			datasource.CreateTableGames(ctx, db.GetPool())
			log.Printf("Starting HTTP server on %s", srv.Addr)
			// Запускаем сервер в отдельной горутине, чтобы не блокировать старт fx
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP server ListenAndServe error: %v", err)
				}
			}()
			return nil
		},

		// OnStop выполняется при завершении приложения (например, по SIGTERM)
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			// Graceful shutdown: дает 5 секунд на завершение текущих запросов
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			err := srv.Shutdown(shutdownCtx)

			db.Close()
			return err
		},
	})
}
