package main

import (
	"log"
	"time"

	"github.com/gabuladze/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) Aggregator {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricsMiddleware{
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
		next:           next,
	}
}

func (mm MetricsMiddleware) AggregateDistance(d types.Distance) (err error) {
	defer func(t time.Time) {
		mm.reqLatencyAgg.Observe(time.Since(t).Seconds())
		mm.reqCounterAgg.Inc()
		if err != nil {
			mm.errCounterAgg.Inc()
		}
	}(time.Now())
	err = mm.next.AggregateDistance(d)
	return
}

func (mm MetricsMiddleware) GenerateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(t time.Time) {
		mm.reqLatencyCalc.Observe(time.Since(t).Seconds())
		mm.reqCounterCalc.Inc()
		if err != nil {
			mm.errCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = mm.next.GenerateInvoice(obuID)
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (lm LogMiddleware) AggregateDistance(d types.Distance) error {
	defer func(t time.Time) {
		log.Printf("AggregateDistance \t took=%v OBUID=%d Value=%.2f Unix=%d", time.Since(t), d.OBUID, d.Value, d.Unix)
	}(time.Now())
	return lm.next.AggregateDistance(d)
}

func (lm LogMiddleware) GenerateInvoice(obuID int) (*types.Invoice, error) {
	defer func(t time.Time) {
		log.Printf("GenerateInvoice \t took=%v OBUID=%d", time.Since(t), obuID)
	}(time.Now())
	return lm.next.GenerateInvoice(obuID)
}
