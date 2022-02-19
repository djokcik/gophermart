package service

import (
	"context"
	"errors"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/reporegistry"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/encrypt"
	"github.com/djokcik/gophermart/pkg/logging"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name=UserService

type UserService interface {
	Authenticate(ctx context.Context, login string, password string) (string, error)
	CreateUser(ctx context.Context, login string, password string) error
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	GenerateToken(ctx context.Context, user model.User) (string, error)
	GetBalance(ctx context.Context, user model.User) (model.UserBalance, error)
}

func NewUserService(cfg config.Config, registry reporegistry.RepoRegistry) UserService {
	return &userService{cfg: cfg, repo: registry.GetUserRepo(), withdrawRepo: registry.GetWithdrawRepo()}
}

type userService struct {
	cfg          config.Config
	repo         storage.UserRepository
	withdrawRepo storage.WithdrawRepository
}

func (u userService) GetBalance(ctx context.Context, user model.User) (model.UserBalance, error) {
	withdrawAmount, err := u.withdrawRepo.AmountWithdrawByUser(ctx, user.ID)
	if err != nil {
		u.Log(ctx).Error().Err(err).Msg("GetBalance:")
		return model.UserBalance{}, nil
	}

	return model.UserBalance{Current: user.Balance, Withdrawn: withdrawAmount}, nil
}

func (u userService) Authenticate(ctx context.Context, login string, password string) (string, error) {
	user, err := u.GetUserByUsername(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			u.Log(ctx).Trace().Err(err).Msg("authenticate: wrong username")
			return "", ErrWrongPassword
		}

		return "", err
	}

	if err := encrypt.CompareHashAndPassword(password+u.cfg.PasswordPepper, user.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			u.Log(ctx).Trace().Err(err).Msg("authenticate: wrong password")
			return "", ErrWrongPassword
		}

		return "", err
	}

	token, err := u.GenerateToken(ctx, user)
	if err != nil {
		return "", err
	}

	return token, err
}

func (u userService) CreateUser(ctx context.Context, login string, password string) error {
	user := model.User{Username: login, Password: password}
	err := user.Validate()
	if err != nil {
		u.Log(ctx).Trace().Err(err).Msgf("invalid validate user")
		return err
	}

	user.Password, err = encrypt.HashAndSalt(user.Password, u.cfg.PasswordPepper)
	if err != nil {
		u.Log(ctx).Trace().Err(err).Msgf("error create hash")
		return err
	}

	err = u.repo.CreateUser(ctx, user)
	if err != nil {
		u.Log(ctx).Trace().Err(err).Msg("invalid create user")
		return err
	}

	u.Log(ctx).Info().
		Str("Username", user.Username).
		Msg("success created user")

	return nil
}

func (u userService) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user, err := u.repo.UserByUsername(ctx, username)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u userService) GenerateToken(ctx context.Context, user model.User) (string, error) {
	token, err := encrypt.CreateToken(u.cfg.Key, user.ID)
	if err != nil {
		u.Log(ctx).Err(err).Msgf("error create token")
		return "", err
	}

	return token, nil
}

func (u userService) Log(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, "user service").Logger()

	return &logger
}
