package consumer

import (
	"context"
	"fmt"

	"github.com/ElladanTasartir/ledger-service/internal/common"
	"github.com/ElladanTasartir/ledger-service/internal/config"
	"github.com/ElladanTasartir/ledger-service/internal/event/handler"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Consumer interface {
	StartConsumers(ctx context.Context) error
}

type consumer struct {
	logger                    *zap.Logger
	consumerURL               string
	consumerGroupID           string
	consumerTopics            []string
	consumer                  *kafka.Consumer
	createrTransactionHandler handler.CreateTransactionHandler
}

func NewConsumerModule() fx.Option {
	return fx.Options(
		fx.Provide(NewConsumer),
	)
}

func NewConsumer(config *config.Config, createTransactionHandler handler.CreateTransactionHandler, logger *zap.Logger) (Consumer, error) {
	consumer := &consumer{
		consumerURL:               config.Kafka.URL,
		consumerGroupID:           config.Kafka.GroupID,
		consumerTopics:            config.Kafka.ConsumerTopics,
		createrTransactionHandler: createTransactionHandler,
		logger: logger.With(
			zap.String("at", "Consumer"),
		),
	}

	if len(consumer.consumerURL) == 0 {
		return nil, fmt.Errorf("missing consumer URL")
	}

	if len(consumer.consumerGroupID) == 0 {
		return nil, fmt.Errorf("missing consumer group ID")
	}

	return consumer, nil
}

func (c *consumer) StartConsumers(ctx context.Context) error {
	var err error
	c.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": c.consumerURL,
		"group.id":          c.consumerGroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		c.logger.Panic("Failed to bootstrap consumer", zap.Error(err))
	}
	defer c.consumer.Close()

	return c.subscribeToTopics(ctx)
}

func (c *consumer) subscribeToTopics(ctx context.Context) error {
	if len(c.consumerTopics) == 0 {
		c.logger.Warn("no topics to subscribe to")
		return nil
	}

	c.consumer.SubscribeTopics(c.consumerTopics, nil)
	for {
		ctx, cancel := context.WithTimeout(ctx, common.ConsumerHandleTimeout)
		defer cancel()

		msg, err := c.consumer.ReadMessage(common.ConsumerReadTimeout)
		if err != nil {
			c.logger.Error("failed to read message", zap.Error(err))
			return err
		}

		switch *msg.TopicPartition.Topic {
		case common.CreateLedgerTopic:
			if err := c.createrTransactionHandler.Handle(ctx); err != nil {
				c.logger.Error("failed to create transaction", zap.Error(err))
				continue
			}

			c.logger.Info("consume create transaction event")
		default:
			c.logger.Warn("Unsupported handler method")
		}
	}
}
