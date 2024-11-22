package boot

import (
	"log"
	"net"

	pb "github.com/aparnasukesh/inter-communication/user_admin"
	"github.com/aparnasukesh/user-admin-super_admin-svc/config"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/admin"
	superadmin "github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/super-admin"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/user"
	"google.golang.org/grpc"
)

func NewGrpcServer(config config.Config, userGrpcHandler user.GrpcHandler, adminGrpcHandler admin.GrpcHandler, superAdminGrpcHandler superadmin.GrpcHandler) (func() error, error) {
	//lis, err := net.Listen("tcp", ":"+config.GrpcPort)
	lis, err := net.Listen("tcp", "0.0.0.0:"+config.GrpcPort)

	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userGrpcHandler)
	pb.RegisterAdminServiceServer(s, &adminGrpcHandler)
	pb.RegisterSuperAdminServiceServer(s, &superAdminGrpcHandler)
	// Assuming `server` is your gRPC server instance
	srv := func() error {
		log.Printf("gRPC server started on port %s", config.GrpcPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return err
		}
		return nil
	}
	return srv, nil
}
