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
	movieBooking       movie_booking.MovieServiceClient
	theaterClient      movie_booking.TheatreServiceClient
}

type Service interface {
	LoginSuperAdmin(ctx context.Context, user Admin) (string, error)
	ListAdminRequests(ctx context.Context) ([]AdminRequests, error)
	AdminApproval(ctx context.Context, email string, isVerified bool) error
	RegisterMovie(ctx context.Context, movie Movie) (uint, error)
	UpdateMovie(ctx context.Context, movie Movie, movieId int) error
	ListMovies(ctx context.Context) ([]Movie, error)
	GetMovieDetails(ctx context.Context, movieId int) (*Movie, error)
	DeleteMovie(ctx context.Context, movieId int) error
	AddTheaterType(ctx context.Context, theaterType TheaterType) error
	DeleteTheaterTypeByID(ctx context.Context, id int) error
	DeleteTheaterTypeByName(ctx context.Context, name string) error
	GetTheaterTypeByID(ctx context.Context, id int) (*TheaterType, error)
	GetTheaterTypeByName(ctx context.Context, name string) (*TheaterType, error)
	UpdateTheaterType(ctx context.Context, id int, theaterType TheaterType) error
	ListTheaterTypes(ctx context.Context) ([]TheaterType, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient, movieBookingClient movie_booking.MovieServiceClient, theaterClient movie_booking.TheatreServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
		movieBooking:       movieBookingClient,
		theaterClient:      theaterClient,
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
	response, err := s.movieBooking.RegisterMovie(ctx, &movie_booking.RegisterMovieRequest{
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

func (s *service) UpdateMovie(ctx context.Context, movie Movie, movieId int) error {
	_, err := s.movieBooking.UpdateMovie(ctx, &movie_booking.UpdateMovieRequest{
		MovieId:     uint32(movieId),
		Title:       movie.Title,
		Description: movie.Description,
		Duration:    int32(movie.Duration),
		Genre:       movie.Genre,
		ReleaseDate: movie.ReleaseDate,
		Rating:      float32(movie.Rating),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteMovie(ctx context.Context, movieId int) error {
	_, err := s.movieBooking.DeleteMovie(ctx, &movie_booking.DeleteMovieRequest{
		MovieId: uint32(movieId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetMovieDetails(ctx context.Context, movieId int) (*Movie, error) {
	response, err := s.movieBooking.GetMovieDetails(ctx, &movie_booking.GetMovieDetailsRequest{
		MovieId: uint32(movieId),
	})
	if err != nil {
		return nil, err
	}
	if response.Movie == nil {
		return nil, errors.New("movie details not found")
	}
	movie := &Movie{
		Title:       response.Movie.Title,
		Description: response.Movie.Description,
		Duration:    int(response.Movie.Duration),
		Genre:       response.Movie.Genre,
		ReleaseDate: response.Movie.ReleaseDate,
		Rating:      float64(response.Movie.Rating),
	}

	return movie, nil
}

func (s *service) ListMovies(ctx context.Context) ([]Movie, error) {
	response, err := s.movieBooking.ListMovies(ctx, &movie_booking.ListMoviesRequest{})
	if err != nil {
		return nil, err
	}
	var movies []Movie
	for _, m := range response.Movies {
		movie := Movie{
			Title:       m.Title,
			Description: m.Description,
			Duration:    int(m.Duration),
			Genre:       m.Genre,
			ReleaseDate: m.ReleaseDate,
			Rating:      float64(m.Rating),
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *service) AddTheaterType(ctx context.Context, theaterType TheaterType) error {
	_, err := s.theaterClient.AddTheaterType(ctx, &movie_booking.AddTheaterTypeRequest{
		TheaterTypeName: theaterType.TheaterTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterTypeByID(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteTheaterTypeByID(ctx, &movie_booking.DeleteTheaterTypeRequest{
		TheaterTypeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterTypeByName(ctx context.Context, name string) error {
	_, err := s.theaterClient.DeleteTheaterTypeByName(ctx, &movie_booking.DeleteTheaterTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterTypeByID(ctx context.Context, id int) (*TheaterType, error) {
	response, err := s.theaterClient.GetTheaterTypeByID(ctx, &movie_booking.GetTheaterTypeByIDRequest{
		TheaterTypeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterType{
		ID:              int(response.TheaterType.Id),
		TheaterTypeName: response.TheaterType.TheaterTypeName,
	}, nil
}

func (s *service) GetTheaterTypeByName(ctx context.Context, name string) (*TheaterType, error) {
	response, err := s.theaterClient.GetTheaterTypeByName(ctx, &movie_booking.GetTheaterTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &TheaterType{
		ID:              int(response.TheaterType.Id),
		TheaterTypeName: response.TheaterType.TheaterTypeName,
	}, nil
}

func (s *service) ListTheaterTypes(ctx context.Context) ([]TheaterType, error) {
	response, err := s.theaterClient.ListTheaterTypes(ctx, &movie_booking.ListTheaterTypesRequest{})
	if err != nil {
		return nil, err
	}
	theaterTypes := []TheaterType{}

	for _, res := range response.TheaterTypes {
		theaterType := TheaterType{
			ID:              int(res.Id),
			TheaterTypeName: res.TheaterTypeName,
		}
		theaterTypes = append(theaterTypes, theaterType)
	}
	return theaterTypes, nil
}

func (s *service) UpdateTheaterType(ctx context.Context, id int, theaterType TheaterType) error {
	_, err := s.theaterClient.UpdateTheaterType(ctx, &movie_booking.UpdateTheaterTypeRequest{
		Id:              int32(id),
		TheaterTypeName: theaterType.TheaterTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}
