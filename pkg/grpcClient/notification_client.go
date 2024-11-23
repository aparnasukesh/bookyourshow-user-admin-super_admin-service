package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/notification"
	"google.golang.org/grpc"
)

func NewNotificationGrpcClient(port string) (pb.EmailServiceClient, error) {
	//address := "notification-svc.default.svc.cluster.local:5051"
	//serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewEmailServiceClient(conn), nil
}
