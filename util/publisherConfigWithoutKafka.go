// +build !kafka

package util

type PublisherConfig struct {
	Type                string              `toml:"type"`
	DebugStream         bool                `toml:"debug_stream"`
	NatsPublisherConfig NatsPublisherConfig `toml:"nats"`
}
