package service

import (
	"context"
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/model"
	serviceMock "github.com/djokcik/gophermart/internal/service/mocks"
	"github.com/djokcik/gophermart/internal/storage/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_userService_GetBalance(t *testing.T) {
	t.Run("should return user balance", func(t *testing.T) {
		m := mocks.WithdrawRepository{Mock: mock.Mock{}}
		m.On("AmountWithdrawByUser", mock.Anything, 666).Return(model.Amount(1000), nil)

		service := userService{withdrawRepo: &m}

		amount, err := service.GetBalance(context.Background(), model.User{ID: 666, Balance: 1234})

		m.AssertNumberOfCalls(t, "AmountWithdrawByUser", 1)
		require.Equal(t, err, nil)
		require.Equal(t, amount, model.UserBalance{Withdrawn: 1000, Current: 1234})
	})
}

func Test_userService_GetUserByUsername(t *testing.T) {
	t.Run("should return user by username", func(t *testing.T) {
		user := model.User{ID: 666, Username: "testUsername"}

		m := mocks.UserRepository{Mock: mock.Mock{}}
		m.On("UserByUsername", mock.Anything, "testUsername").Return(user, nil)

		service := userService{repo: &m}

		result, err := service.GetUserByUsername(context.Background(), "testUsername")

		m.AssertNumberOfCalls(t, "UserByUsername", 1)
		require.Equal(t, err, nil)
		require.Equal(t, result, user)
	})
}

func Test_userService_GenerateToken(t *testing.T) {
	t.Run("should be generated token", func(t *testing.T) {
		m := serviceMock.AuthService{Mock: mock.Mock{}}
		m.On("CreateToken", "key", 666).Return("secretToken", nil)

		service := userService{auth: &m, cfg: config.Config{Key: "key"}}

		token, err := service.GenerateToken(context.Background(), model.User{ID: 666})

		m.AssertNumberOfCalls(t, "CreateToken", 1)
		require.Equal(t, err, nil)
		require.Equal(t, token, "secretToken")
	})
}

func Test_userService_CreateUser(t *testing.T) {
	t.Run("should be created user", func(t *testing.T) {
		authMock := serviceMock.AuthService{Mock: mock.Mock{}}
		authMock.On("HashAndSalt", "userPassword", "pepper").
			Return("HashedPassword", nil)

		repoMock := mocks.UserRepository{Mock: mock.Mock{}}
		repoMock.On("CreateUser", mock.Anything, model.User{Username: "UserLogin", Password: "HashedPassword"}).
			Return(nil)

		service := userService{auth: &authMock, repo: &repoMock, cfg: config.Config{PasswordPepper: "pepper"}}

		err := service.CreateUser(context.Background(), "UserLogin", "userPassword")

		authMock.AssertNumberOfCalls(t, "HashAndSalt", 1)
		repoMock.AssertNumberOfCalls(t, "CreateUser", 1)
		require.Equal(t, err, nil)
	})
}

func Test_userService_Authenticate(t *testing.T) {
	t.Run("should authenticate user", func(t *testing.T) {
		repoMock := mocks.UserRepository{Mock: mock.Mock{}}
		repoMock.On("UserByUsername", mock.Anything, "UserLogin").
			Return(model.User{ID: 666, Password: "HashedPassword"}, nil)

		authMock := serviceMock.AuthService{Mock: mock.Mock{}}
		authMock.On("CompareHashAndPassword", "userPasswordpepper", "HashedPassword").
			Return(nil)
		authMock.On("CreateToken", "key", 666).Return("secretToken", nil)

		service := userService{
			auth: &authMock,
			repo: &repoMock,
			cfg:  config.Config{PasswordPepper: "pepper", Key: "key"},
		}

		token, err := service.Authenticate(context.Background(), "UserLogin", "userPassword")

		authMock.AssertNumberOfCalls(t, "CompareHashAndPassword", 1)
		authMock.AssertNumberOfCalls(t, "CreateToken", 1)
		repoMock.AssertNumberOfCalls(t, "UserByUsername", 1)
		require.Equal(t, err, nil)
		require.Equal(t, token, "secretToken")
	})
}
