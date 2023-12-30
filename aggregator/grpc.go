package main

import "github.com/gabuladze/tolling/types"

type GRPCAggregatorService struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCService(svc Aggregator) *GRPCAggregatorService {
	return &GRPCAggregatorService{
		svc: svc,
	}
}

func (s *GRPCAggregatorService) AggregateDistance(req *types.AggregateRequest) error {
	dist := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return s.svc.AggregateDistance(dist)
}
