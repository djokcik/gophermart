package service

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/internal/storage"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/djokcik/gophermart/provider"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=OrderService

type OrderService interface {
	ProcessOrder(ctx context.Context, orderID model.OrderID) error
	OrdersByUser(ctx context.Context, userID int) ([]model.Order, error)
	OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error)
	UpdateForAccrual(ctx context.Context, order model.Order, accrual provider.AccrualResponse) error
}

func NewOrderService(cfg config.Config, registry reporegistry.RepoRegistry) OrderService {
	return &orderService{cfg: cfg, repo: registry.GetOrderRepo()}
}

type orderService struct {
	cfg  config.Config
	repo storage.OrderRepository
}

func (o orderService) UpdateForAccrual(ctx context.Context, order model.Order, accrual provider.AccrualResponse) error {
	err := o.repo.UpdateForAccrual(ctx, order, accrual)
	if err != nil {
		o.Log(ctx).Trace().Err(err).Msg("UpdateForAccrual:")
		return err
	}

	return nil
}

func (o orderService) OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	orders, err := o.repo.OrdersByStatus(ctx, status)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("OrdersByStatus:")
		return nil, err
	}

	return orders, nil
}

func (o orderService) OrdersByUser(ctx context.Context, userID int) ([]model.Order, error) {
	orders, err := o.repo.OrdersByUserID(ctx, userID)
	if err != nil {
		o.Log(ctx).Err(err).Msg("OrdersByUser:")
		return nil, err
	}

	return orders, nil
}

func (o orderService) ProcessOrder(ctx context.Context, orderID model.OrderID) error {
	user := appContext.User(ctx)
	if user == nil {
		o.Log(ctx).Err(ErrNotAuthenticated).Msg("")
		return ErrNotAuthenticated
	}

	order, err := o.repo.OrderByID(ctx, orderID)
	if err == nil {
		if user.ID == order.UserID {
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
		ID:     orderID,
		UserID: user.ID,
		Status: model.StatusNew,
	}

	err = o.repo.CreateOrder(ctx, order)
	if err != nil {
		o.Log(ctx).Trace().Err(err).Msg("service: invalid create order")
		return err
	}

	o.Log(ctx).Trace().
		Str("orderID", string(orderID)).
		Msg("success order stored in DB")

	return nil
}

func (o orderService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "orderService").Logger()

	return &logger
}
