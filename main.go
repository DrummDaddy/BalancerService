package main

import (
	"log"
	"net"

	"BalancerService/config"
	"BalancerService/internal/handlers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "BalancerService/proto/service"
)

func main() {
	config := config.LoadConfig()
	grpcOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(10000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * 60,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             2 * 60,
			PermitWithoutStream: true,
		}),
	}
	server := grpc.NewServer(grpcOptions...)
	handler := handlers.NewBalancerHandler(config)
	pb.RegisterBalancerServiceServer(server, handler)

	listener, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("Failed to listen on port 443: %v", err)
	}

	log.Println("Balancer service is running on port 443")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
