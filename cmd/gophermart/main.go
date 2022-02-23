package main

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/internal/service"
	helpers "github.com/djokcik/gophermart/pkg/helper"
	"github.com/djokcik/gophermart/pkg/logging"
	serverMiddleware "github.com/djokcik/gophermart/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.NewConfig()

	logging.
		NewLogger().
		Info().
		Msgf("config: %+v", cfg)

	mux := chi.NewMux()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(serverMiddleware.GzipHandle)
	mux.Use(serverMiddleware.LoggerMiddleware())

	repoRegistry, err := reporegistry.NewPostgreSQL(ctx, cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msgf("Doesn`t open database connection")
		os.Exit(1)
	}

	accrualService := service.NewAccrualService(cfg, repoRegistry)
	go helpers.SetTicker(accrualService.Poller(ctx), 5*time.Second)

	makeMetricRoutes(ctx, mux, cfg, repoRegistry)

	go func() {
		err := http.ListenAndServe(cfg.Address, mux)
		if err != nil {
			logging.NewLogger().Fatal().Err(err).Msg("server stopped")
		}

	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	cancel()
	logging.NewLogger().Info().Msg("Shutdown Server ...")
}
