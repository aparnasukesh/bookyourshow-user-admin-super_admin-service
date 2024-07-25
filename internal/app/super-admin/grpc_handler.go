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
	token, err := h.svc.LoginSuperAdmin(ctx, Admin{
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
func (h *GrpcHandler) ListAdminRequests(ctx context.Context, req *user_admin.ListAdminRequestsRequest) (*user_admin.ListAdminRequestsResponse, error) {
	adminLists, err := h.svc.ListAdminRequests(ctx)
	if err != nil {
		return nil, err
	}

	emails := make([]*user_admin.Email, len(adminLists))
	for i, admin := range adminLists {
		emails[i] = &user_admin.Email{
			Email: admin.Email,
		}
	}

	return &user_admin.ListAdminRequestsResponse{
		Email: emails,
	}, nil
}

func (h *GrpcHandler) AdminApproval(ctx context.Context, req *user_admin.AdminApprovalRequest) (*user_admin.AdminApprovalResponse, error) {
	err := h.svc.AdminApproval(ctx, req.Email, req.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user_admin.AdminApprovalResponse{}, nil

}

func (h *GrpcHandler) RegisterMovie(ctx context.Context, req *user_admin.RegisterMovieRequest) (*user_admin.RegisterMovieResponse, error) {
	movieId, err := h.svc.RegisterMovie(ctx, Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    int(req.Duration),
		Genre:       req.Genre,
		ReleaseDate: req.ReleaseDate,
		Rating:      float64(req.Rating),
		Language:    req.Language,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.RegisterMovieResponse{
		MovieId: uint32(movieId),
		Message: "create movie successfull",
	}, nil
}
