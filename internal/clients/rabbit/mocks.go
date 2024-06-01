package rabbit

import (
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type PublisherMock struct {
	Logger *zap.Logger
}

func (m *PublisherMock) Publish(msg []byte) error {
	m.Logger.Sugar().Debugf("published through publisher mock: %b", msg)
	return nil
}

func (m *PublisherMock) Close() error {
	m.Logger.Debug("closed publisher mock")
	return nil
}

type ConsumerMock struct {
	Logger *zap.Logger
}

func (m *ConsumerMock) Close() error {
	m.Logger.Debug("closed consumer mock")
	return nil
}

func (m *ConsumerMock) Run(handler rabbitmq.Handler) error {
	m.Logger.Debug("run through rabbitmq mock")
	return nil
}
