package psql

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_userRepository_CreateUser(t *testing.T) {
	t.Run("1. Should create user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &userRepository{db: db}

		mock.
			ExpectExec("INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\)").
			WithArgs("test", "userPassword").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.CreateUser(
			context.Background(),
			model.User{Username: "test", Password: "userPassword"},
		)

		require.Equal(t, err, nil)
	})

	t.Run("2. Should return error as duplicate username", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &userRepository{db: db}

		mock.
			ExpectExec("INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\)").
			WithArgs("test", "userPassword").
			WillReturnError(pgx.PgError{Code: pgerrcode.UniqueViolation})

		err = repo.CreateUser(
			context.Background(),
			model.User{Username: "test", Password: "userPassword"},
		)

		require.Equal(t, err, storage.ErrLoginAlreadyExists)
	})
}

func Test_userRepository_UserByUsername(t *testing.T) {
	t.Run("should return user by username", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &userRepository{db: db}
		now := time.Now()

		row := sqlmock.NewRows([]string{"id", "password", "created_at", "balance"}).
			AddRow(666, "testPassword", now, 1000)
		mock.ExpectQuery("SELECT id, password, created_at, balance from users where username=\\$1").
			WithArgs("testUsername").
			WillReturnRows(row)

		user, err := repo.UserByUsername(
			context.Background(),
			"testUsername",
		)

		require.Equal(t, err, nil)
		require.Equal(t, user, model.User{
			ID:        666,
			Username:  "testUsername",
			Password:  "testPassword",
			CreatedAt: now,
			Balance:   1000,
		},
		)
	})
}

func Test_userRepository_UserByID(t *testing.T) {
	t.Run("should return user by id", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		repo := &userRepository{db: db}
		now := time.Now()

		row := sqlmock.NewRows([]string{"username", "password", "created_at", "balance"}).
			AddRow("testUsername", "testPassword", now, 1000)
		mock.ExpectQuery("SELECT username, password, created_at, balance from users where id=\\$1").
			WithArgs(666).
			WillReturnRows(row)

		user, err := repo.UserByID(
			context.Background(),
			666,
		)

		require.Equal(t, err, nil)
		require.Equal(t, user, model.User{
			ID:        666,
			Username:  "testUsername",
			Password:  "testPassword",
			CreatedAt: now,
			Balance:   1000,
		},
		)
	})
}
