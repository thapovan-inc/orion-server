// +build kafka

package publisher

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
