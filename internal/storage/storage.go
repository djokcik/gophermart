package storage

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
)

//go:generate mockery --name=OrderRepository

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	UserByUsername(ctx context.Context, username string) (model.User, error)
	UserById(ctx context.Context, id int) (model.User, error)
}

type OrderRepository interface {
	FindOrderById(ctx context.Context, id model.OrderId) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	FindOrdersByUserId(ctx context.Context, userId int) ([]model.Order, error)
}

var (
	ErrNotFound           = errors.New("storage: not found")
	ErrLoginAlreadyExists = errors.New("storage: login already exists")
	//
	//ErrAlreadyProcessed = errors.New("storage: already processed")
	//ErrInvalidStatus    = errors.New("storage: non-processed order has PROCESSED status")
	//
	//ErrInsufficientPoints = errors.New("storage: insufficient points to perform withdrawal operation")
	//
	//// ErrInvalidInput is threw when accrual or withdrawal amount less than zero.
	//ErrInvalidInput = errors.New("storage: amount less than zero")
)
