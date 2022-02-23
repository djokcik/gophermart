package handler

import (
	"bytes"
	"context"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service"
	"github.com/djokcik/gophermart/internal/service/mocks"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_UploadOrderHandler(t *testing.T) {
	t.Run("1. should order uploaded with status `202`", func(t *testing.T) {
		m := mocks.OrderService{Mock: mock.Mock{}}
		m.On("ProcessOrder", mock.Anything, model.OrderID("9278923470")).Return(nil)

		body := bytes.NewReader([]byte(`9278923470`))

		request := httptest.NewRequest(http.MethodPost, "/user/orders", body)

		h := Handler{order: &m, Mux: chi.NewMux()}
		h.Post("/user/orders", h.UploadOrderHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "ProcessOrder", 1)
		require.Equal(t, string(resBody), "OK")
		require.Equal(t, res.StatusCode, http.StatusAccepted)
	})

	t.Run("2. should return error when orderID is invalid", func(t *testing.T) {
		body := bytes.NewReader([]byte(`1`))

		request := httptest.NewRequest(http.MethodPost, "/user/orders", body)

		h := Handler{Mux: chi.NewMux()}
		h.Post("/user/orders", h.UploadOrderHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		require.Equal(t, string(resBody), "invalid orderID\n")
		require.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
	})

	t.Run("3. should return error when order upload another user", func(t *testing.T) {
		m := mocks.OrderService{Mock: mock.Mock{}}
		m.On("ProcessOrder", mock.Anything, model.OrderID("9278923470")).
			Return(service.ErrOrderAlreadyUploadedAnotherUser)

		body := bytes.NewReader([]byte(`9278923470`))

		request := httptest.NewRequest(http.MethodPost, "/user/orders", body)

		h := Handler{order: &m, Mux: chi.NewMux()}
		h.Post("/user/orders", h.UploadOrderHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "ProcessOrder", 1)
		require.Equal(t, string(resBody), "order already uploaded\n")
		require.Equal(t, res.StatusCode, http.StatusConflict)
	})

	t.Run("4. shouldn`t upload order when order already uploaded", func(t *testing.T) {
		m := mocks.OrderService{Mock: mock.Mock{}}
		m.On("ProcessOrder", mock.Anything, model.OrderID("9278923470")).
			Return(service.ErrOrderAlreadyUploaded)

		body := bytes.NewReader([]byte(`9278923470`))

		request := httptest.NewRequest(http.MethodPost, "/user/orders", body)

		h := Handler{order: &m, Mux: chi.NewMux()}
		h.Post("/user/orders", h.UploadOrderHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "ProcessOrder", 1)
		require.Equal(t, string(resBody), "OK")
		require.Equal(t, res.StatusCode, http.StatusOK)
	})
}

func TestHandler_GetOrdersHandler(t *testing.T) {
	t.Run("1. should return order list", func(t *testing.T) {
		uploaded, _ := time.Parse(time.RFC3339, "2020-12-10T15:15:45+03:00")

		orders := []model.Order{
			{ID: "1", UserID: 666, Status: model.StatusProcessed, Accrual: 50000, UploadedAt: model.UploadedTime(uploaded)},
			{ID: "2", UserID: 666, Status: model.StatusProcessing, Accrual: 10012, UploadedAt: model.UploadedTime(uploaded)},
		}

		m := mocks.OrderService{Mock: mock.Mock{}}
		m.On("OrdersByUser", mock.Anything, 666).Return(orders, nil)

		request := httptest.NewRequest(http.MethodGet, "/user/orders", nil)
		request = request.WithContext(appContext.WithUser(context.Background(), &model.User{ID: 666}))

		h := Handler{order: &m, Mux: chi.NewMux()}
		h.Get("/user/orders", h.GetOrdersHandler())

		w := httptest.NewRecorder()

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		resBody, _ := io.ReadAll(res.Body)

		m.AssertNumberOfCalls(t, "OrdersByUser", 1)
		require.Equal(t, string(resBody), `[{"number":"1","status":"PROCESSED","uploaded_at":"2020-12-10T15:15:45+03:00","accrual":500},{"number":"2","status":"PROCESSING","uploaded_at":"2020-12-10T15:15:45+03:00","accrual":100.12}]`)
	})
}
