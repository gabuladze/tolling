package main

import (
	"log"
	"net/http"

	"github.com/gabuladze/tolling/types"
	"github.com/gorilla/websocket"
)

const listenAddr = ":3000"

type receiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

func newReceiver() *receiver {
	return &receiver{
		msgch: make(chan types.OBUData, 128),
	}
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
		rec.msgch <- data
	}
}

func main() {
	rec := newReceiver()
	http.HandleFunc("/", rec.handleWS)
	log.Println("starting receiver at ", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatal(err)
	}
}
