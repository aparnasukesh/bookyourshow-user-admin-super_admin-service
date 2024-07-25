package user

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type GrpcHandler struct {
	svc Service
	user_admin.UnimplementedUserServiceServer
}

func NewGrpcHandler(service Service) GrpcHandler {
	return GrpcHandler{
		svc: service,
	}
}
func (h *GrpcHandler) RegisterUser(ctx context.Context, req *user_admin.RegisterUserRequest) (*user_admin.RegisterUserResponse, error) {

	userData := User{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.Phone,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
	}
	if err := h.svc.RegisterUser(ctx, userData); err != nil {
		return nil, err
	}
	return &user_admin.RegisterUserResponse{
		Status:     "Redirect: http://localhost:8080/gateway/register/validate",
		StatusCode: 200,
	}, nil
}

func (h *GrpcHandler) ValidateUser(ctx context.Context, req *user_admin.ValidateUserRequest) (*user_admin.ValidateUserResponse, error) {
	userData := User{
		Email: req.Email,
		Otp:   req.Otp,
	}
	if err := h.svc.ValidateUser(ctx, userData); err != nil {
		return nil, err
	}
	return &user_admin.ValidateUserResponse{
		Status:     "sign-up successfull",
		StatusCode: 200,
	}, nil
}

func (h *GrpcHandler) LoginUser(ctx context.Context, req *user_admin.LoginUserRequest) (*user_admin.LoginUserResponse, error) {
	userData := User{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := h.svc.LoginUser(ctx, userData)
	if err != nil {
		return nil, err
	}
	return &user_admin.LoginUserResponse{
		Status:     "login successfull",
		StatusCode: 200,
		Token:      token,
	}, nil
}

func (h *GrpcHandler) GetUserProfile(ctx context.Context, req *user_admin.GetProfileRequest) (*user_admin.GetProfileResponse, error) {
	profileDetails, err := h.svc.GetProfileDetails(ctx, int(req.UserId))
	if err != nil {
		return &user_admin.GetProfileResponse{
			Status:         "Failed to retrieve user profile",
			StatusCode:     404,
			ProfileDetails: nil,
		}, err
	}
	return &user_admin.GetProfileResponse{
		Status:     "User profile retrieved successfully",
		StatusCode: 200,
		ProfileDetails: &user_admin.UpdateUserProfileRequest{
			Username:    profileDetails.Username,
			Phone:       profileDetails.PhoneNumber,
			Email:       profileDetails.Email,
			FirstName:   profileDetails.FirstName,
			LastName:    profileDetails.LastName,
			Gender:      profileDetails.Gender,
			DateOfBirth: profileDetails.DateOfBirth,
		},
	}, nil
}