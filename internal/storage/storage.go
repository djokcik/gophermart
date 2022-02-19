package storage

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/provider"
)

//go:generate mockery --name=UserRepository
//go:generate mockery --name=OrderRepository
//go:generate mockery --name=WithdrawRepository

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	UserByUsername(ctx context.Context, username string) (model.User, error)
	UserById(ctx context.Context, id int) (model.User, error)
}

type OrderRepository interface {
	OrderById(ctx context.Context, id model.OrderId) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error)
	OrdersByUserId(ctx context.Context, userId int) ([]model.Order, error)
	UpdateForAccrual(ctx context.Context, order model.Order, accrual provider.AccrualResponse) error
}

type WithdrawRepository interface {
	ProcessWithdraw(ctx context.Context, withdraw model.Withdraw) error
	WithdrawLogsByUserId(ctx context.Context, userId int) ([]model.Withdraw, error)
	AmountWithdrawByUser(ctx context.Context, userId int) (model.Amount, error)
}

var (
	ErrNotFound           = errors.New("storage: not found")
	ErrLoginAlreadyExists = errors.New("storage: login already exists")
	ErrInsufficientFunds  = errors.New("storage: insufficient funds")
)
