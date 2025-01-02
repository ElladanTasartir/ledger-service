package main

import (
	"context"
	"log"

	"github.com/ElladanTasartir/ledger-service/internal/config"
	"github.com/ElladanTasartir/ledger-service/internal/event/consumer"
	"github.com/ElladanTasartir/ledger-service/internal/event/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.Provide(zap.NewProduction),
		config.NewConfigModule("config"),
		handler.NewCreateTransactionHandlerModule(),
		consumer.NewConsumerModule(),
		fx.Invoke(BootstrapApplication),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatalf("failed to start application. err = %v", err)
	}

	<-app.Done()

	if err := app.Stop(context.Background()); err != nil {
		log.Fatalf("failed to shutdown gracefully. err = %v", err)
	}
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
