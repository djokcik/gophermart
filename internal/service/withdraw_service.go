package service

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/internal/storage"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=WithdrawService

type WithdrawService interface {
	ProcessWithdraw(ctx context.Context, orderId model.OrderId, sum model.Amount) error
	WithdrawLogsByUserId(ctx context.Context, userId int) ([]model.Withdraw, error)
	AmountWithdrawByUser(ctx context.Context, userId int) (model.Amount, error)
}

func NewWithdrawService(cfg config.Config, registry reporegistry.RepoRegistry) WithdrawService {
	return &withdrawService{cfg: cfg, repo: registry.GetWithdrawRepo()}
}

type withdrawService struct {
	cfg  config.Config
	repo storage.WithdrawRepository
}

func (o withdrawService) AmountWithdrawByUser(ctx context.Context, userId int) (model.Amount, error) {
	amount, err := o.repo.AmountWithdrawByUser(ctx, userId)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("AmountWithdrawByUser:")
		return 0, err
	}

	return amount, nil
}

func (o withdrawService) WithdrawLogsByUserId(ctx context.Context, userId int) ([]model.Withdraw, error) {
	withdrawLogs, err := o.repo.WithdrawLogsByUserId(ctx, userId)
	if err != nil {
		o.Log(ctx).Error().Err(err).Msg("WithdrawLogsByUserId:")
		return nil, err
	}

	return withdrawLogs, nil
}

func (o withdrawService) ProcessWithdraw(ctx context.Context, orderId model.OrderId, sum model.Amount) error {
	user := appContext.User(ctx)
	if user == nil {
		o.Log(ctx).Err(ErrNotAuthenticated).Msg("")
		return ErrNotAuthenticated
	}

	err := o.repo.ProcessWithdraw(ctx, model.Withdraw{OrderId: orderId, Sum: sum, UserId: user.Id})
	if err != nil {
		o.Log(ctx).Warn().Err(err).Msg("ProcessWithdraw:")
		return err
	}

	return nil
}

func (o withdrawService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "withdrawService").Logger()

	return &logger
}
