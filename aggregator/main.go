package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gabuladze/tolling/types"
)

func main() {
	listenAddr := ":3001"
	store := NewMemoryStore()
	distAgg := NewDistanceAggregator(store)
	err := makeHTTPTransport(listenAddr, distAgg)
	if err != nil {
		log.Fatal(err)
	}
}

func makeHTTPTransport(listenAddr string, da Aggregator) error {
	log.Println("HTTP server running on ", listenAddr)
	http.HandleFunc("/", handleAggregate(da))
	return http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(da Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
