package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gabuladze/tolling/types"
	"github.com/gorilla/websocket"
)

const (
	listenAddr  = ":3000"
	kafkaServer = "localhost"
)

var kafkaTopic = "obudata"

type receiver struct {
	prod *kafka.Producer
	conn *websocket.Conn
}

func newReceiver() (*receiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return &receiver{
		prod: p,
	}, nil
}

func (rec *receiver) produceMsg(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = rec.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          b,
	}, nil)

	return err
}

func (rec *receiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	rec.conn = conn

	go rec.receiveLoop()
}

func (rec *receiver) receiveLoop() {
	log.Println("OBU connected!")
	for {
		var data types.OBUData
		if err := rec.conn.ReadJSON(&data); err != nil {
			log.Println("error when reading data: ", err)
			continue
		}
		log.Printf("received OBU data from [%d]: <lat: %.2f, long: %.2f>\n", data.OBUID, data.Lat, data.Long)
		rec.produceMsg(data)
	}
}

func main() {
	rec, err := newReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rec.handleWS)
	log.Println("starting receiver at ", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
