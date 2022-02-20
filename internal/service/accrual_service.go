package service

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/djokcik/gophermart/provider"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=AccrualService

type AccrualService interface {
	Poller(ctx context.Context) func()
}

func NewAccrualService(cfg config.Config, registry reporegistry.RepoRegistry) AccrualService {
	return &accrualService{
		client: provider.NewAccrualClient(cfg),
		order:  NewOrderService(cfg, registry),
	}
}

type accrualService struct {
	client provider.AccrualClient
	order  OrderService
}

func (a accrualService) Poller(ctx context.Context) func() {
	return func() {
		orders := a.getOrders(ctx)

		for _, order := range orders {
			a.ProcessOrder(ctx, order)
		}
	}
}

func (a accrualService) ProcessOrder(ctx context.Context, order model.Order) {
	response, err := a.client.GetOrder(ctx, order.ID)
	if err != nil {
		var apiErr *provider.ErrAccrualResponse
		if errors.As(err, &apiErr) {
			a.Log(ctx).Warn().Err(apiErr).Msg("ProcessOrder:")
			return
		}

		a.Log(ctx).Error().Err(err).Msg("ProcessOrder:")
		return
	}

	err = a.order.UpdateForAccrual(ctx, order, response)
	if err != nil {
		a.Log(ctx).Error().Err(err).Msg("UpdateForAccrual:")
		return
	}
}

func (a accrualService) getOrders(ctx context.Context) []model.Order {
	orders := make([]model.Order, 0)
	for _, status := range []model.Status{model.StatusNew, model.StatusProcessing} {
		ords, err := a.order.OrdersByStatus(ctx, status)
		if err != nil {
			a.Log(ctx).Err(err).Msg("getOrders: failed get orders by status")
			continue
		}

		orders = append(orders, ords...)
	}

	return orders
}

func (a accrualService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "accrualService").Logger()

	return &logger
}
