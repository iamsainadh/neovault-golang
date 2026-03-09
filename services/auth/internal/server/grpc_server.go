package server

import (
	"net"

	pb "github.com/iamsainadh/neovault/auth/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer struct {
	Server *grpc.Server
	health *health.Server
}

func NewGRPCServerConn(handler pb.UserServiceServer) *GRPCServer {

	s := grpc.NewServer()

	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(s, healthServer)

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	pb.RegisterUserServiceServer(s, handler)

	return &GRPCServer{
		Server: s,
	}
}

func (g *GRPCServer) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return g.Server.Serve(lis)
}

func (g *GRPCServer) Stop() {
	g.health.SetServingStatus(
		"",
		healthpb.HealthCheckResponse_NOT_SERVING,
	)
	g.Server.GracefulStop()
}
