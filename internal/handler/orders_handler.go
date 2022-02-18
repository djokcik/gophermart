package handler

import (
	"encoding/json"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/logging"
	"io"
	"net/http"
)

func (h *Handler) UploadOrderHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "UploadOrderHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Trace().Err(err).Msg("invalid parse order")
			http.Error(rw, "invalid parse order", http.StatusBadRequest)
			return
		}

		orderId := model.OrderId(body)
		if !orderId.Valid() {
			logger.Trace().Err(err).Msg("invalid orderId")
			http.Error(rw, "invalid orderId", http.StatusUnprocessableEntity)
			return
		}

		err = h.order.ProcessOrder(ctx, orderId)
		if err != nil {
			if errors.Is(err, service.ErrOrderAlreadyUploadedAnotherUser) {
				logger.Trace().Err(err).Msg("")
				http.Error(rw, "invalid orderId", http.StatusConflict)
				return
			}

			if errors.Is(err, service.ErrOrderAlreadyUploaded) {
				logger.Trace().Err(err).Msg("")
				rw.Write([]byte("OK"))
				return
			}

			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte("OK"))
	}
}

func (h *Handler) GetOrdersHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "GetOrdersHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		user := appContext.User(ctx)
		if user == nil {
			h.Log(ctx).Err(ErrNotAuthenticated).Msg("")
			http.Error(rw, "user not found", http.StatusUnauthorized)
			return
		}

		orders, err := h.order.FindOrdersByUser(ctx, user.Id)
		if err != nil {
			h.Log(ctx).Err(err).Msg("invalid find users")
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")

		h.Log(ctx).Info().Msg("get orders handled")

		bytes, _ := json.Marshal(orders)

		rw.Write(bytes)
	}
}
