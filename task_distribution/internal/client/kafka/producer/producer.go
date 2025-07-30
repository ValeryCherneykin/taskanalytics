package producer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) *kafkaProducer {
	return &kafkaProducer{producer: producer}
}

func (p *kafkaProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("[Kafka Producer] failed to send message: %v", err)
		return err
	}

	log.Printf("[Kafka Producer] message sent to topic %s [partition: %d, offset: %d]", topic, partition, offset)
	return nil
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
