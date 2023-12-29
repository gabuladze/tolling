package main

import (
	"github.com/gabuladze/tolling/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
	GenerateInvoice(obuID int) (*types.Invoice, error)
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

func (da DistanceAggregator) GenerateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := da.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   dist * 1.2,
	}
	return &inv, nil
}
