package publisher

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gogo/protobuf/proto"
	"github.com/thapovan-inc/orion-server/util"
	"github.com/thapovan-inc/orionproto"
	"sync"
)

type KafkaPublisher struct {
	ConfigMap          kafka.ConfigMap
	connected          bool
	producer           *kafka.Producer
	closeChan          chan int
	pubEventChan       chan kafka.Event
	deliveryHandlerWG  sync.WaitGroup
	debugStreamEnabled bool
}

func (p *KafkaPublisher) connect() error {
	logger := util.GetLogger("publisher", "KafkaPublisher::connect")
	logger.Info("Connecting")
	if p.ConfigMap == nil {
		panic("KafkaPublisher ConfigMap is nil")
	}
	var err error

	value := p.ConfigMap["go.produce.channel.size"]
	if value != nil {
		int64Value, valid := value.(int64)
		if valid {
			p.ConfigMap["go.produce.channel.size"] = int(int64Value)
		}
	}
	p.producer, err = kafka.NewProducer(&p.ConfigMap)
	if err != nil {
		p.connected = false
		logger.WithError(err).Error("Error when trying to connect to kafka server")
	} else {
		closeChan := make(chan int, 1)
		pubEventChan := make(chan kafka.Event, 16)
		p.closeChan = closeChan
		p.pubEventChan = pubEventChan
		p.closeChan <- 1
		p.deliveryHandlerWG.Wait()
		go p.deliveryHandler()
		p.connected = true
	}
	return err
}

func (p *KafkaPublisher) deliveryHandler() {
	logger := util.GetLogger("publisher", "KafkaPublisher::deliveryHandler")
	p.deliveryHandlerWG.Add(1)
	defer p.deliveryHandlerWG.Done()
HandlerLoop:
	for {
		select {
		case _ = <-p.closeChan:
			break HandlerLoop
		case ev, ok := <-p.producer.Events():
			if !ok {
				p.pubEventChan = nil
				continue
			}
			switch event := ev.(type) {
			case *kafka.Message:
				if event.TopicPartition.Error != nil {
					logger.WithError(event.TopicPartition.Error).Error("Received error when trying to publish key: ",
						string(event.Key), " to topic ", event.TopicPartition.Topic)
				} else {
					logger.Debugln("Published ", string(event.Key), " to topic ", *event.TopicPartition.Topic)
				}
			}
		}
	}
}

func (p *KafkaPublisher) isConnected() bool {
	return p.connected
}

func (p *KafkaPublisher) Publish(topic string, key, value []byte) error {
	logger := util.GetLogger("publisher", "KafkaPublisher::Publish")
	if !p.isConnected() {
		p.connect()
		logger.Info("Publisher not connected. Attempting to connect")
		err := p.connect()
		if err != nil {
			return err
		}
	}
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}
	p.producer.ProduceChannel() <- msg
	return nil
}

func (p *KafkaPublisher) PublishSpan(topic string, key []byte, spanData *orionproto.Span) error {
	data, _ := proto.Marshal(spanData)
	if p.debugStreamEnabled {
		jsonData, _ := orionproto.ProtoToJson(spanData)
		go func() {
			p.Publish(topic+"-debug", key, []byte(jsonData))
		}()
	}
	return p.Publish(topic, key, data)
}

func (p *KafkaPublisher) Close() error {
	if p.isConnected() {
		p.producer.Close()
		close(p.pubEventChan)
		p.closeChan <- 1
		p.deliveryHandlerWG.Wait()
		p.connected = false
	}
	return nil
}
