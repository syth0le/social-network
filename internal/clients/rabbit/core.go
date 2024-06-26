package rabbit

import (
	"fmt"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
)

const defaultContentType = "application/json"

type Publisher interface {
	Publish(msg []byte, routingKey string, withExpiration bool) error
	Close() error
}

type PublisherImpl struct {
	Conn         *rabbitmq.Conn
	Publisher    *rabbitmq.Publisher
	ExchangeName string
}

func NewRabbitPublisher(logger *zap.Logger, cfg configuration.RabbitConfig) (Publisher, error) {
	if !cfg.Enable {
		return &PublisherMock{
			Logger: logger,
		}, nil
	}

	conn, err := rabbitmq.NewConn(
		cfg.Address,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection: %w", err)
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(cfg.ExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create publisher: %w", err)
	}

	return &PublisherImpl{
		Conn:         conn,
		Publisher:    publisher,
		ExchangeName: cfg.ExchangeName,
	}, nil
}

func (p *PublisherImpl) Close() error {
	p.Publisher.Close()
	return p.Conn.Close()
}

func (p *PublisherImpl) Publish(msg []byte, routingKey string, withExpiration bool) error {
	optionFuncs := []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsContentType(defaultContentType),
		rabbitmq.WithPublishOptionsExchange(p.ExchangeName),
	}
	if withExpiration {
		optionFuncs = append(optionFuncs, rabbitmq.WithPublishOptionsExpiration("5000"))
	}

	return p.Publisher.Publish(
		msg,
		[]string{routingKey},
		optionFuncs...,
	)
}

type Consumer interface {
	Close() error
	Run(handler rabbitmq.Handler) error
}

type ConsumerImpl struct {
	Conn     *rabbitmq.Conn
	Consumer *rabbitmq.Consumer
}

func NewRabbitConsumer(logger *zap.Logger, cfg configuration.RabbitConfig, routingKey string) (Consumer, error) {
	if !cfg.Enable {
		return &ConsumerMock{
			Logger: logger,
		}, nil
	}

	conn, err := rabbitmq.NewConn(
		cfg.Address,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection: %w", err)
	}

	consumer, err := rabbitmq.NewConsumer(
		conn,
		cfg.QueueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(cfg.ExchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create consumer: %w", err)
	}

	return &ConsumerImpl{
		Conn:     conn,
		Consumer: consumer,
	}, nil
}

func (c *ConsumerImpl) Close() error {
	c.Consumer.Close()
	return c.Conn.Close()
}

func (c *ConsumerImpl) Run(handler rabbitmq.Handler) error {
	return c.Consumer.Run(handler)
}
