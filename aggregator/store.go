package main

import (
	"errors"

	"github.com/gabuladze/tolling/types"
)

type Storer interface {
	Insert(types.Distance) error
	Get(obuID int) (float64, error)
}

type MemoryStorer struct {
	data map[int]float64
}

func NewMemoryStore() Storer {
	return &MemoryStorer{
		data: make(map[int]float64),
	}
}

func (ms *MemoryStorer) Insert(d types.Distance) error {
	ms.data[d.OBUID] += d.Value
	return nil
}

func (ms *MemoryStorer) Get(obuID int) (float64, error) {
	dist, ok := ms.data[obuID]
	if !ok {
		return 0.0, errors.New("obu not found")
	}

	return dist, nil
}
