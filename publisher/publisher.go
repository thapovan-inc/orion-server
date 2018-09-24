package publisher

import (
	"fmt"
	"github.com/thapovan-inc/orion-server/util"
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

func InitSpanPublisherFromConfig() error {
	logger := util.GetLogger("publisher", "InitSpanPublisherFromConfig")
	serverConfig := util.GetConfig()
	switch serverConfig.PublisherConfig.Type {
	case NATS:
		natsConfig := serverConfig.PublisherConfig.NatsPublisherConfig
		publisher = &NatsPublisher{URL: natsConfig.URL, ClientID: natsConfig.ClientID, ClusterID: natsConfig.ClusterID,
			debugStreamEnabled: serverConfig.PublisherConfig.DebugStream}
		err := publisher.connect()
		if err != nil {
			logger.Debug(err)
			return err
		} else {
			return nil
		}
	case KAFKA:
		kafkaConfig := serverConfig.PublisherConfig.KafkaPublisherConfig
		publisher = &KafkaPublisher{ConfigMap: kafkaConfig,
			debugStreamEnabled: serverConfig.PublisherConfig.DebugStream}
		err := publisher.connect()
		if err != nil {
			logger.WithError(err).Debug("Error when connecting")
			return err
		} else {
			return nil
		}
	default:
		publisher = nil
		return fmt.Errorf("unable to find publisher backend configuration")
	}
}

func GetSpanPublisher() (SpanPublisher, error) {
	if publisher == nil {
		return nil, fmt.Errorf("publisher not yet initialized. Try publisher::InitPublisherConfig")
	}
	return publisher.(SpanPublisher), nil
}

func GetPublisher() (Publisher, error) {
	return GetSpanPublisher()
}
