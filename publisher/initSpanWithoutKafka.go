// +build !kafka

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
	default:
		publisher = nil
		return fmt.Errorf("unable to find publisher backend configuration")
	}
}
