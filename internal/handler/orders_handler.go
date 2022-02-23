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

		orderID := model.OrderID(body)
		if !orderID.Valid() {
			logger.Trace().Err(err).Msg("invalid orderID")
			http.Error(rw, "invalid orderID", http.StatusUnprocessableEntity)
			return
		}

		err = h.order.ProcessOrder(ctx, orderID)
		if err != nil {
			if errors.Is(err, service.ErrOrderAlreadyUploadedAnotherUser) {
				logger.Trace().Err(err).Msg("")
				http.Error(rw, "order already uploaded", http.StatusConflict)
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

		orders, err := h.order.OrdersByUser(ctx, user.ID)
		if err != nil {
			h.Log(ctx).Err(err).Msg("invalid find users")
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		h.Log(ctx).Info().Msg("get orders handled")

		if len(orders) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		rw.Header().Set("Content-Type", "application/json")

		bytes, _ := json.Marshal(orders)
		rw.Write(bytes)
	}
}
