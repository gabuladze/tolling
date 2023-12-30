package main

import (
	"context"

	"github.com/gabuladze/tolling/types"
)

type GRPCAggregatorService struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCService(svc Aggregator) *GRPCAggregatorService {
	return &GRPCAggregatorService{
		svc: svc,
	}
}

func (s *GRPCAggregatorService) AggregateDistance(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	dist := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(dist)
}
