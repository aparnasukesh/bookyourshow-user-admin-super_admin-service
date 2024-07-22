package admin

import (
	"context"
	"errors"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/utils"
	"gorm.io/gorm"
)

type service struct {
	repo               Repository
	notificationClient notificationClient.EmailServiceClient
	authClient         authClient.JWT_TokenServiceClient
}

type Service interface {
	RegisterAdmin(ctx context.Context, user User) error
	LoginAdmin(ctx context.Context, user User) (string, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
	}
}

func (s *service) RegisterAdmin(ctx context.Context, user User) error {
	existingUser, err := s.repo.GetAdminByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if existingUser != nil {
		isVerified, err := s.repo.CheckIsVerified(ctx, existingUser.Email)
		if err != nil {
			return err
		}
		if !isVerified {
			return errors.New("user not verified")
		}

		exists, err := s.repo.CheckAdminExist(ctx, existingUser.Email)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if exists {
			status, err := s.repo.CheckAdminStatus(ctx, existingUser.Email)
			if err != nil {
				return err
			}

			switch status {
			case "pending":
				return errors.New("request pending already")
			case "approved":
				return errors.New("already an admin")
			case "rejected":
				return errors.New("already rejected")
			}
		}

		if !exists {
			if err := s.repo.CreateAdminStatus(ctx, &AdminStatus{
				Status: "pending",
				Email:  existingUser.Email,
			}); err != nil {
				return err
			}
		}

		return nil
	}

	if err := s.repo.CreateAdmin(ctx, user); err != nil {
		return err
	}
	if err := s.repo.CreateAdminStatus(ctx, &AdminStatus{
		Status: "pending",
		Email:  user.Email,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) LoginAdmin(ctx context.Context, user User) (string, error) {
	res, err := s.repo.GetAdminByEmail(ctx, user.Email)
	if err != nil {
		return "", nil
	}

	if !res.IsVerified {
		return "", errors.New("")
	}
	isVerified := utils.VerifyPassword(user.Password, res.Password)
	if user.Email != res.Email || !isVerified {
		return "", errors.New("incorrect password")

	}
	exist, err := s.repo.CheckAdminRole(ctx, res.ID)
	if !exist && err != nil {
		return "", errors.New("no admin exist")
	}
	response, err := s.authClient.GenerateJWt(ctx, &authClient.GenerateRequest{
		RoleId: ADMIN_ROLE,
		Email:  res.Email,
		UserId: int32(res.ID),
	})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}
