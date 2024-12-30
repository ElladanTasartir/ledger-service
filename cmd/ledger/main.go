package main

import (
	"context"

	"github.com/ElladanTasartir/ledger-service/internal/config"
	"github.com/ElladanTasartir/ledger-service/internal/event/consumer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(zap.NewProduction),
		config.NewConfigModule("config"),
		consumer.NewConsumerModule(),
		fx.Invoke(BootstrapApplication),
	).Run()
}

func BootstrapApplication(
	lc fx.Lifecycle,
	consumer consumer.Consumer,
	logger *zap.Logger,
) {
	logger = logger.With(zap.String("at", "Main"))

	lc.Append(fx.StartStopHook(
		func() error {
			go func() {
				if err := consumer.StartConsumers(context.Background()); err != nil {
					logger.Error("failed to start consumers", zap.Error(err))
				}
			}()

			return nil
		},
		func() error {
			logger.Info("Shutting down application")
			return nil
		},
	))
}
