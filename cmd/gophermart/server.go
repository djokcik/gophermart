package main

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/handler"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/pkg/middleware"
	"github.com/go-chi/chi/v5"
)

func makeMetricRoutes(_ context.Context, mux *chi.Mux, cfg config.Config, registry reporegistry.RepoRegistry) *handler.Handler {
	h := handler.NewHandler(mux, cfg, registry)

	h.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.RegisterUserHandler())
		r.Post("/login", h.SignInHandler())

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.UserContext(registry.GetUserRepo(), cfg))

			r.Post("/orders", h.UploadOrderHandler())
			r.Get("/orders", h.GetOrdersHandler())
			r.Get("/balance", h.GetBalanceHandler())
			r.Post("/balance/withdraw", h.WithdrawHandler())
			r.Get("/withdrawals", h.WithdrawLogsHandler())
		})
	})

	return h
}
