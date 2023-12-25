package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gabuladze/tolling/types"
	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:3000"

var sendInterval = time.Second

func genRandCoord() float64 {
	return float64(rand.Intn(100)) + rand.Float64()
}

func genLatLong() (float64, float64) {
	return genRandCoord(), genRandCoord()
}

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ws connection established")

	obuIDs := genOBUIDs(20)
	for {
		for _, id := range obuIDs {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: id,
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}
