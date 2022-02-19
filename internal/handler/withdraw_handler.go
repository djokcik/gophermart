package handler

import (
	"encoding/json"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/logging"
	"net/http"
)

func (h *Handler) WithdrawHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "WithdrawHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var withdrawDto model.WithdrawRequestDto
		err := json.NewDecoder(r.Body).Decode(&withdrawDto)
		if err != nil {
			h.Log(ctx).Trace().Err(err).Msg("WithdrawHandler: invalid parse body")
			http.Error(rw, "invalid parse body", http.StatusBadRequest)

			return
		}

		err = withdrawDto.Validate()
		if err != nil {
			h.Log(ctx).Trace().Err(err).Msg("validate orderID")
			if errors.Is(err, model.ErrInvalidOrderID) {
				http.Error(rw, "invalid orderID", http.StatusUnprocessableEntity)
				return
			}

			http.Error(rw, "invalid request body", http.StatusBadRequest)
			return
		}

		user := appContext.User(ctx)
		if user == nil {
			h.Log(ctx).Trace().Err(ErrNotAuthenticated).Msg("")
			http.Error(rw, "user not found", http.StatusUnauthorized)

			return
		}

		err = h.withdraw.ProcessWithdraw(ctx, withdrawDto.OrderID, withdrawDto.Sum)
		if err != nil {
			if errors.Is(err, storage.ErrInsufficientFunds) {
				h.Log(ctx).Trace().Err(storage.ErrInsufficientFunds).Msg("WithdrawHandler:")
				http.Error(rw, "insufficient funds", http.StatusPaymentRequired)

				return
			}

			h.Log(ctx).Trace().Err(err).Msg("WithdrawHandler:")
			http.Error(rw, "insufficient funds", http.StatusInternalServerError)
			return
		}

		rw.Write([]byte("OK"))
	}
}

func (h *Handler) WithdrawLogsHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "WithdrawLogsHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		user := appContext.User(ctx)
		if user == nil {
			h.Log(ctx).Trace().Err(ErrNotAuthenticated).Msg("")
			http.Error(rw, "user not found", http.StatusUnauthorized)

			return
		}

		withdrawLogs, err := h.withdraw.WithdrawLogsByUserID(ctx, user.ID)
		if err != nil {
			h.Log(ctx).Err(err).Msg("invalid find users")
			http.Error(rw, "internal error", http.StatusInternalServerError)
			return
		}

		h.Log(ctx).Info().Msg("get withdraw logs handled")

		if len(withdrawLogs) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		rw.Header().Set("Content-Type", "application/json")

		bytes, _ := json.Marshal(withdrawLogs)
		rw.Write(bytes)
	}
}
