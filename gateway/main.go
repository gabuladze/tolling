package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gabuladze/tolling/aggregator/client"
	"github.com/gabuladze/tolling/types"
)

type InvoiceHandler struct {
	client client.Client
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	c := client.NewHTTPClient("http://localhost:3001")
	invoiceHandler := NewInvoiceHandler(c)
	http.HandleFunc("/invoice", makeAPIHandlerFunc(invoiceHandler.handleGetInvoice))
	log.Println("http server listening")
	log.Fatal(http.ListenAndServe(":3003", nil))
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (c *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	obuIDParam := r.URL.Query().Get("obu")
	if len(obuIDParam) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
		return nil
	}

	obuID, err := strconv.Atoi(obuIDParam)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
		return nil
	}
	req := types.GetInvoiceRequest{
		ObuID: int64(obuID),
	}
	inv, err := c.client.GetInvoice(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return nil
	}
	writeJSON(w, http.StatusOK, inv)
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeAPIHandlerFunc(fn handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
