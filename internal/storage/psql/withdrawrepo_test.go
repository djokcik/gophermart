package psql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_withdrawRepository_AmountWithdrawByUser(t *testing.T) {
	t.Run("should return amount withdraw by user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &withdrawRepository{db: db}

		row := sqlmock.NewRows([]string{"amount"}).
			AddRow(1000)
		mock.ExpectQuery("SELECT SUM\\(sum\\) as amount from withdraw_log GROUP BY user_id = \\$1").
			WithArgs(666).
			WillReturnRows(row)

		amount, err := repo.AmountWithdrawByUser(
			context.Background(),
			666,
		)

		require.Equal(t, err, nil)
		require.Equal(t, amount, model.Amount(1000))
	})
}

func Test_withdrawRepository_WithdrawLogsByUserID(t *testing.T) {
	t.Run("should return withdraw logs by user id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &withdrawRepository{db: db}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "sum", "processed_at", "order_id"}).
			AddRow(1, 1000, now, "123").
			AddRow(2, 4312, now, "555")
		mock.ExpectQuery("SELECT id, sum, processed_at, order_id from withdraw_log WHERE user_id = \\$1 ORDER BY processed_at").
			WithArgs(666).
			WillReturnRows(rows)

		result, err := repo.WithdrawLogsByUserID(context.Background(), 666)

		require.Equal(t, err, nil)
		require.Equal(t, result, []model.Withdraw{
			{ID: 1, OrderID: "123", Sum: 1000, ProcessedAt: model.UploadedTime(now), UserID: 666},
			{ID: 2, OrderID: "555", Sum: 4312, ProcessedAt: model.UploadedTime(now), UserID: 666},
		})
	})
}
