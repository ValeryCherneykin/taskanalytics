package config

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
)

type KafkaProducerConfig interface {
	Brokers() []string
	Config() *sarama.Config
}

type kafkaProducerConfig struct {
	brokers []string
}

func NewKafkaProducerConfig() (KafkaProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers not found")
	}

	brokers := strings.Split(brokersStr, ",")

	return &kafkaProducerConfig{
		brokers: brokers,
	}, nil
}

func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	return config
}
