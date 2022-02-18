package main

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/handler"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/djokcik/gophermart/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"os"
)

func makeMetricRoutes(ctx context.Context, mux *chi.Mux, cfg config.Config) *handler.Handler {
	repoRegistry, err := reporegistry.NewPostgreSQL(ctx, cfg)
	if err != nil {
		logging.NewLogger().Fatal().Err(err).Msgf("Doesn`t open database connection")
		os.Exit(1)
	}

	h := handler.NewHandler(mux, cfg, repoRegistry)

	h.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUserHandler())
		r.Post("/login", h.SignInHandler())

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.UserContext(repoRegistry.GetUserRepo(), cfg))

			r.Post("/orders", h.UploadOrderHandler())
			r.Get("/orders", h.GetOrdersHandler())
		})
	})

	return h
}
