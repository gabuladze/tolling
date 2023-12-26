package main

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gabuladze/tolling/types"
)

type DataConsumer interface {
	Start()
	readMessageLoop()
}

type KafkaConsumer struct {
	isRunning bool
	consumer  *kafka.Consumer
	dcs       DistanceCalculator
}

func NewKafkaConsumer(topic string, dcs DistanceCalculator) (DataConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer: c,
		dcs:      dcs,
	}, nil
}

func (kc *KafkaConsumer) Start() {
	log.Println("starting kafka consumer")
	kc.isRunning = true
	kc.readMessageLoop()
}

func (kc *KafkaConsumer) readMessageLoop() {
	for kc.isRunning {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			log.Println("kafka readMessage error", err)
		}

		var d types.OBUData
		err = json.Unmarshal(msg.Value, &d)
		if err != nil {
			log.Println("json unmarshal error", err)
		}

		distance, err := kc.dcs.CalculateDistance(d)
		log.Printf("calculating distance obuID=%v distance=%.2f", d.OBUID, distance)
		if err != nil {
			log.Println(err)
		}

		_ = distance
	}
}
