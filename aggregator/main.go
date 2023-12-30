package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gabuladze/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := ":3001"
	grpcListenAddr := ":3002"
	store := NewMemoryStore()
	distAgg := NewDistanceAggregator(store)
	distAgg = NewLogMiddleware(distAgg)

	go makeGRPCTransport(grpcListenAddr, distAgg)
	makeHTTPTransport(httpListenAddr, distAgg)
}

func makeGRPCTransport(listenAddr string, da Aggregator) error {
	log.Println("GRPC service running on ", listenAddr)
	ln, err := net.Listen("TCP", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewAggregatorGRPCService(da))

	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, da Aggregator) error {
	log.Println("HTTP server running on ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(da))
	http.HandleFunc("/invoice", handleGetInvoice(da))
	return http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(da Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuIDParam := r.URL.Query().Get("obu")
		if len(obuIDParam) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}

		obuID, err := strconv.Atoi(obuIDParam)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}

		inv, err := da.GenerateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "invoice generation failed"})
			return
		}

		writeJSON(w, http.StatusOK, inv)
		return
	}
}

func handleAggregate(da Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dist types.Distance
		if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := da.AggregateDistance(dist); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
