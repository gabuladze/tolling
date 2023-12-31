package client

import (
	"context"

	"github.com/gabuladze/tolling/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial("localhost:3002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) AggregateDistance(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := c.client.AggregateDistance(ctx, aggReq)
	return err
}
