package admin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type GrpcHandler struct {
	svc Service
	user_admin.UnimplementedAdminServiceServer
}

func NewGrpcHandler(service Service) GrpcHandler {
	return GrpcHandler{
		svc: service,
	}
}

func (h *GrpcHandler) RegisterAdmin(ctx context.Context, req *user_admin.RegisterAdminRequest) (*user_admin.RegisterAdminResponse, error) {

	userData := Admin{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.Phone,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
	}
	if err := h.svc.RegisterAdmin(ctx, userData); err != nil {
		return nil, err
	}
	return &user_admin.RegisterAdminResponse{
		Status:     "request pending",
		StatusCode: 200,
	}, nil
}

func (h *GrpcHandler) LoginAdmin(ctx context.Context, req *user_admin.LoginAdminRequest) (*user_admin.LoginAdminResponse, error) {
	userData := Admin{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := h.svc.LoginAdmin(ctx, userData)
	if err != nil {
		return nil, err
	}
	return &user_admin.LoginAdminResponse{
		Status:     "login successfull",
		StatusCode: 200,
		Token:      token,
	}, nil
}
