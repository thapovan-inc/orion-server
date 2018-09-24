// Copyright 2018-Present Thapovan Info Systems Pvt. Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package publisher

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"github.com/vmihailenco/msgpack"
	"log"
)

type NatsMessage struct {
	Key   []byte
	Value []byte
}

type NatsPublisher struct {
	URL                string
	ClusterID          string
	ClientID           string
	nc                 stan.Conn
	connected          bool
	debugStreamEnabled bool
}

func (n *NatsPublisher) connect() error {
	logger := util.GetLogger("publisher", "NatsPublisher::connect")
	logger.Info("Connecting")
	if n.URL == "" {
		logger.Warn("nats URL not provided. Using default URL ", stan.DefaultNatsURL)
		n.URL = stan.DefaultNatsURL
	}
	nc, err := stan.Connect(n.ClusterID, n.ClientID, stan.NatsURL(n.URL), stan.SetConnectionLostHandler(func(conn stan.Conn, e error) {
		logger.Errorf("Connection to nats server lost. Error: %v", e)
		n.connected = false
	}))
	if err != nil {
		n.connected = false
		logger.Error("Error when trying to connect to nats server at ", n.URL)
	} else {
		n.nc = nc
		n.connected = true
	}
	return err
}

func (n *NatsPublisher) isConnected() bool {
	return n.connected
}

func (n *NatsPublisher) publishMessage(topic string, key, messageBytes []byte) error {
	logger := util.GetLogger("publisher", "NatsPublisher::publishMessage")
	var err error
	_, err = n.nc.PublishAsync(topic, messageBytes, func(rGuid string, e error) {
		logger.Debugf("Published key %s to topic %s", string(key), topic)
		if err != nil {
			log.Fatalf("Received error when trying to publish key %s with guid %s to topic %s. Error: %v",
				string(key), rGuid, topic, e)
		}
	})
	return err
}

func (n *NatsPublisher) Publish(topic string, key, value []byte) error {
	logger := util.GetLogger("publisher", "NatsPublisher::Publish")
	if !n.isConnected() {
		logger.Info("Publisher not connected. Attempting to connect")
		err := n.connect()
		if err != nil {
			return err
		}
	}
	messageBytes, err := msgpack.Marshal(&NatsMessage{Key: key, Value: value})
	if err != nil {
		logger.WithError(err).Errorln("Error when marshaling message")
		return err
	}
	return n.publishMessage(topic, key, messageBytes)
}

func (n *NatsPublisher) PublishSpan(topic string, key []byte, spanData *orionproto.Span) error {
	data, _ := proto.Marshal(spanData)
	if n.debugStreamEnabled {
		jsonData, _ := orionproto.ProtoToJson(spanData)
		go func() {
			n.Publish(topic+"-debug", key, []byte(jsonData))
		}()
	}
	return n.Publish(topic, key, data)
}

func (n *NatsPublisher) Close() error {
	if n.isConnected() {
		return n.nc.Close()
	}
	return nil
}
