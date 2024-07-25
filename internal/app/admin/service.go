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
	RegisterAdmin(ctx context.Context, admin Admin) error
	LoginAdmin(ctx context.Context, admin Admin) (string, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
	}
}

func (s *service) RegisterAdmin(ctx context.Context, admin Admin) error {
	existingAdmin, err := s.repo.GetAdminByEmail(ctx, admin.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	if existingAdmin != nil && err == nil {
		return errors.New("email already exist")
	}
	hashPass := utils.HashPassword(admin.Password)
	admin.Password = hashPass
	if existingAdmin != nil {
		exists, err := s.repo.CheckAdminExist(ctx, existingAdmin.Email)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if exists {
			status, err := s.repo.CheckAdminStatus(ctx, existingAdmin.Email)
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
				Email:  existingAdmin.Email,
			}); err != nil {
				return err
			}
		}
		return nil
	}

	if err := s.repo.CreateAdmin(ctx, admin); err != nil {
		return err
	}
	if err := s.repo.CreateAdminStatus(ctx, &AdminStatus{
		Status: "pending",
		Email:  admin.Email,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) LoginAdmin(ctx context.Context, admin Admin) (string, error) {
	res, err := s.repo.GetAdminByEmail(ctx, admin.Email)
	if err != nil {
		return "", nil
	}

	if !res.IsVerified {
		return "", errors.New("")
	}
	isVerified := utils.VerifyPassword(admin.Password, res.Password)
	if admin.Email != res.Email || !isVerified {
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
