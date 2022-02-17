package context

import (
	"context"
	"fmt"
	"github.com/djokcik/gophermart/internal/model"
)

type ContextKey string

const (
	userKey ContextKey = "userKey"
)

func (c ContextKey) String() string {
	return fmt.Sprintf("%s%s", contextKeyPrefix, string(c))
}

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *model.User {
	if ctxValue := ctx.Value(userKey); ctxValue != nil {
		if user, ok := ctxValue.(*model.User); ok {
			return user
		}
	}

	return nil
}

const (
	contextKeyPrefix = "gophermartLogging-"
)
