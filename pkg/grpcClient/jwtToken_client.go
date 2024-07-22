package grpcclient

import (
	pb "github.com/aparnasukesh/inter-communication/auth"
	"google.golang.org/grpc"
)

func NewJWTGrpcClint(port string) (pb.JWT_TokenServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewJWT_TokenServiceClient(conn), nil
}
