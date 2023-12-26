package main

import (
	"log"
	"time"

	"github.com/gabuladze/tolling/types"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}
func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		log.Printf("producing to kafka took=%v OBUID=%d lat=%.2f long=%.2f\n", time.Since(start), data.OBUID, data.Lat, data.Long)
	}(time.Now())
	return l.next.ProduceData(data)
}
