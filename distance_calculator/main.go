package main

import (
	"log"
)

const kafkaTopic = "obudata"

func main() {
	distanceCalculatorService := NewDistanceCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, distanceCalculatorService)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
	log.Println("calculating distance")
}
