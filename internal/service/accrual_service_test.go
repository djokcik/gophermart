package service

import (
	"context"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service/mocks"
	"github.com/djokcik/gophermart/provider"
	providerMocks "github.com/djokcik/gophermart/provider/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_accrualService_ProcessOrder(t *testing.T) {
	t.Run("should update order for accrual", func(t *testing.T) {
		order := model.Order{ID: "1", Status: model.StatusNew}
		accrualResponse := provider.AccrualResponse{Order: order.ID, Accrual: 1000}

		mockClient := providerMocks.AccrualClient{Mock: mock.Mock{}}
		mockClient.On("GetOrder", mock.Anything, order.ID).
			Return(accrualResponse, nil)

		mockOrder := mocks.OrderService{Mock: mock.Mock{}}
		mockOrder.On("UpdateForAccrual", mock.Anything, order, accrualResponse).
			Return(nil)

		service := accrualService{order: &mockOrder, client: &mockClient}

		service.ProcessOrder(context.Background(), order)

		mockClient.AssertNumberOfCalls(t, "GetOrder", 1)
		mockOrder.AssertNumberOfCalls(t, "UpdateForAccrual", 1)
	})
}

func Test_accrualService_getOrders(t *testing.T) {
	t.Run("should return orders which need update", func(t *testing.T) {
		newOrders := []model.Order{{ID: "5", UserID: 666, Status: model.StatusNew}}

		processingOrders := []model.Order{
			{ID: "1", Status: model.StatusProcessing, UserID: 666},
			{ID: "2", Status: model.StatusProcessing, UserID: 111},
		}

		m := mocks.OrderService{Mock: mock.Mock{}}
		m.On("OrdersByStatus", mock.Anything, model.StatusNew).Return(newOrders, nil)
		m.On("OrdersByStatus", mock.Anything, model.StatusProcessing).Return(processingOrders, nil)

		service := accrualService{order: &m}

		orders := service.getOrders(context.Background())

		m.AssertNumberOfCalls(t, "OrdersByStatus", 2)
		require.Equal(t, orders, []model.Order{
			{ID: "5", UserID: 666, Status: model.StatusNew},
			{ID: "1", Status: model.StatusProcessing, UserID: 666},
			{ID: "2", Status: model.StatusProcessing, UserID: 111},
		})
	})
}
