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

type HandleMessage = func(ctx context.Context, message []byte) error

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
	createLedgerHandler       handler.CreateLedgerHandler
	supportedHandlers         map[string]HandleMessage
}

func NewConsumerModule() fx.Option {
	return fx.Options(
		fx.Provide(NewConsumer),
	)
}

func NewConsumer(
	config *config.Config,
	createLedgerHandler handler.CreateLedgerHandler,
	createTransactionHandler handler.CreateTransactionHandler,
	logger *zap.Logger,
) (Consumer, error) {
	consumer := &consumer{
		consumerURL:               config.Kafka.URL,
		consumerGroupID:           config.Kafka.GroupID,
		consumerTopics:            config.Kafka.ConsumerTopics,
		createrTransactionHandler: createTransactionHandler,
		createLedgerHandler:       createLedgerHandler,
		logger: logger.With(
			zap.String("at", "Consumer"),
		),
	}

	consumer.supportedHandlers = map[string]HandleMessage{
		common.CreateLedgerTopic:      consumer.handleCreateLedgerMessage,
		common.CreateTransactionTopic: consumer.handleCreateTransactionMessage,
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
		"bootstrap.servers":  c.consumerURL,
		"group.id":           c.consumerGroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
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

		if err := c.handleMessage(ctx, msg); err != nil {
			c.logger.Error("failed to handle message", zap.Error(err))
		} else {
			c.logger.Info("message consumed successfully")
		}
	}
}

func (c *consumer) handleMessage(ctx context.Context, message *kafka.Message) error {
	topic := *message.TopicPartition.Topic
	logger := c.logger.With(zap.String("topic", topic))

	if messageHandler, ok := c.supportedHandlers[topic]; ok {
		if err := messageHandler(ctx, message.Value); err != nil {
			logger.Error("failed to handle message", zap.Error(err))
			return err
		}

		if _, err := c.consumer.CommitMessage(message); err != nil {
			logger.Error("failed to commit message")
			return err
		}

		logger.Info("consume event")
		return nil
	}

	logger.Warn("unsupported handler method")
	return nil
}

func (c *consumer) handleCreateLedgerMessage(ctx context.Context, message []byte) error {
	return nil
}

func (c *consumer) handleCreateTransactionMessage(ctx context.Context, message []byte) error {
	return nil
}
