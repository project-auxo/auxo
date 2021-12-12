package grpc

import (
	"context"

	pb "github.com/project-auxo/auxo/oracle/proto"
	"github.com/project-auxo/auxo/oracle/registry"
)

type OracleBackendServer struct {
	pb.UnimplementedOracleBackendServiceServer
}

func (s *OracleBackendServer) CheckServiceExists(
	ctx context.Context, req *pb.CheckServiceExistsReq) (*pb.CheckServiceExistsRep, error) {
	rep := &pb.CheckServiceExistsRep{}
	rep.Exists = registry.ServiceExists(req.ServiceName)
	return rep, nil
}
