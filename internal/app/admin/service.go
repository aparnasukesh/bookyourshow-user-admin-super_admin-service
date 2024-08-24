package admin

import (
	"context"
	"errors"
	"time"

	authClient "github.com/aparnasukesh/inter-communication/auth"
	"github.com/aparnasukesh/inter-communication/movie_booking"
	notificationClient "github.com/aparnasukesh/inter-communication/notification"
	"github.com/aparnasukesh/user-admin-super_admin-svc/internal/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type service struct {
	repo               Repository
	notificationClient notificationClient.EmailServiceClient
	authClient         authClient.JWT_TokenServiceClient
	movieBooking       movie_booking.MovieServiceClient
	theaterClient      movie_booking.TheatreServiceClient
}

type Service interface {
	RegisterAdmin(ctx context.Context, admin Admin) error
	LoginAdmin(ctx context.Context, admin Admin) (string, error)
	//Theater
	AddTheater(ctx context.Context, theater *Theater) error
	DeleteTheaterByID(ctx context.Context, id int) error
	DeleteTheaterByName(ctx context.Context, name string) error
	GetTheaterByID(ctx context.Context, id int) (*Theater, error)
	GetTheaterByName(ctx context.Context, name string) ([]Theater, error)
	UpdateTheater(ctx context.Context, id int, theater Theater) error
	ListTheaters(ctx context.Context) ([]Theater, error)
	//Seat categories
	ListSeatCategories(ctx context.Context) ([]SeatCategory, error)
	//Movies
	ListMovies(ctx context.Context) ([]Movie, error)
	//Theater-type
	ListTheaterTypes(ctx context.Context) ([]TheaterType, error)
	//Sreen type
	ListScreenTypes(ctx context.Context) ([]ScreenType, error)
	//Theater screen
	AddTheaterScreen(ctx context.Context, theaterScreen TheaterScreen) error
	DeleteTheaterScreenByID(ctx context.Context, id int) error
	DeleteTheaterScreenByNumber(ctx context.Context, theaterID int, screenNumber int) error
	GetTheaterScreenByID(ctx context.Context, id int) (*TheaterScreen, error)
	GetTheaterScreenByNumber(ctx context.Context, theaterID int, screenNumber int) (*TheaterScreen, error)
	UpdateTheaterScreen(ctx context.Context, id int, theaterScreen TheaterScreen) error
	ListTheaterScreens(ctx context.Context, theaterId int) ([]TheaterScreen, error)
	//Show time
	AddShowtime(ctx context.Context, showtime Showtime) error
	DeleteShowtimeByID(ctx context.Context, id int) error
	DeleteShowtimeByDetails(ctx context.Context, movieID int, screenID int, showDate time.Time, showTime time.Time) error
	GetShowtimeByID(ctx context.Context, id int) (*Showtime, error)
	GetShowtimeByDetails(ctx context.Context, movieID int, screenID int, showDate time.Time, showTime time.Time) (*Showtime, error)
	UpdateShowtime(ctx context.Context, id int, showtime Showtime) error
	ListShowtimes(ctx context.Context, movieID int) ([]Showtime, error)
}

func NewService(repo Repository, notificationClient notificationClient.EmailServiceClient, authClient authClient.JWT_TokenServiceClient, movieBooking movie_booking.MovieServiceClient, theaterClient movie_booking.TheatreServiceClient) Service {
	return &service{
		repo:               repo,
		notificationClient: notificationClient,
		authClient:         authClient,
		theaterClient:      theaterClient,
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
		return "", errors.New("no admin exist")
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

// Theater
func (s *service) AddTheater(ctx context.Context, theater *Theater) error {
	_, err := s.theaterClient.AddTheater(ctx, &movie_booking.AddTheaterRequest{
		Name:            theater.Name,
		Place:           theater.Place,
		City:            theater.City,
		District:        theater.District,
		State:           theater.State,
		OwnerId:         uint32(theater.OwnerID),
		NumberOfScreens: int32(theater.NumberOfScreens),
		TheaterTypeId:   int32(theater.TheaterTypeID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterByID(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteTheaterByID(ctx, &movie_booking.DeleteTheaterRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterByName(ctx context.Context, name string) error {
	_, err := s.theaterClient.DeleteTheaterByName(ctx, &movie_booking.DeleteTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterByID(ctx context.Context, id int) (*Theater, error) {
	response, err := s.theaterClient.GetTheaterByID(ctx, &movie_booking.GetTheaterByIDRequest{
		TheaterId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &Theater{
		Name:            response.Theater.Name,
		Place:           response.Theater.Place,
		City:            response.Theater.City,
		District:        response.Theater.District,
		State:           response.Theater.State,
		OwnerID:         uint(response.Theater.OwnerId),
		NumberOfScreens: int(response.Theater.NumberOfScreens),
		TheaterTypeID:   int(response.Theater.TheaterTypeId),
	}, nil
}

func (s *service) GetTheaterByName(ctx context.Context, name string) ([]Theater, error) {
	response, err := s.theaterClient.GetTheaterByName(ctx, &movie_booking.GetTheaterByNameRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	var theaterResponses []Theater
	for _, theater := range response.Theater {
		theaterResponse := Theater{
			ID:              uint(theater.TheaterId),
			Name:            theater.Name,
			Place:           theater.Place,
			City:            theater.City,
			District:        theater.District,
			State:           theater.State,
			OwnerID:         uint(theater.OwnerId),
			NumberOfScreens: int(theater.NumberOfScreens),
			TheaterTypeID:   int(theater.TheaterTypeId),
		}
		theaterResponses = append(theaterResponses, theaterResponse)
	}
	return theaterResponses, nil
}

func (s *service) ListTheaters(ctx context.Context) ([]Theater, error) {
	response, err := s.theaterClient.ListTheaters(ctx, &movie_booking.ListTheatersRequest{})
	if err != nil {
		return nil, err
	}
	theaters := []Theater{}

	for _, res := range response.Theaters {
		theater := Theater{
			ID:              uint(res.TheaterId),
			Name:            res.Name,
			Place:           res.Place,
			City:            res.City,
			District:        res.District,
			State:           res.State,
			OwnerID:         uint(res.OwnerId),
			NumberOfScreens: int(res.NumberOfScreens),
			TheaterTypeID:   int(res.TheaterTypeId),
		}
		theaters = append(theaters, theater)
	}
	return theaters, nil
}

func (s *service) UpdateTheater(ctx context.Context, id int, theater Theater) error {
	_, err := s.theaterClient.UpdateTheater(ctx, &movie_booking.UpdateTheaterRequest{
		TheaterId:       int32(id),
		Name:            theater.Name,
		Place:           theater.Place,
		City:            theater.City,
		District:        theater.District,
		State:           theater.State,
		OwnerId:         uint32(theater.OwnerID),
		NumberOfScreens: int32(theater.NumberOfScreens),
		TheaterTypeId:   int32(theater.TheaterTypeID),
	},
	)
	if err != nil {
		return err
	}
	return nil
}

// Movies
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

//Theater-types

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

// Screen types
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

//Seat categories

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

// TheaterScreen
func (s *service) AddTheaterScreen(ctx context.Context, theaterScreen TheaterScreen) error {
	_, err := s.theaterClient.AddTheaterScreen(ctx, &movie_booking.AddTheaterScreenRequest{
		TheaterScreen: &movie_booking.TheaterScreen{
			ID:           uint32(theaterScreen.ID),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterScreenByID(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteTheaterScreenByID(ctx, &movie_booking.DeleteTheaterScreenRequest{
		TheaterScreenId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTheaterScreenByNumber(ctx context.Context, theaterID, screenNumber int) error {
	_, err := s.theaterClient.DeleteTheaterScreenByNumber(ctx, &movie_booking.DeleteTheaterScreenByNumberRequest{
		TheaterID:    int32(theaterID),
		ScreenNumber: int32(screenNumber),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTheaterScreenByID(ctx context.Context, id int) (*TheaterScreen, error) {
	response, err := s.theaterClient.GetTheaterScreenByID(ctx, &movie_booking.GetTheaterScreenByIDRequest{
		TheaterScreenId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterScreen{
		ID:           uint(response.TheaterScreen.ID),
		TheaterID:    int(response.TheaterScreen.TheaterID),
		ScreenNumber: int(response.TheaterScreen.ScreenNumber),
		SeatCapacity: int(response.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(response.TheaterScreen.ScreenTypeID),
	}, nil
}

func (s *service) GetTheaterScreenByNumber(ctx context.Context, theaterID, screenNumber int) (*TheaterScreen, error) {
	response, err := s.theaterClient.GetTheaterScreenByNumber(ctx, &movie_booking.GetTheaterScreenByNumberRequest{
		TheaterID:    int32(theaterID),
		ScreenNumber: int32(screenNumber),
	})
	if err != nil {
		return nil, err
	}
	return &TheaterScreen{
		ID:           uint(response.TheaterScreen.ID),
		TheaterID:    int(response.TheaterScreen.TheaterID),
		ScreenNumber: int(response.TheaterScreen.ScreenNumber),
		SeatCapacity: int(response.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(response.TheaterScreen.ScreenTypeID),
	}, nil
}

func (s *service) ListTheaterScreens(ctx context.Context, theaterID int) ([]TheaterScreen, error) {
	response, err := s.theaterClient.ListTheaterScreens(ctx, &movie_booking.ListTheaterScreensRequest{
		TheaterID: int32(theaterID),
	})
	if err != nil {
		return nil, err
	}
	theaterScreens := []TheaterScreen{}

	for _, res := range response.TheaterScreens {
		theaterScreen := TheaterScreen{
			ID:           uint(res.ID),
			TheaterID:    int(res.TheaterID),
			ScreenNumber: int(res.ScreenNumber),
			SeatCapacity: int(res.SeatCapacity),
			ScreenTypeID: int(res.ScreenTypeID),
		}
		theaterScreens = append(theaterScreens, theaterScreen)
	}
	return theaterScreens, nil
}

func (s *service) UpdateTheaterScreen(ctx context.Context, id int, theaterScreen TheaterScreen) error {
	_, err := s.theaterClient.UpdateTheaterScreen(ctx, &movie_booking.UpdateTheaterScreenRequest{
		TheaterScreen: &movie_booking.TheaterScreen{
			ID:           uint32(id),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// Show time
// Show Times
func (s *service) AddShowtime(ctx context.Context, showtime Showtime) error {
	_, err := s.theaterClient.AddShowtime(ctx, &movie_booking.AddShowtimeRequest{
		Showtime: &movie_booking.Showtime{
			Id:       uint32(showtime.ID),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShowtimeByID(ctx context.Context, id int) error {
	_, err := s.theaterClient.DeleteShowtimeByID(ctx, &movie_booking.DeleteShowtimeRequest{
		ShowtimeId: int32(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteShowtimeByDetails(ctx context.Context, movieID, screenID int, showDate, showTime time.Time) error {
	_, err := s.theaterClient.DeleteShowtimeByDetails(ctx, &movie_booking.DeleteShowtimeByDetailsRequest{
		MovieId:  int32(movieID),
		ScreenId: int32(screenID),
		ShowDate: timestamppb.New(showDate),
		ShowTime: timestamppb.New(showTime),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetShowtimeByID(ctx context.Context, id int) (*Showtime, error) {
	response, err := s.theaterClient.GetShowtimeByID(ctx, &movie_booking.GetShowtimeByIDRequest{
		ShowtimeId: int32(id),
	})
	if err != nil {
		return nil, err
	}
	return &Showtime{
		ID:       uint(response.Showtime.Id),
		MovieID:  int(response.Showtime.MovieId),
		ScreenID: int(response.Showtime.ScreenId),
		ShowDate: response.Showtime.ShowDate.AsTime(),
		ShowTime: response.Showtime.ShowTime.AsTime(),
	}, nil
}

func (s *service) GetShowtimeByDetails(ctx context.Context, movieID, screenID int, showDate, showTime time.Time) (*Showtime, error) {
	response, err := s.theaterClient.GetShowtimeByDetails(ctx, &movie_booking.GetShowtimeByDetailsRequest{
		MovieId:  int32(movieID),
		ScreenId: int32(screenID),
		ShowDate: timestamppb.New(showDate),
		ShowTime: timestamppb.New(showTime),
	})
	if err != nil {
		return nil, err
	}
	return &Showtime{
		ID:       uint(response.Showtime.Id),
		MovieID:  int(response.Showtime.MovieId),
		ScreenID: int(response.Showtime.ScreenId),
		ShowDate: response.Showtime.ShowDate.AsTime(),
		ShowTime: response.Showtime.ShowTime.AsTime(),
	}, nil
}

func (s *service) UpdateShowtime(ctx context.Context, id int, showtime Showtime) error {
	_, err := s.theaterClient.UpdateShowtime(ctx, &movie_booking.UpdateShowtimeRequest{
		Showtime: &movie_booking.Showtime{
			Id:       uint32(id),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) ListShowtimes(ctx context.Context, movieID int) ([]Showtime, error) {
	response, err := s.theaterClient.ListShowtimes(ctx, &movie_booking.ListShowtimesRequest{
		MovieId: int32(movieID),
	})
	if err != nil {
		return nil, err
	}
	showtimes := []Showtime{}

	for _, res := range response.Showtimes {
		showtime := Showtime{
			ID:       uint(res.Id),
			MovieID:  int(res.MovieId),
			ScreenID: int(res.ScreenId),
			ShowDate: res.ShowDate.AsTime(),
			ShowTime: res.ShowTime.AsTime(), // Convert string to time.Time if needed
		}
		showtimes = append(showtimes, showtime)
	}
	return showtimes, nil
}
