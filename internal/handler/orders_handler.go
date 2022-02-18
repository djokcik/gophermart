package handler

import (
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service"
	"github.com/djokcik/gophermart/pkg/logging"
	"io"
	"net/http"
)

func (h *Handler) OrdersHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "OrdersHandler").Logger()
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
