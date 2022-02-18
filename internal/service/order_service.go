package service

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=OrderService

type OrderService interface {
	ProcessOrder(ctx context.Context, orderId model.OrderId) error
}

func NewOrderService(cfg config.Config, repo storage.OrderRepository) OrderService {
	return &orderService{cfg: cfg, repo: repo}
}

type orderService struct {
	cfg  config.Config
	repo storage.OrderRepository
}

func (o orderService) ProcessOrder(ctx context.Context, orderId model.OrderId) error {
	user := appContext.User(ctx)
	if user == nil {
		o.Log(ctx).Err(ErrNotAuthenticated).Msg("")
		return ErrNotAuthenticated
	}

	order, err := o.repo.FindOrderById(ctx, orderId)
	if err == nil {
		if user.Id == order.UserId {
			o.Log(ctx).Trace().Err(ErrOrderAlreadyUploaded).Msg("")
			return ErrOrderAlreadyUploaded
		}

		o.Log(ctx).Trace().Err(ErrOrderAlreadyUploadedAnotherUser).Msg("")
		return ErrOrderAlreadyUploadedAnotherUser
	}

	if !errors.Is(err, storage.ErrNotFound) {
		o.Log(ctx).Trace().Err(err).Msg("")
		return err
	}

	order = model.Order{
		Id:      orderId,
		UserId:  user.Id,
		Status:  model.StatusNew,
		Accrual: 0,
	}

	err = o.repo.CreateOrder(ctx, order)
	if err != nil {
		o.Log(ctx).Trace().Err(err).Msg("service: invalid create order")
		return err
	}

	o.Log(ctx).Trace().
		Str("orderId", string(orderId)).
		Msg("success order stored in DB")

	return nil
}

func (o orderService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "orderService").Logger()

	return &logger
}
