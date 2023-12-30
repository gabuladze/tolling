package main

import (
	"context"
	"log"
	"time"

	"github.com/gabuladze/tolling/aggregator/client"
	"github.com/gabuladze/tolling/types"
)

func main() {
	c, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.AggregateDistance(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
