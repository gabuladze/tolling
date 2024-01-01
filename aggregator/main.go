package main

import (
	"log"
)

func main() {
	httpListenAddr := ":3001"
	grpcListenAddr := ":3002"
	store := NewMemoryStore()
	distAgg := NewDistanceAggregator(store)
	distAgg = NewLogMiddleware(distAgg)
	distAgg = NewMetricsMiddleware(distAgg)

	go func() {
		log.Fatal(MakeGRPCTransport(grpcListenAddr, distAgg))
	}()
	log.Fatal(MakeHTTPTransport(httpListenAddr, distAgg))
}
