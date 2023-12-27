package main

import (
	"log"

	"github.com/gabuladze/tolling/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type DistanceAggregator struct {
	store Storer
}

func NewDistanceAggregator(store Storer) Aggregator {
	return &DistanceAggregator{
		store: store,
	}
}

func (da DistanceAggregator) AggregateDistance(dist types.Distance) error {
	log.Panicln("aggregateDistance", dist)
	return da.store.Insert(dist)
}
