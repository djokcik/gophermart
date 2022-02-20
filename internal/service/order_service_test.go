package service

import (
	"context"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/internal/storage/mocks"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/provider"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_orderService_UpdateForAccrual(t *testing.T) {
	t.Run("should update for accrual", func(t *testing.T) {
		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("UpdateForAccrual", mock.Anything, model.Order{ID: "1"}, provider.AccrualResponse{Order: "1"}).
			Return(nil)

		service := orderService{repo: &m}

		err := service.UpdateForAccrual(context.Background(), model.Order{ID: "1"}, provider.AccrualResponse{Order: "1"})

		m.AssertNumberOfCalls(t, "UpdateForAccrual", 1)
		require.Equal(t, err, nil)
	})
}

func Test_orderService_OrdersByStatus(t *testing.T) {
	t.Run("should return orders by status", func(t *testing.T) {
		orders := []model.Order{
			{ID: "1", UserID: 666, Accrual: 1000},
			{ID: "2", UserID: 666, Accrual: 1234},
		}

		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("OrdersByStatus", mock.Anything, model.StatusNew).Return(orders, nil)

		service := orderService{repo: &m}

		results, err := service.OrdersByStatus(context.Background(), model.StatusNew)

		m.AssertNumberOfCalls(t, "OrdersByStatus", 1)
		require.Equal(t, err, nil)
		require.Equal(t, results, orders)
	})
}

func Test_orderService_OrdersByUser(t *testing.T) {
	t.Run("should return orders by user", func(t *testing.T) {
		orders := []model.Order{
			{ID: "1", UserID: 666, Accrual: 1000},
			{ID: "2", UserID: 666, Accrual: 1234},
		}

		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("OrdersByUserID", mock.Anything, 666).Return(orders, nil)

		service := orderService{repo: &m}

		results, err := service.OrdersByUser(context.Background(), 666)

		m.AssertNumberOfCalls(t, "OrdersByUserID", 1)
		require.Equal(t, err, nil)
		require.Equal(t, results, orders)
	})
}

func Test_orderService_ProcessOrder(t *testing.T) {
	t.Run("1. should call createOrder", func(t *testing.T) {
		o := model.Order{ID: "1", UserID: 666, Status: model.StatusNew}

		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("OrderByID", mock.Anything, model.OrderID("1")).Return(model.Order{}, storage.ErrNotFound)
		m.On("CreateOrder", mock.Anything, o).Return(nil)

		service := orderService{repo: &m}

		ctx := appContext.WithUser(context.Background(), &model.User{ID: 666})

		err := service.ProcessOrder(ctx, "1")

		m.AssertNumberOfCalls(t, "OrderByID", 1)
		m.AssertNumberOfCalls(t, "CreateOrder", 1)
		require.Equal(t, err, nil)
	})

	t.Run("2. should return error that order already uploaded", func(t *testing.T) {
		o := model.Order{ID: "1", UserID: 666, Status: model.StatusNew}

		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("OrderByID", mock.Anything, model.OrderID("1")).Return(o, nil)

		service := orderService{repo: &m}

		ctx := appContext.WithUser(context.Background(), &model.User{ID: 666})

		err := service.ProcessOrder(ctx, "1")

		m.AssertNumberOfCalls(t, "OrderByID", 1)
		require.Equal(t, err, ErrOrderAlreadyUploaded)
	})

	t.Run("3. should return error that order already uploaded for another user", func(t *testing.T) {
		o := model.Order{ID: "1", UserID: 111, Status: model.StatusNew}

		m := mocks.OrderRepository{Mock: mock.Mock{}}
		m.On("OrderByID", mock.Anything, model.OrderID("1")).Return(o, nil)

		service := orderService{repo: &m}

		ctx := appContext.WithUser(context.Background(), &model.User{ID: 666})

		err := service.ProcessOrder(ctx, "1")

		m.AssertNumberOfCalls(t, "OrderByID", 1)
		require.Equal(t, err, ErrOrderAlreadyUploadedAnotherUser)
	})
}
