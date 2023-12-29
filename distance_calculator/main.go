package main

import (
	"log"

	"github.com/gabuladze/tolling/aggregator/client"
)

const (
	kafkaTopic     = "obudata"
	clientEndpoint = "http://localhost:3001/aggregate"
)

func main() {
	distanceCalculatorService := NewDistanceCalculatorService()
	aggClient := client.NewHTTPClient(clientEndpoint)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, distanceCalculatorService, aggClient)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
	log.Println("calculating distance")
}
