package main

import (
	"log"
	"net/http"

	"github.com/gabuladze/tolling/types"
	"github.com/gorilla/websocket"
)

const (
	listenAddr  = ":3000"
	kafkaServer = "localhost"
)

var kafkaTopic = "obudata"

type receiver struct {
	prod DataProducer
	conn *websocket.Conn
}

func newReceiver() (*receiver, error) {
	var (
		p   DataProducer
		err error
	)
	p, err = NewKafkaDataProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)

	return &receiver{
		prod: p,
	}, nil
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
		rec.prod.ProduceData(data)
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
