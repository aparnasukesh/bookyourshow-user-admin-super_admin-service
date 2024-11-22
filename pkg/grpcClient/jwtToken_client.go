package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/auth"
	"google.golang.org/grpc"
)

func NewJWTGrpcClint(port string) (pb.JWT_TokenServiceClient, error) {
	address := "auth-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewJWT_TokenServiceClient(conn), nil
}
