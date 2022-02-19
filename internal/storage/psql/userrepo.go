package psql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
)

func NewUserRepository(db *sql.DB) storage.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *sql.DB
}

func (r userRepository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		if err, ok := err.(pgx.PgError); ok && err.Code == pgerrcode.UniqueViolation /* or just == "23505" */ {
			return storage.ErrLoginAlreadyExists
		}

		return err
	}

	return err
}

func (r userRepository) UserByUsername(ctx context.Context, username string) (model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, password, created_at, balance from users where username=$1", username)

	user := model.User{Username: username}
	err := row.Scan(&user.Id, &user.Password, &user.CreatedAt, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, storage.ErrNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r userRepository) UserById(ctx context.Context, id int) (model.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT username, password, created_at, balance from users where id=$1", id)

	user := model.User{Id: id}
	err := row.Scan(&user.Username, &user.Password, &user.CreatedAt, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, storage.ErrNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r userRepository) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "database userRepository").Logger()

	return &logger
}
