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
	//movies
	RegisterMovie(ctx context.Context, movie Movie) (uint, error)
	UpdateMovie(ctx context.Context, movie Movie, movieId int) error
	ListMovies(ctx context.Context) ([]Movie, error)
	GetMovieDetails(ctx context.Context, movieId int) (*Movie, error)
	DeleteMovie(ctx context.Context, movieId int) error
	//theater types
	AddTheaterType(ctx context.Context, theaterType TheaterType) error
	DeleteTheaterTypeByID(ctx context.Context, id int) error
	DeleteTheaterTypeByName(ctx context.Context, name string) error
	GetTheaterTypeByID(ctx context.Context, id int) (*TheaterType, error)
	GetTheaterTypeByName(ctx context.Context, name string) (*TheaterType, error)
	UpdateTheaterType(ctx context.Context, id int, theaterType TheaterType) error
	ListTheaterTypes(ctx context.Context) ([]TheaterType, error)
	//screen types
	AddScreenType(ctx context.Context, data ScreenType) error
	DeleteScreenTypeById(ctx context.Context, id int) error
	DeleteScreenTypeByName(ctx context.Context, screenName string) error
	GetScreenTypeByID(ctx context.Context, id int) (*ScreenType, error)
	GetScreenTypeByName(ctx context.Context, name string) (*ScreenType, error)
	UpdateScreenType(ctx context.Context, id int, screenType ScreenType) error
	ListScreenTypes(ctx context.Context) ([]ScreenType, error)
	// Seat category
	AddSeatCategory(ctx context.Context, seatCategory SeatCategory) error
	DeleteSeatCategoryByID(ctx context.Context, id int) error
	DeleteSeatCategoryByName(ctx context.Context, name string) error
	GetSeatCategoryByID(ctx context.Context, id int) (*SeatCategory, error)
	GetSeatCategoryByName(ctx context.Context, name string) (*SeatCategory, error)
	UpdateSeatCategory(ctx context.Context, id int, seatCategory SeatCategory) error
	ListSeatCategories(ctx context.Context) ([]SeatCategory, error)
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
		return errors.New("no admin request found with the provided email ID")
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
		Language:    movie.Language,
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
		Language:    response.Movie.Language,
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
			Language:    m.Language,
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// theater type
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

// screen types
func (s *service) AddScreenType(ctx context.Context, screenType ScreenType) error {
	_, err := s.theaterClient.AddScreenType(ctx, &movie_booking.AddScreenTypeRequest{
		ScreenTypeName: screenType.ScreenTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteScreenTypeById(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteScreenTypeByID(ctx, &movie_booking.DeleteScreenTypeRequest{
		ScreenTypeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteScreenTypeByName(ctx context.Context, name string) error {
	_, err := s.theaterClient.DeleteScreenTypeByName(ctx, &movie_booking.DeleteScreenTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetScreenTypeByID(ctx context.Context, id int) (*ScreenType, error) {
	response, err := s.theaterClient.GetScreenTypeByID(ctx, &movie_booking.GetScreenTypeByIDRequest{
		ScreenTypeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &ScreenType{
		ID:             int(response.ScreenType.Id),
		ScreenTypeName: response.ScreenType.ScreenTypeName,
	}, nil
}

func (s *service) GetScreenTypeByName(ctx context.Context, name string) (*ScreenType, error) {
	response, err := s.theaterClient.GetScreenTypeByName(ctx, &movie_booking.GetScreenTypeByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &ScreenType{
		ID:             int(response.ScreenType.Id),
		ScreenTypeName: response.ScreenType.ScreenTypeName,
	}, nil
}

func (s *service) ListScreenTypes(ctx context.Context) ([]ScreenType, error) {
	response, err := s.theaterClient.ListScreenTypes(ctx, &movie_booking.ListScreenTypesRequest{})
	if err != nil {
		return nil, err
	}
	screenTypes := []ScreenType{}

	for _, res := range response.ScreenTypes {
		screenType := ScreenType{
			ID:             int(res.Id),
			ScreenTypeName: res.ScreenTypeName,
		}
		screenTypes = append(screenTypes, screenType)
	}
	return screenTypes, nil
}

func (s *service) UpdateScreenType(ctx context.Context, id int, screenType ScreenType) error {
	_, err := s.theaterClient.UpdateScreenType(ctx, &movie_booking.UpdateScreenTypeRequest{
		Id:             int32(id),
		ScreenTypeName: screenType.ScreenTypeName,
	})
	if err != nil {
		return err
	}
	return nil
}

// seat category
func (s *service) AddSeatCategory(ctx context.Context, seatCategory SeatCategory) error {
	_, err := s.theaterClient.AddSeatCategory(ctx, &movie_booking.AddSeatCategoryRequest{
		SeatCategory: &movie_booking.SeatCategory{
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteSeatCategoryByID(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteSeatCategoryByID(ctx, &movie_booking.DeleteSeatCategoryRequest{
		SeatCategoryId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteSeatCategoryByName(ctx context.Context, name string) error {
	_, err := s.theaterClient.DeleteSeatCategoryByName(ctx, &movie_booking.DeleteSeatCategoryByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetSeatCategoryByID(ctx context.Context, id int) (*SeatCategory, error) {
	response, err := s.theaterClient.GetSeatCategoryByID(ctx, &movie_booking.GetSeatCategoryByIDRequest{
		SeatCategoryId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &SeatCategory{
		ID:               int(response.SeatCategory.Id),
		SeatCategoryName: response.SeatCategory.SeatCategoryName,
	}, nil
}

func (s *service) GetSeatCategoryByName(ctx context.Context, name string) (*SeatCategory, error) {
	response, err := s.theaterClient.GetSeatCategoryByName(ctx, &movie_booking.GetSeatCategoryByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &SeatCategory{
		ID:               int(response.SeatCategory.Id),
		SeatCategoryName: response.SeatCategory.SeatCategoryName,
	}, nil
}

func (s *service) ListSeatCategories(ctx context.Context) ([]SeatCategory, error) {
	response, err := s.theaterClient.ListSeatCategories(ctx, &movie_booking.ListSeatCategoriesRequest{})
	if err != nil {
		return nil, err
	}
	seatCategories := []SeatCategory{}

	for _, res := range response.SeatCategories {
		seatCategory := SeatCategory{
			ID:               int(res.Id),
			SeatCategoryName: res.SeatCategoryName,
		}
		seatCategories = append(seatCategories, seatCategory)
	}
	return seatCategories, nil
}

func (s *service) UpdateSeatCategory(ctx context.Context, id int, seatCategory SeatCategory) error {
	_, err := s.theaterClient.UpdateSeatCategory(ctx, &movie_booking.UpdateSeatCategoryRequest{
		Id: int32(id),
		SeatCategory: &movie_booking.SeatCategory{
			Id:               int32(id),
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
