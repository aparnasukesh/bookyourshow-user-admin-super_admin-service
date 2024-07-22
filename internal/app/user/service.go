package user

import (
	"context"
	"errors"
	"fmt"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/utils"
)

type service struct {
	repo               UserRepository
	notificationClient notificationClient.EmailServiceClient
	authClient         authClient.JWT_TokenServiceClient
}

type Service interface {
	RegisterUser(ctx context.Context, user User) error
	ValidateUser(ctx context.Context, user User) error
	LoginUser(ctx context.Context, user User) (string, error)
}

func NewService(repo UserRepository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
	}
}

func (s *service) RegisterUser(ctx context.Context, user User) error {
	res, err := s.repo.GetUserByEmail(ctx, user.Email)
	if res != nil && err == nil {
		return errors.New("email already exist")
	}
	hashPass := utils.HashPassword(user.Password)
	user.Password = hashPass
	otp, err := utils.GenCaptchaCode()
	if err != nil {
		return err
	}
	user.Otp = otp
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}
	_, err = s.notificationClient.SendEmail(ctx, &notificationClient.EmailRequest{
		Email:       user.Email,
		Otp:         otp,
		BodyMessage: "",
	})
	if err != nil {
		return fmt.Errorf("failed to send email")
	}
	return nil

}

func (s *service) ValidateUser(ctx context.Context, user User) error {
	enterdOtp := user.Otp
	res, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return errors.New("invalid email address")
	}

	if user.Email != res.Email || enterdOtp != res.Otp {
		err = s.repo.DeleteUserByEmail(ctx, user)
		if err != nil {
			return err
		}
		return errors.New("invalid otp")
	}

	if err := s.repo.UserApproval(ctx, user.Email); err != nil {
		return err
	}
	userRoles := createUserRoles(res.ID, USER_ROLE)
	if err := s.repo.CreateUserRoles(ctx, userRoles); err != nil {
		return err
	}

	return nil
}

func (s *service) LoginUser(ctx context.Context, user User) (string, error) {
	res, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", nil
	}

	if !res.IsVerified {
		return "", errors.New("invalid user")
	}
	isVerified := utils.VerifyPassword(user.Password, res.Password)
	if user.Email != res.Email || !isVerified {
		return "", errors.New("incorrect password")

	}
	exist, err := s.repo.CheckUserRole(ctx, res.ID)
	if !exist && err != nil {
		return "", err
	}
	response, err := s.authClient.GenerateJWt(ctx, &authClient.GenerateRequest{
		RoleId: USER_ROLE,
		Email:  res.Email,
		UserId: int32(res.ID),
	})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}
