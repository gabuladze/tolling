package main

import "github.com/gabuladze/tolling/types"

type Storer interface {
	Insert(types.Distance) error
}

type MemoryStorer struct{}

func NewMemoryStore() Storer {
	return &MemoryStorer{}
}

func (ms MemoryStorer) Insert(d types.Distance) error {
	return nil
}
