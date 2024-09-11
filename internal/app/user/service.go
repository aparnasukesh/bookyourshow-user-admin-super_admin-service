package user

import (
	"context"
	"errors"
	"fmt"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/utils"
	"gorm.io/gorm"
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
	GetProfileDetails(ctx context.Context, userId int) (*UserProfileDetails, error)
	UpdateUserProfile(ctx context.Context, id int, admin UserProfileDetails) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, data ResetPassword) error
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
	_, err = s.notificationClient.SendEmail(ctx, &notificationClient.EmailRequest{
		Email:       user.Email,
		Otp:         otp,
		BodyMessage: "",
	})
	if err != nil {
		return fmt.Errorf("failed to send email")
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
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
			return errors.New("email id not exist")
		}
		return errors.New("invalid otp or ivalid email address")
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
		return "", err
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

func (s *service) GetProfileDetails(ctx context.Context, userId int) (*UserProfileDetails, error) {
	profileDetails, err := s.repo.GetProfileDetails(ctx, userId)
	if err != nil {
		return nil, err
	}
	return profileDetails, nil
}

func (s *service) UpdateUserProfile(ctx context.Context, id int, user UserProfileDetails) error {
	data, err := s.repo.GetUserDetails(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user.Username != "" {
		data.Username = user.Username
	}
	if user.PhoneNumber != "" {
		data.PhoneNumber = user.PhoneNumber
	}
	if user.FirstName != "" {
		data.FirstName = user.FirstName
	}
	if user.LastName != "" {
		data.LastName = user.LastName
	}
	if user.DateOfBirth != "" {
		data.DateOfBirth = user.DateOfBirth
	}
	if user.Gender != "" {
		data.Gender = user.Gender
	}
	if err := s.repo.UpdateUserProfile(ctx, data); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("no user found with this email id")
	}

	otp, err := utils.GenCaptchaCode()
	if err != nil {
		return err
	}
	user.Otp = otp
	_, err = s.notificationClient.SendResetPassWordEmail(ctx, &notificationClient.EmailRequest{
		Email:       user.Email,
		Otp:         otp,
		BodyMessage: "",
	})
	if err != nil {
		return fmt.Errorf("failed to send the password reset email")
	}
	err = s.repo.UpdateOtp(ctx, email, otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, data ResetPassword) error {
	user, err := s.repo.GetUserByEmail(ctx, data.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("no user found with this email id")
	}

	password := utils.HashPassword(data.NewPassword)
	if data.Email != user.Email || data.Otp != user.Otp {
		return errors.New("invalid otp")
	}

	if password == user.Password {
		return errors.New("try another password")
	}
	err = s.repo.ResetPassword(ctx, user.Email, password)
	if err != nil {
		return err
	}
	return nil
}
