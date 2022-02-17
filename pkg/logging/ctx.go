package logging

import (
	"context"
	appContext "github.com/djokcik/gophermart/pkg/context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	contextKeyLogger  = appContext.ContextKey("Logger")
	contextKeyTraceID = appContext.ContextKey("TraceID")
)

func GetCtxLogger(ctx context.Context) (context.Context, zerolog.Logger) {
	if ctxValue := ctx.Value(contextKeyLogger); ctxValue != nil {
		if ctxLogger, ok := ctxValue.(zerolog.Logger); ok {
			return ctx, ctxLogger
		}
	}

	traceID, _ := uuid.NewUUID()
	logger := NewLogger().With().Str(TraceIDKey, traceID.String()).Logger()

	ctx = context.WithValue(ctx, contextKeyTraceID, traceID.String())

	return SetCtxLogger(ctx, logger), logger
}

func SetCtxLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}
