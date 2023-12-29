package main

import (
	"log"
	"time"

	"github.com/gabuladze/tolling/types"
)

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
