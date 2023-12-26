package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gabuladze/tolling/types"
)

type DataProducer interface {
	ProduceData(types.OBUData) error
}

type KafkaDataProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaDataProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	return &KafkaDataProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (kdp *KafkaDataProducer) ProduceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return kdp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)
}
