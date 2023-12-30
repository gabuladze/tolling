package client

import (
	"context"

	"github.com/gabuladze/tolling/types"
)

type Client interface {
	AggregateDistance(context.Context, *types.AggregateRequest) error
}
