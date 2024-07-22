package superadmin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type GrpcHandler struct {
	svc Service
	user_admin.UnimplementedSuperAdminServiceServer
}

func NewGrpcHandler(service Service) GrpcHandler {
	return GrpcHandler{
		svc: service,
	}
}

func (h *GrpcHandler) LoginSuperAdmin(ctx context.Context, req *user_admin.LoginSuperAdminRequest) (*user_admin.LoginSuperAdminResponse, error) {
	token, err := h.svc.LoginSuperAdmin(ctx, User{
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.LoginSuperAdminResponse{
		Status:     "login successfull",
		StatusCode: 200,
		Token:      token,
	}, nil
}

//func (h *GrpcHandler) AdminApproval(context.Context, *user_admin.AdminApprovalRequest) (*user_admin.AdminApprovalResponse, error)
