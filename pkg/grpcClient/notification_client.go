package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/notification"
	"google.golang.org/grpc"
)

// func NewNotificationGrpcClint(port string) (pb.EmailServiceClient, error) {
// 	address := "notification-svc.default.svc.cluster.local:" + port
// 	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
// 	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
// 	if err != nil {
// 		log.Printf("Failed to connect to gRPC service: %v", err)
// 		return nil, err
// 	}
// 	return pb.NewEmailServiceClient(conn), nil
// }

func NewNotificationGrpcClient(port string) (pb.EmailServiceClient, error) {
	address := "notification-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`

	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(serviceConfig),
		grpc.WithBlock(), // Ensures the call blocks until the connection is established
	)
	if err != nil {
		log.Printf("Failed to connect to gRPC service at %s: %v", address, err)
		return nil, err
	}

	client := pb.NewEmailServiceClient(conn)
	log.Printf("Connected to gRPC service at %s", address)
	return client, nil
}
