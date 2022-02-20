package psql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_orderRepository_OrdersByStatus(t *testing.T) {
	t.Run("should return orders by status", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &orderRepository{db: db}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "user_id", "uploaded_at", "accrual"}).
			AddRow("1", 666, now, 1000).
			AddRow("2", 666, now, 1555)
		mock.ExpectQuery("SELECT id, user_id, uploaded_at, accrual from orders WHERE status = \\$1").
			WithArgs(model.StatusNew).
			WillReturnRows(rows)

		result, err := repo.OrdersByStatus(context.Background(), model.StatusNew)

		require.Equal(t, err, nil)
		require.Equal(t, result, []model.Order{
			{ID: "1", Status: model.StatusNew, Accrual: 1000, UploadedAt: model.UploadedTime(now), UserID: 666},
			{ID: "2", Status: model.StatusNew, Accrual: 1555, UploadedAt: model.UploadedTime(now), UserID: 666},
		})
	})
}

func Test_orderRepository_OrdersByUserID(t *testing.T) {
	t.Run("should return orders by userID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &orderRepository{db: db}

		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "status", "uploaded_at", "accrual"}).
			AddRow("1", model.StatusNew, now, 1000).
			AddRow("2", model.StatusProcessed, now, 1555)
		mock.ExpectQuery("SELECT id, status, uploaded_at, accrual from orders WHERE user_id = \\$1 ORDER BY uploaded_at").
			WithArgs(666).
			WillReturnRows(rows)

		result, err := repo.OrdersByUserID(context.Background(), 666)

		require.Equal(t, err, nil)
		require.Equal(t, result, []model.Order{
			{ID: "1", Status: model.StatusNew, Accrual: 1000, UploadedAt: model.UploadedTime(now), UserID: 666},
			{ID: "2", Status: model.StatusProcessed, Accrual: 1555, UploadedAt: model.UploadedTime(now), UserID: 666},
		})
	})
}

func Test_orderRepository_CreateOrder(t *testing.T) {
	t.Run("should create order", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &orderRepository{db: db}

		mock.
			ExpectExec("INSERT INTO orders \\(id, user_id, status, accrual\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
			WithArgs("1", 666, model.StatusNew, 1000).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateOrder(
			context.Background(),
			model.Order{ID: "1", UserID: 666, Status: model.StatusNew, Accrual: 1000},
		)

		require.Equal(t, err, nil)
	})
}

func Test_orderRepository_OrderByID(t *testing.T) {
	t.Run("should return order by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &orderRepository{db: db}
		now := time.Now()

		row := sqlmock.NewRows([]string{"user_id", "status", "uploaded_at", "accrual"}).
			AddRow(666, model.StatusNew, now, 1000)
		mock.ExpectQuery("SELECT user_id, status, uploaded_at, accrual from orders where id=\\$1").
			WithArgs("1").
			WillReturnRows(row)

		order, err := repo.OrderByID(
			context.Background(),
			"1",
		)

		require.Equal(t, err, nil)
		require.Equal(t, order, model.Order{
			ID:         "1",
			UserID:     666,
			Status:     model.StatusNew,
			UploadedAt: model.UploadedTime(now),
			Accrual:    1000,
		},
		)
	})
}
