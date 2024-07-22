package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/notification"
	"google.golang.org/grpc"
)

func NewNotificationGrpcClint(port string) (pb.EmailServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewEmailServiceClient(conn), nil
}
