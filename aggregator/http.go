package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gabuladze/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MakeHTTPTransport(listenAddr string, da Aggregator) error {
	httpHandler := NewHTTPHandler(da)
	aggregateHandler := makeHTTPHandlerFunc(
		newHTTPMetricsMiddleware("aggregate").instrument(httpHandler.handlePostAggregate),
	)
	invoiceHandler := makeHTTPHandlerFunc(
		newHTTPMetricsMiddleware("invoice").instrument(httpHandler.handleGetInvoice),
	)

	log.Println("HTTP server running on ", listenAddr)
	http.HandleFunc("/aggregate", aggregateHandler)
	http.HandleFunc("/invoice", invoiceHandler)
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAddr, nil)
}

type HTTPHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandlerFunc(fn HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

type HTTPMetricsMiddleware struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func newHTTPMetricsMiddleware(route string) *HTTPMetricsMiddleware {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_request_counter", route),
		Name:      "aggregator",
	})
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_error_counter", route),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_request_latency", route),
		Name:      "aggregator",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricsMiddleware{
		reqCounter: reqCounter,
		errCounter: errCounter,
		reqLatency: reqLatency,
	}
}

func (m HTTPMetricsMiddleware) instrument(next HTTPHandlerFunc) HTTPHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(t time.Time) {
			m.reqLatency.Observe(time.Since(t).Seconds())
			m.reqCounter.Inc()
			if err != nil {
				m.errCounter.Inc()
			}
		}(time.Now())
		err = next(w, r)
		return err
	}
}

type APIError struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func (ae APIError) Error() string {
	return ae.Err.Error()
}

type HTTPHandler struct {
	aggregatorService Aggregator
}

func NewHTTPHandler(agg Aggregator) *HTTPHandler {
	return &HTTPHandler{
		aggregatorService: agg,
	}
}

func (h HTTPHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return APIError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid method"),
		}
	}
	obuIDParam := r.URL.Query().Get("obu")
	if len(obuIDParam) == 0 {
		return APIError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid obu id"),
		}
	}

	obuID, err := strconv.Atoi(obuIDParam)
	if err != nil {
		return APIError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid obu id"),
		}
	}

	inv, err := h.aggregatorService.GenerateInvoice(obuID)
	if err != nil {
		return APIError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("invoice generation failed"),
		}
	}

	writeJSON(w, http.StatusOK, inv)
	return nil
}

func (h HTTPHandler) handlePostAggregate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return APIError{
			Code: http.StatusBadRequest,
			Err:  errors.New("invalid method"),
		}
	}
	var dist types.Distance
	if err := json.NewDecoder(r.Body).Decode(&dist); err != nil {
		return APIError{
			Code: http.StatusBadRequest,
			Err:  errors.New("failed to decode request body"),
		}
	}
	if err := h.aggregatorService.AggregateDistance(dist); err != nil {
		return APIError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("failed to aggregate distange"),
		}
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
