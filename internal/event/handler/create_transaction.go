package handler

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CreateTransactionHandler interface {
	Handle(ctx context.Context) error
}

type createTransactionHandler struct {
	logger *zap.Logger
}

func NewCreateTransactionHandlerModule() fx.Option {
	return fx.Provide(
		NewCreateTransactionHandler,
	)
}

func NewCreateTransactionHandler(logger *zap.Logger) CreateTransactionHandler {
	return &createTransactionHandler{
		logger: logger.With(
			zap.String("at", "CreateTransactionHandler")),
	}
}

func (c *createTransactionHandler) Handle(ctx context.Context) error {
	c.logger.Info("handled message")
	return nil
}
