package handler

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/internal/service"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

var (
	CookieName = "gophermart_token"
)

type Handler struct {
	*chi.Mux
	user     service.UserService
	order    service.OrderService
	withdraw service.WithdrawService
}

func NewHandler(mux *chi.Mux, cfg config.Config, repoRegistry reporegistry.RepoRegistry) *Handler {
	return &Handler{
		Mux:      mux,
		user:     service.NewUserService(cfg, repoRegistry),
		order:    service.NewOrderService(cfg, repoRegistry),
		withdraw: service.NewWithdrawService(cfg, repoRegistry),
	}
}

func (h *Handler) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "Handler").Logger()

	return &logger
}
