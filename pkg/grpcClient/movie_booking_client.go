package grpcclient

import (
	"log"

	pb "github.com/aparnasukesh/inter-communication/movie_booking"
	"google.golang.org/grpc"
)

func NewMovieBookingGrpcClint(port string) (pb.MovieServiceClient, error) {
	address := "movies-booking-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewMovieServiceClient(conn), nil
}

func NewTheaterGrpcClient(port string) (pb.TheatreServiceClient, error) {
	address := "movies-booking-svc.default.svc.cluster.local:" + port
	serviceConfig := `{"loadBalancingPolicy": "round_robin"}`
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(serviceConfig))
	if err != nil {
		log.Printf("Failed to connect to gRPC service: %v", err)
		return nil, err
	}
	return pb.NewTheatreServiceClient(conn), nil
}
