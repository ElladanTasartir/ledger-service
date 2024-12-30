package main

import (
	"fmt"

	"github.com/ElladanTasartir/ledger-service/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to setup logger")
		panic(1)
	}

	config, err := config.NewConfig("config")
	if err != nil {
		logger.Panic("Failed to setup config", zap.Error(err))
	}

	logger.Info("Start", zap.Any("config", config))
}
