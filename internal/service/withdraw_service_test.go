package service

import (
	"context"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage/mocks"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_withdrawService_AmountWithdrawByUser(t *testing.T) {
	t.Run("should return amount withdraw by user", func(t *testing.T) {
		m := mocks.WithdrawRepository{Mock: mock.Mock{}}
		m.On("AmountWithdrawByUser", mock.Anything, 666).Return(model.Amount(1000), nil)

		service := withdrawService{repo: &m}

		amount, err := service.AmountWithdrawByUser(context.Background(), 666)

		m.AssertNumberOfCalls(t, "AmountWithdrawByUser", 1)
		require.Equal(t, err, nil)
		require.Equal(t, amount, model.Amount(1000))
	})
}

func Test_withdrawService_WithdrawLogsByUserID(t *testing.T) {
	t.Run("should return withdraw logs by userID", func(t *testing.T) {
		logs := []model.Withdraw{
			{ID: 1, UserID: 666, OrderID: "123", Sum: 1000},
			{ID: 2, UserID: 666, OrderID: "333", Sum: 2123},
		}

		m := mocks.WithdrawRepository{Mock: mock.Mock{}}
		m.On("WithdrawLogsByUserID", mock.Anything, 666).Return(logs, nil)

		service := withdrawService{repo: &m}

		amount, err := service.WithdrawLogsByUserID(context.Background(), 666)

		m.AssertNumberOfCalls(t, "WithdrawLogsByUserID", 1)
		require.Equal(t, err, nil)
		require.Equal(t, amount, logs)
	})
}

func Test_withdrawService_ProcessWithdraw(t *testing.T) {
	t.Run("should call repo process withdraw", func(t *testing.T) {
		m := mocks.WithdrawRepository{Mock: mock.Mock{}}
		m.On("ProcessWithdraw", mock.Anything, model.Withdraw{OrderID: "1", Sum: 1000, UserID: 666}).
			Return(nil)

		service := withdrawService{repo: &m}

		ctx := appContext.WithUser(context.Background(), &model.User{ID: 666})

		err := service.ProcessWithdraw(ctx, "1", 1000)

		m.AssertNumberOfCalls(t, "ProcessWithdraw", 1)
		require.Equal(t, err, nil)
	})
}
