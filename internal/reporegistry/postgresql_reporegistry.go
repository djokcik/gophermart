package reporegistry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/internal/storage/psql"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/stdlib"
)

//go:generate mockery --name=RepoRegistry

// RepoRegistry .
type RepoRegistry interface {
	GetUserRepo() storage.UserRepository
	GetOrderRepo() storage.OrderRepository
	GetWithdrawRepo() storage.WithdrawRepository
}

type postgresqlRepoRegistry struct {
	db *sql.DB
}

func NewPostgreSQL(ctx context.Context, cfg config.Config) (RepoRegistry, error) {
	err := autoMigrate(ctx, "file://internal/storage/psql/migrations", cfg)
	if err != nil {
		return nil, err
	}

	db, err := open(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &postgresqlRepoRegistry{db: db}, nil
}

func autoMigrate(ctx context.Context, path string, cfg config.Config) error {
	_, logger := logging.GetCtxLogger(ctx)

	m, err := migrate.New(path, cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("psql: autoMigrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("psql: autoMigrate: %w", err)
	}

	logger.Info().Msg("success auto migrate")
	return nil
}

func open(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	_, logger := logging.GetCtxLogger(ctx)

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Unable to connect to database")
		return nil, err
	}

	return db, nil
}

func (r postgresqlRepoRegistry) GetUserRepo() storage.UserRepository {
	return psql.NewUserRepository(r.db)
}

func (r postgresqlRepoRegistry) GetOrderRepo() storage.OrderRepository {
	return psql.NewOrderRepository(r.db)
}

func (r postgresqlRepoRegistry) GetWithdrawRepo() storage.WithdrawRepository {
	return psql.NewWithdrawRepository(r.db)
}
