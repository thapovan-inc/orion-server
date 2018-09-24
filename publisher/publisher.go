package publisher

import (
	"fmt"
	"github.com/thapovan-inc/orionproto"
)

const (
	KAFKA = "kafka"
	NATS  = "nats"
)

type Publisher interface {
	connect() error
	isConnected() bool
	Publish(topic string, key, value []byte) error
	Close() error
}

type SpanPublisher interface {
	Publisher
	PublishSpan(topic string, key []byte, spanData *orionproto.Span) error
}

var publisher SpanPublisher = nil

func GetSpanPublisher() (SpanPublisher, error) {
	if publisher == nil {
		return nil, fmt.Errorf("publisher not yet initialized. Try publisher::InitPublisherConfig")
	}
	return publisher.(SpanPublisher), nil
}

func GetPublisher() (Publisher, error) {
	return GetSpanPublisher()
}
