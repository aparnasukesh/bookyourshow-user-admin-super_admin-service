package di

import (
	"log"

	"github.com/aparnasukesh/user-admin-super_admin-svc/config"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/admin"
	superadmin "github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/super-admin"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/app/user"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/boot"
	grpcClient "github.com/aparnasukesh/user-admin-super_admin-svc/pkg/grpcClient"
	"github.com/aparnasukesh/user-admin-super_admin-svc/pkg/sql"
)

func InitResources(cfg config.Config) (func() error, error) {

	// Db initialization
	db, err := sql.NewSql(cfg)
	if err != nil {
		log.Fatal(err)
	}
	notificationClient, err := grpcClient.NewNotificationGrpcClint(cfg.GrpcNotificationPort)
	if err != nil {
		return nil, err
	}

	authClient, err := grpcClient.NewJWTGrpcClint(cfg.GrpcAuthPort)
	if err != nil {
		return nil, err
	}

	movieBookingClient, err := grpcClient.NewMovieBookingGrpcClint(cfg.GrpcMovieBookingPort)
	if err != nil {
		return nil, err
	}
	// // Admin Module Initialization
	adminRepo := admin.NewRepository(db)
	adminService := admin.NewService(adminRepo, notificationClient, authClient)
	adminGrpcHandler := admin.NewGrpcHandler(adminService)

	// User Module initialization
	repo := user.NewRepository(db)
	service := user.NewService(repo, notificationClient, authClient)
	userGrpcHandler := user.NewGrpcHandler(service)

	// SuperAdmin Module initialization
	superAdminRepo := superadmin.NewRepository(db)
	superAdminService := superadmin.NewService(superAdminRepo, notificationClient, authClient, movieBookingClient)
	superAdminGrpcHandler := superadmin.NewGrpcHandler(superAdminService)

	// Server initialization
	server, err := boot.NewGrpcServer(cfg, userGrpcHandler, adminGrpcHandler, superAdminGrpcHandler)
	if err != nil {
		log.Fatal(err)
	}
	return server, nil
}
