// +build kafka

package util

import "github.com/confluentinc/confluent-kafka-go/kafka"

type PublisherConfig struct {
	Type                 string              `toml:"type"`
	DebugStream          bool                `toml:"debug_stream"`
	NatsPublisherConfig  NatsPublisherConfig `toml:"nats"`
	KafkaPublisherConfig kafka.ConfigMap     `toml:"kafka"`
}
