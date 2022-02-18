package handler

import (
	"encoding/json"
	"errors"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/service"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/logging"
	"net/http"
)

func (h *Handler) RegisterUserHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "RegisterUserHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var userDto model.UserRequestDto
		err := json.NewDecoder(r.Body).Decode(&userDto)
		if err != nil {
			logger.Trace().Err(err).Msg("")
			http.Error(rw, "invalid parse decoder", http.StatusBadRequest)
			return
		}

		err = h.user.CreateUser(ctx, userDto.Login, userDto.Password)
		if err != nil {
			if errors.Is(err, storage.ErrLoginAlreadyExists) {
				logger.Trace().Err(err).Msg("login already exists")
				http.Error(rw, "login already exists", http.StatusConflict)
				return
			}

			logger.Trace().Err(err).Msg("failed created user")
			http.Error(rw, "invalid create user", http.StatusBadRequest)
			return
		}

		user, err := h.user.GetUserByUsername(ctx, userDto.Login)
		if err != nil {
			logger.Trace().Err(err).Msg("failed to find by username")
			http.Error(rw, "failed find by username", http.StatusInternalServerError)
			return
		}

		token, err := h.user.GenerateToken(ctx, user)
		if err != nil {
			logger.Error().Err(err).Msg("invalid generate token")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{Name: CookieName, Value: token}
		http.SetCookie(rw, &cookie)
		rw.Header().Set("Content-Type", "application/json")

		bytes, _ := json.Marshal(model.UserResponseDto{Token: token})
		rw.Write(bytes)
	}
}

func (h *Handler) SignInHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := h.Log(ctx).With().Str(logging.ServiceKey, "SignInHandler").Logger()
		ctx = logging.SetCtxLogger(ctx, logger)

		var user model.UserRequestDto
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Trace().Err(err).Msg("failed parse data")
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := h.user.Authenticate(ctx, user.Login, user.Password)
		if err != nil {
			if errors.Is(err, service.ErrWrongPassword) {
				logger.Trace().Err(err).Msg("invalid password")
				http.Error(rw, "invalid password", http.StatusUnauthorized)
				return
			}

			logger.Trace().Err(err).Msg("invalid authenticate")
			http.Error(rw, "", http.StatusBadRequest)
			return
		}

		cookie := http.Cookie{Name: CookieName, Value: token}
		http.SetCookie(rw, &cookie)

		rw.Header().Set("Content-Type", "application/json")

		bytes, _ := json.Marshal(model.UserResponseDto{Token: token})
		rw.Write(bytes)
	}
}
