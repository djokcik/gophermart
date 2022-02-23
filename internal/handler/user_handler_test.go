package handler

import (
	"bytes"
	"context"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service/mocks"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetBalanceHandler(t *testing.T) {
	t.Run("should return user balance", func(t *testing.T) {
		m := mocks.UserService{Mock: mock.Mock{}}
		m.On("GetBalance", mock.Anything, model.User{ID: 666}).
			Return(model.UserBalance{Current: 1012, Withdrawn: 5555}, nil)

		request := httptest.NewRequest(http.MethodGet, "/user/balance", nil)
		request = request.WithContext(appContext.WithUser(context.Background(), &model.User{ID: 666}))

		h := Handler{user: &m, Mux: chi.NewMux()}
		h.Get("/user/balance", h.GetBalanceHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "GetBalance", 1)
		require.Equal(t, string(resBody), `{"current":10.12,"withdrawn":55.55}`)
	})
}

func TestHandler_RegisterUserHandler(t *testing.T) {
	t.Run("should user be registered", func(t *testing.T) {
		m := mocks.UserService{Mock: mock.Mock{}}
		m.On("CreateUser", mock.Anything, "userLogin", "userPassword").Return(nil)
		m.On("GetUserByUsername", mock.Anything, "userLogin").
			Return(model.User{ID: 666}, nil)
		m.On("GenerateToken", mock.Anything, model.User{ID: 666}).Return("secretToken", nil)

		body := bytes.NewReader([]byte(`{"login":"userLogin","password":"userPassword"}`))

		request := httptest.NewRequest(http.MethodPost, "/user/register", body)

		h := Handler{user: &m, Mux: chi.NewMux()}
		h.Post("/user/register", h.RegisterUserHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "CreateUser", 1)
		m.AssertNumberOfCalls(t, "GetUserByUsername", 1)
		m.AssertNumberOfCalls(t, "GenerateToken", 1)
		require.Equal(t, string(resBody), `{"token":"secretToken"}`)
		require.Equal(t, string(resBody), `{"token":"secretToken"}`)
		require.Equal(t, res.Header.Get("Authorization"), "Bearer: secretToken")
		require.Equal(t, res.Header.Get("Content-Type"), "application/json")
		require.Equal(t, res.Cookies()[0].Value, "secretToken")
	})
}

func TestHandler_SignInHandler(t *testing.T) {
	t.Run("should user be authorized", func(t *testing.T) {
		m := mocks.UserService{Mock: mock.Mock{}}
		m.On("Authenticate", mock.Anything, "userLogin", "userPassword").
			Return("secretToken", nil)

		body := bytes.NewReader([]byte(`{"login":"userLogin","password":"userPassword"}`))
		request := httptest.NewRequest(http.MethodPost, "/user/login", body)

		h := Handler{user: &m, Mux: chi.NewMux()}
		h.Post("/user/login", h.SignInHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "Authenticate", 1)
		require.Equal(t, string(resBody), `{"token":"secretToken"}`)
		require.Equal(t, res.Header.Get("Authorization"), "Bearer: secretToken")
		require.Equal(t, res.Header.Get("Content-Type"), "application/json")
		require.Equal(t, res.Cookies()[0].Value, "secretToken")
	})
}
