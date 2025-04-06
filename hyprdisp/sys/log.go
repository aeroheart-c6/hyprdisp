package sys

import (
	"context"
	"fmt"
	"log/slog"
)

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

func GetLogger(ctx context.Context) (*slog.Logger, error) {
	var (
		logger *slog.Logger
		ok     bool
		value  any
	)

	value = ctx.Value(ContextKeyLogger)
	if value == nil {
		return nil, fmt.Errorf("no value for context key \"%v\" found", ContextKeyLogger)
	}

	logger, ok = value.(*slog.Logger)
	if !ok {
		return nil, fmt.Errorf("found value with context key \"%v\" but could not convert to logger instance", ContextKeyLogger)
	}

	return logger, nil
}
