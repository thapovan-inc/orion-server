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

// +build kafka

package publisher

import (
	"fmt"
	"github.com/thapovan-inc/orion-server/util"
)

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
