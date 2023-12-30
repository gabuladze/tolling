package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gabuladze/tolling/aggregator/client"
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
	aggClient client.Client
}

func NewKafkaConsumer(topic string, dcs DistanceCalculator, aggClient client.Client) (DataConsumer, error) {
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
		consumer:  c,
		dcs:       dcs,
		aggClient: aggClient,
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

		req := &types.AggregateRequest{
			ObuID: int32(d.OBUID),
			Value: distance,
			Unix:  time.Now().Unix(),
		}
		if err := kc.aggClient.AggregateDistance(context.Background(), req); err != nil {
			log.Fatalf("aggregate error: %v", err)
			continue
		}
	}
}
