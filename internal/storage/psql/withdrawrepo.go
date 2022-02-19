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

func NewWithdrawRepository(db *sql.DB) storage.WithdrawRepository {
	return &withdrawRepository{db: db}
}

type withdrawRepository struct {
	db *sql.DB
}

func (r withdrawRepository) AmountWithdrawByUser(ctx context.Context, userId int) (model.Amount, error) {
	row := r.db.QueryRow("SELECT SUM(sum) as amount from withdraw_log GROUP BY user_id")

	var amount model.Amount
	err := row.Scan(&amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		r.Log(ctx).Error().Err(err).Msg("AmountWithdrawByUser: invalid scan")
		return 0, err
	}

	return amount, nil
}

func (r withdrawRepository) ProcessWithdraw(ctx context.Context, withdraw model.Withdraw) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.Log(ctx).Error().Err(err).Msg("ProcessWithdraw: prepare transaction")
		return err
	}

	row := tx.QueryRow("SELECT balance FROM users where id = $1", withdraw.UserId)
	var balance model.Amount
	if err = row.Scan(&balance); err != nil {
		return err
	}

	if balance < withdraw.Sum {
		if err = tx.Rollback(); err != nil {
			r.Log(ctx).Error().Err(err).Msgf("ProcessWithdraw: unable to rollback")
			return err
		}

		return storage.ErrInsufficientFunds
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO withdraw_log (user_id, sum, order_id) VALUES ($1, $2, $3)`,
		withdraw.UserId, withdraw.Sum, withdraw.OrderId); err != nil {
		r.Log(ctx).Error().Err(err).Msg("ProcessWithdraw: exec withdraw_log")
		if err = tx.Rollback(); err != nil {
			r.Log(ctx).Error().Err(err).Msgf("ProcessWithdraw: unable to rollback")
			return err
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = balance - $1 WHERE id = $2`, withdraw.Sum, withdraw.UserId)
	if err != nil {
		r.Log(ctx).Error().Err(err).Msg("ProcessWithdraw: exec users")
		if err = tx.Rollback(); err != nil {
			r.Log(ctx).Error().Err(err).Msgf("ProcessWithdraw: unable to rollback")
			return err
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		r.Log(ctx).Error().Err(err).Msgf("ProcessWithdraw: unable to commit")
		return err
	}

	return nil
}

func (r withdrawRepository) WithdrawLogsByUserId(ctx context.Context, userId int) ([]model.Withdraw, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, sum, processed_at, order_id 
		from withdraw_log WHERE user_id = $1 ORDER BY processed_at`, userId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	withdrawLogs := make([]model.Withdraw, 0)
	for rows.Next() {
		withdrawLog := model.Withdraw{UserId: userId}
		err = rows.Scan(&withdrawLog.Id, &withdrawLog.Sum, &withdrawLog.ProcessedAt, &withdrawLog.OrderId)
		if err != nil {
			return nil, err
		}

		withdrawLogs = append(withdrawLogs, withdrawLog)
	}

	return withdrawLogs, nil
}

func (r withdrawRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "database withdrawRepository").Logger()

	return &logger
}
