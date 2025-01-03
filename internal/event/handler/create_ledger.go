package handler

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CreateLedgerHandler interface {
	Handle(ctx context.Context) error
}

type createLedgerHandler struct {
	logger *zap.Logger
}

func NewCreateLedgerHandlerModule() fx.Option {
	return fx.Provide(
		NewCreateLedgerHandler,
	)
}

func NewCreateLedgerHandler(logger *zap.Logger) CreateLedgerHandler {
	return &createLedgerHandler{
		logger: logger,
	}
}

func (c *createLedgerHandler) Handle(ctx context.Context) error {
	return nil
}
