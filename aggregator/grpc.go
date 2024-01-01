package main

import (
	"context"
	"log"
	"net"

	"github.com/gabuladze/tolling/types"
	"google.golang.org/grpc"
)

func MakeGRPCTransport(listenAddr string, da Aggregator) error {
	log.Println("GRPC service running on ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewAggregatorGRPCService(da))

	return server.Serve(ln)
}

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
