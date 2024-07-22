package superadmin

import (
	"context"
	"errors"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
)

type service struct {
	repo               Repository
	notificationClient notificationClient.EmailServiceClient
	authClient         authClient.JWT_TokenServiceClient
}

type Service interface {
	LoginSuperAdmin(ctx context.Context, user User) (string, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
	}
}

func (s *service) LoginSuperAdmin(ctx context.Context, user User) (string, error) {
	res, err := s.repo.GetSuperAdminByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if user.Email != res.Email || user.Password != res.Password {
		return "", errors.New("incorrect password")

	}
	response, err := s.authClient.GenerateJWt(ctx, &authClient.GenerateRequest{
		RoleId: SUPER_ADMIN_ROLE,
		Email:  res.Email,
		UserId: int32(res.ID),
	})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}
