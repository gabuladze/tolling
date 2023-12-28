package main

import "github.com/gabuladze/tolling/types"

type Storer interface {
	Insert(types.Distance) error
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
