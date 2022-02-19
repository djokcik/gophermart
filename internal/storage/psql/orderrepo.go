package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/djokcik/gophermart/provider"
	"github.com/rs/zerolog"
)

func NewOrderRepository(db *sql.DB) storage.OrderRepository {
	return &orderRepository{db: db}
}

type orderRepository struct {
	db *sql.DB
}

func (r orderRepository) OrdersByStatus(ctx context.Context, status model.Status) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, uploaded_at, accrual 
		from orders WHERE status = $1`, status)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		order := model.Order{Status: status}
		err = rows.Scan(&order.Id, &order.UserId, &order.UploadedAt, &order.Accrual)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r orderRepository) UpdateForAccrual(ctx context.Context, order model.Order, accrual provider.AccrualResponse) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.Log(ctx).Error().Err(err).Msg("UpdateForAccrual: prepare transaction")
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE orders SET status = $1, accrual = $2 
			WHERE id = $3`, accrual.Status, accrual.Accrual, order.Id)
	if err != nil {
		r.Log(ctx).Error().Err(err).Msg("UpdateForAccrual: exec orders")
		if err = tx.Rollback(); err != nil {
			r.Log(ctx).Error().Err(err).Msgf("UpdateForAccrual: unable to rollback")
			return err
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, accrual.Accrual, order.UserId)
	if err != nil {
		r.Log(ctx).Error().Err(err).Msg("UpdateForAccrual: exec users")
		if err = tx.Rollback(); err != nil {
			r.Log(ctx).Error().Err(err).Msgf("UpdateForAccrual: unable to rollback")
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		r.Log(ctx).Error().Err(err).Msgf("UpdateForAccrual: unable to commit")
		return err
	}

	return nil
}

func (r orderRepository) OrdersByUserId(ctx context.Context, userId int) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, status, uploaded_at, accrual 
		from orders WHERE user_id = $1 ORDER BY uploaded_at`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.Order, 0)
	for rows.Next() {
		order := model.Order{UserId: userId}
		err = rows.Scan(&order.Id, &order.Status, &order.UploadedAt, &order.Accrual)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
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

func (r orderRepository) OrderById(ctx context.Context, orderId model.OrderId) (model.Order, error) {
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
