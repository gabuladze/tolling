package main

import (
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
	return da.store.Insert(dist)
}
