package publisher

import (
	"github.com/thapovan-inc/orionproto"
	"go.uber.org/zap"
)

type consolePublisher struct {
	logger *zap.SugaredLogger
}

func (cp *consolePublisher) connect() error {
	if cp.logger == nil {
		config := zap.NewDevelopmentConfig()
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.Encoding = "console"
		logger, _ := config.Build()
		cp.logger = logger.Sugar()
	}
	return nil
}

func (cp *consolePublisher) isConnected() bool {
	return cp.logger != nil
}

func (cp *consolePublisher) Publish(topic string, key, value []byte) error {
	cp.logger.Info("Topic: ", topic, "key: ", string(key), " value: ", string(value))
	return nil
}

func (cp *consolePublisher) Close() error {
	cp.logger.Sync()
	cp.logger = nil
	return nil
}

func (cp *consolePublisher) PublishSpan(topic string, key []byte, spanData *orionproto.Span) error {
	json, _ := orionproto.ProtoToJson(spanData)
	cp.logger.Info("Received Span Event: ", json)
	return nil
}
