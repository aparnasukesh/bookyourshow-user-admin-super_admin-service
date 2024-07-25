package superadmin

import (
	"context"
	"errors"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/movie_booking"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
)

type service struct {
	repo               Repository
	notificationClient notificationClient.EmailServiceClient
	authClient         authClient.JWT_TokenServiceClient
	movieBookingClient movie_booking.MovieServiceClient
}

type Service interface {
	LoginSuperAdmin(ctx context.Context, user Admin) (string, error)
	ListAdminRequests(ctx context.Context) ([]AdminRequests, error)
	AdminApproval(ctx context.Context, email string, isVerified bool) error
	RegisterMovie(ctx context.Context, movie Movie) (uint, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient, movieBookingClient movie_booking.MovieServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
		movieBookingClient: movieBookingClient,
	}
}

func (s *service) LoginSuperAdmin(ctx context.Context, user Admin) (string, error) {
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

func (s *service) ListAdminRequests(ctx context.Context) ([]AdminRequests, error) {
	adminLists, err := s.repo.ListAdiminRequests(ctx)
	if err != nil {
		return nil, err
	}
	adminRequests := buildAdminRequestsList(adminLists)
	return adminRequests, nil
}

func (s *service) AdminApproval(ctx context.Context, email string, isVerified bool) error {
	if err := s.repo.AdminApproval(ctx, email, isVerified); err != nil {
		return err
	}
	res, err := s.repo.GetAdminByEmail(ctx, email)
	if err != nil {
		return errors.New("invalid email address")
	}
	adminRole := AdminRole{
		AdminID: res.ID,
		RoleID:  ADMIN_ROLE,
		Admin:   Admin{},
		Role:    Role{},
	}
	if isVerified {
		if err := s.repo.UpdateIsVerified(ctx, email); err != nil {
			return err
		}
		err := s.repo.CreateAdminRoles(ctx, adminRole)
		if err != nil {
			return err
		}
		return nil
	}
	if !isVerified {
		if err := s.repo.DeleteAdminByEmail(ctx, *res); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) RegisterMovie(ctx context.Context, movie Movie) (uint, error) {
	response, err := s.movieBookingClient.RegisterMovie(ctx, &movie_booking.RegisterMovieRequest{
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		Genre:       movie.Genre,
		ReleaseDate: movie.ReleaseDate,
		Rating:      float32(movie.Rating),
		Language:    movie.Language,
	})
	if err != nil {
		return 0, err
	}
	return uint(response.MovieId), nil
}
