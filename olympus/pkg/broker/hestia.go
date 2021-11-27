package broker

import (
	"context"
	"math/rand"

	pb "github.com/project-auxo/auxo/olympus/proto/hestia"
)

type hestiaFrontendServer struct {
	pb.UnimplementedHestiaFrontendServiceServer
	broker *Broker // The currently running broker
}

// GetNumberOfAgents returns the number of agents that are currently connected
// to Olympus.
func (s *hestiaFrontendServer) GetNumberOfAgents(
	ctx context.Context, req *pb.GetNumberOfAgentsReq) (*pb.GetNumberOfAgentsRep, error) {
	// TODO(bellabah): Should make use of information within the broker.
	numAgents := rand.Intn(10)
	return &pb.GetNumberOfAgentsRep{Number: int32(numAgents)}, nil
}
