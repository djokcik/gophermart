package middleware

import (
	"github.com/djokcik/gophermart/internal/config"
	"github.com/djokcik/gophermart/internal/handler"
	"github.com/djokcik/gophermart/internal/model"
	"github.com/djokcik/gophermart/internal/storage"
	"github.com/djokcik/gophermart/pkg/context"
	"github.com/djokcik/gophermart/pkg/encrypt"
	"github.com/djokcik/gophermart/pkg/logging"
	"net/http"
)

func UserContext(userRepo storage.UserRepository, cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			_, logger := logging.GetCtxLogger(ctx)
			logger = logger.With().Str(logging.ServiceKey, "UserContext").Logger()

			var token string

			cookie, err := r.Cookie(handler.CookieName)
			if err != nil {
				token, err = encrypt.GetJwtTokenByAuthHeader(r.Header.Get("Authorization"))
				if err != nil {
					logger.Trace().Err(err).Msg("RequireUser: don`t have token")
					http.Error(rw, "Unauthorized", http.StatusUnauthorized)
					return
				}
			} else {
				token = cookie.Value
			}

			id, err := encrypt.ParseToken(token, cfg.Key)
			if err != nil {
				status := http.StatusBadRequest
				if err == model.ErrInvalidAccessToken {
					status = http.StatusUnauthorized
				}

				logger.Trace().Err(err).Msg("RequireUser: invalid token")
				http.Error(rw, "Unauthorized", status)
				return
			}

			user, err := userRepo.UserByID(ctx, id)
			if err != nil {
				logger.Trace().Err(err).Msgf("RequireUser: user with id %d not found", id)
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)

				return
			}

			ctx = context.WithUser(r.Context(), &user)
			logger.Trace().Str("user", user.Username).Msg("RequireUser: successfully authorized")
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
