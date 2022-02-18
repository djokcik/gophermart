package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/rs/zerolog"
)

func NewOrderRepository(db *sql.DB) storage.OrderRepository {
	return &orderRepository{db: db}
}

type orderRepository struct {
	db *sql.DB
}

func (r orderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO orders (id, user_id, status, accrual) VALUES ($1, $2, $3, $4)",
		order.Id,
		order.UserId,
		order.Status,
		order.Accrual,
	)

	if err != nil {
		r.Log(ctx).Err(err).Msg("invalid save order")
		return err
	}

	return err
}

func (r orderRepository) FindOrderById(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT user_id, status, uploaded_at, accrual from orders where id=$1",
		orderId,
	)

	order := model.Order{Id: orderId}
	err := row.Scan(&order.UserId, &order.Status, &order.UploadedAt, &order.Accrual)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, storage.ErrNotFound
		}

		return model.Order{}, err
	}

	return order, nil
}

func (r orderRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "database orderRepository").Logger()

	return &logger
}
