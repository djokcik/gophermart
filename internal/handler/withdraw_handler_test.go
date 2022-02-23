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

func TestHandler_WithdrawLogsHandler(t *testing.T) {
	t.Run("should return withdraw logs", func(t *testing.T) {
		m := mocks.WithdrawService{Mock: mock.Mock{}}
		m.On("WithdrawLogsByUserID", mock.Anything, 666).Return([]model.Withdraw{
			{ID: 1, OrderID: "111", Sum: 1000},
			{ID: 2, OrderID: "222", Sum: 1234},
		}, nil)

		request := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
		request = request.WithContext(appContext.WithUser(context.Background(), &model.User{ID: 666}))

		h := Handler{withdraw: &m, Mux: chi.NewMux()}
		h.Get("/withdrawals", h.WithdrawLogsHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "WithdrawLogsByUserID", 1)
		require.Equal(t, string(resBody),
			`[{"order":"111","sum":10,"processed_at":"0001-01-01T00:00:00Z"},{"order":"222","sum":12.34,"processed_at":"0001-01-01T00:00:00Z"}]`,
		)
	})
}

func TestHandler_WithdrawHandler(t *testing.T) {
	t.Run("should be correct withdraw", func(t *testing.T) {
		m := mocks.WithdrawService{Mock: mock.Mock{}}
		m.On("ProcessWithdraw", mock.Anything, model.OrderID("9278923470"), model.Amount(1012)).
			Return(nil)

		body := bytes.NewReader([]byte(`{"order":"9278923470","sum":10.12}`))

		request := httptest.NewRequest(http.MethodPost, "/balance/withdraw", body)

		h := Handler{withdraw: &m, Mux: chi.NewMux()}
		h.Post("/balance/withdraw", h.WithdrawHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		require.Equal(t, string(resBody), "OK")
		m.AssertNumberOfCalls(t, "ProcessWithdraw", 1)
	})
}
