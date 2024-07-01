package kafka

import (
	"context"
	"errors"
	log1 "github.com/sniperCore/core/log"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

const (
	DefaultDialTimeout  = 10
	DefaultReadTimeout  = 10
	DefaultWriteTimeout = 10
)

type Kafka struct {
	groupId       string
	topic         string
	ConsumerGroup sarama.ConsumerGroup
	Producer      sarama.SyncProducer
}

func NewKafka(conf *KafkaConfig) (*Kafka, error) {
	//init kafka config
	config := initKafkaConfig(conf)

	//init log
	sarama.Logger = log.New(log1.GetLogger(log1.SingletonMain).Writer(), "sarama ", log.LstdFlags)

	//init consumer client
	group, err := sarama.NewConsumerGroup(conf.Hosts, conf.GroupId, config)
	if err != nil {
		return nil, err
	}

	//init producer client
	producer, err := sarama.NewSyncProducer(conf.Hosts, nil)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		topic:         conf.Topic,
		groupId:       conf.GroupId,
		ConsumerGroup: group,
		Producer:      producer,
	}, nil
}

/**
 * kafka init configs
 */
func initKafkaConfig(conf *KafkaConfig) *sarama.Config {
	//deal with config params
	if conf.DialTimeout == 0 {
		conf.DialTimeout = DefaultDialTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = DefaultReadTimeout
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = DefaultWriteTimeout
	}

	config := sarama.NewConfig()
	config.Net.DialTimeout = time.Duration(conf.DialTimeout) * time.Second
	config.Net.ReadTimeout = time.Duration(conf.ReadTimeout) * time.Second
	config.Net.WriteTimeout = time.Duration(conf.WriteTimeout) * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V0_10_2_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = false

	return config
}

// consumer func
func (kafka *Kafka) Consume(consumerHandler sarama.ConsumerGroupHandler) error {
	if consumerHandler == nil {
		return errors.New("consumerGroupHandler is not empty")
	}

	// Track errors ++ 临时处理，后续error是否需要统一接口处理++
	go func() {
		for err := range kafka.ConsumerGroup.Errors() {
			panic(err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{kafka.topic}
		err := kafka.ConsumerGroup.Consume(ctx, topics, consumerHandler)
		if err != nil {
			return err
		}
	}
}

//producer func
func (kafka *Kafka) Produce(message string) (partition int32, offset int64, err error) {
	if message == "" {
		err = errors.New("producer message is not empty")
		return
	}

	msg := &sarama.ProducerMessage{Topic: kafka.topic, Value: sarama.StringEncoder(message)}
	partition, offset, err = kafka.Producer.SendMessage(msg)

	return
}
