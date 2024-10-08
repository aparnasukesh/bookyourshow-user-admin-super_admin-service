package admin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcHandler struct {
	svc Service
	user_admin.UnimplementedAdminServiceServer
}

func NewGrpcHandler(service Service) GrpcHandler {
	return GrpcHandler{
		svc: service,
	}
}

// Seats
func (h *GrpcHandler) CreateSeats(ctx context.Context, req *user_admin.CreateSeatsRequest) (*user_admin.CreateSeatsResponse, error) {
	rowSeatCatetoryPrice := []RowSeatCategoryPrice{}
	for _, res := range req.RowseatCategories {
		rowSeats := RowSeatCategoryPrice{
			RowStart:          res.RowStart,
			RowEnd:            res.RowEnd,
			SeatCategoryId:    int(res.SeatCategoryId),
			SeatCategoryPrice: float32(res.SeatCategoryPrice),
		}
		rowSeatCatetoryPrice = append(rowSeatCatetoryPrice, rowSeats)
	}
	err := h.svc.CreateSeats(ctx, CreateSeatsRequest{
		ScreenId:     int(req.ScreenId),
		TotalRows:    int(req.TotalRows),
		TotalColumns: int(req.TotalColumns),
		SeatRequest:  rowSeatCatetoryPrice,
	}, int(req.OwnerId))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) GetSeatsByScreenID(ctx context.Context, req *user_admin.GetSeatsByScreenIDRequest) (*user_admin.GetSeatsByScreenIDResponse, error) {
	seats, err := h.svc.GetSeatsByScreenId(ctx, int(req.ScreenId))
	if err != nil {
		return nil, err
	}

	grpcSeats := []*user_admin.Seat{}
	for _, seat := range seats {
		grpcSeat := &user_admin.Seat{
			Id:                int32(seat.ID),
			ScreenId:          int32(seat.ScreenID),
			SeatNumber:        seat.SeatNumber,
			Row:               seat.Row,
			Column:            int32(seat.Column),
			SeatCategoryId:    int32(seat.SeatCategoryID),
			SeatCategoryPrice: float64(seat.SeatCategoryPrice),
		}
		grpcSeats = append(grpcSeats, grpcSeat)
	}

	return &user_admin.GetSeatsByScreenIDResponse{Seats: grpcSeats}, nil
}

func (h *GrpcHandler) GetSeatByID(ctx context.Context, req *user_admin.GetSeatByIdRequest) (*user_admin.GetSeatByIdResponse, error) {
	seat, err := h.svc.GetSeatById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &user_admin.GetSeatByIdResponse{
		Seat: &user_admin.Seat{
			Id:                int32(seat.ID),
			ScreenId:          int32(seat.ScreenID),
			SeatNumber:        seat.SeatNumber,
			Row:               seat.Row,
			Column:            int32(seat.Column),
			SeatCategoryId:    int32(seat.SeatCategoryID),
			SeatCategoryPrice: float64(seat.SeatCategoryPrice),
		},
	}, nil
}

func (h *GrpcHandler) GetSeatBySeatNumberAndScreenID(ctx context.Context, req *user_admin.GetSeatBySeatNumberAndScreenIdRequest) (*user_admin.GetSeatBySeatNumberAndScreenIdResponse, error) {
	seat, err := h.svc.GetSeatBySeatNumberAndScreenId(ctx, int(req.ScreenId), req.SeatNumber)
	if err != nil {
		return nil, err
	}

	return &user_admin.GetSeatBySeatNumberAndScreenIdResponse{
		Seat: &user_admin.Seat{
			Id:                int32(seat.ID),
			ScreenId:          int32(seat.ScreenID),
			SeatNumber:        seat.SeatNumber,
			Row:               seat.Row,
			Column:            int32(seat.Column),
			SeatCategoryId:    int32(seat.SeatCategoryID),
			SeatCategoryPrice: float64(seat.SeatCategoryPrice),
		},
	}, nil
}

func (h *GrpcHandler) DeleteSeatByID(ctx context.Context, req *user_admin.DeleteSeatByIdRequest) (*user_admin.DeleteSeatByIdResponse, error) {
	err := h.svc.DeleteSeatById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &user_admin.DeleteSeatByIdResponse{}, nil
}

func (h *GrpcHandler) DeleteSeatBySeatNumberAndScreenID(ctx context.Context, req *user_admin.DeleteSeatBySeatNumberAndScreenIDRequest) (*user_admin.DeleteSeatBySeatNumberAndScreenIDResponse, error) {
	err := h.svc.DeleteSeatBySeatNumberAndScreenId(ctx, int(req.ScreenId), req.SeatNumber)
	if err != nil {
		return nil, err
	}

	return &user_admin.DeleteSeatBySeatNumberAndScreenIDResponse{}, nil
}

// Admin
func (h *GrpcHandler) RegisterAdmin(ctx context.Context, req *user_admin.RegisterAdminRequest) (*user_admin.RegisterAdminResponse, error) {

	userData := Admin{
		Username:    req.Username,
		Password:    req.Password,
		PhoneNumber: req.Phone,
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
	}
	if err := h.svc.RegisterAdmin(ctx, userData); err != nil {
		return nil, err
	}
	return &user_admin.RegisterAdminResponse{
		Status:     "request pending",
		StatusCode: 200,
	}, nil
}

func (h *GrpcHandler) LoginAdmin(ctx context.Context, req *user_admin.LoginAdminRequest) (*user_admin.LoginAdminResponse, error) {
	userData := Admin{
		Email:    req.Email,
		Password: req.Password,
	}
	token, err := h.svc.LoginAdmin(ctx, userData)
	if err != nil {
		return nil, err
	}
	return &user_admin.LoginAdminResponse{
		Status:     "login successfull",
		StatusCode: 200,
		Token:      token,
	}, nil
}

func (h *GrpcHandler) GetAdminProfile(ctx context.Context, req *user_admin.GetProfileRequest) (*user_admin.GetProfileResponse, error) {
	admin, err := h.svc.GetAdminProfile(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetProfileResponse{
		Status:     "success",
		StatusCode: 200,
		ProfileDetails: &user_admin.Admin{
			Id:          int32(admin.ID),
			Username:    admin.Username,
			Phone:       admin.PhoneNumber,
			Email:       admin.Email,
			FirstName:   admin.FirstName,
			LastName:    admin.LastName,
			Gender:      admin.Gender,
			DateOfBirth: admin.DateOfBirth,
			IsVerified:  admin.IsVerified,
		},
	}, nil
}

func (h *GrpcHandler) UpdateAdminProfile(ctx context.Context, req *user_admin.UpdateAdminProfileRequest) (*user_admin.UpdateAdminProfileResponse, error) {
	err := h.svc.UpdateAdminProfile(ctx, int(req.UserId), AdminProfileDetails{
		Username:    req.Username,
		PhoneNumber: req.Phone,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) ForgotAdminPassword(ctx context.Context, req *user_admin.ForgotPasswordRequest) (*user_admin.ForgotPasswordResponse, error) {
	err := h.svc.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) ResetAdminPassword(ctx context.Context, req *user_admin.ResetPasswordRequest) (*user_admin.ResetPasswordResponse, error) {
	err := h.svc.ResetPassword(ctx, ResetPassword{
		Email:       req.Email,
		Otp:         req.Otp,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Theater
func (h *GrpcHandler) AddTheater(ctx context.Context, req *user_admin.AddTheaterRequest) (*user_admin.AddTheaterResponse, error) {
	if err := h.svc.AddTheater(ctx, &Theater{
		Name:            req.Name,
		Place:           req.Place,
		City:            req.City,
		District:        req.District,
		State:           req.State,
		OwnerID:         uint(req.OwnerId),
		NumberOfScreens: int(req.NumberOfScreens),
		TheaterTypeID:   int(req.TheaterTypeId),
	}); err != nil {
		return &user_admin.AddTheaterResponse{}, err
	}
	return &user_admin.AddTheaterResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterByID(ctx context.Context, req *user_admin.DeleteTheaterRequest) (*user_admin.DeleteTheaterResponse, error) {
	if err := h.svc.DeleteTheaterByID(ctx, int(req.TheaterId)); err != nil {
		return &user_admin.DeleteTheaterResponse{}, err
	}
	return &user_admin.DeleteTheaterResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterByName(ctx context.Context, req *user_admin.DeleteTheaterByNameRequest) (*user_admin.DeleteTheaterByNameResponse, error) {
	if err := h.svc.DeleteTheaterByName(ctx, req.Name); err != nil {
		return &user_admin.DeleteTheaterByNameResponse{}, err
	}
	return &user_admin.DeleteTheaterByNameResponse{}, nil
}

func (h *GrpcHandler) GetTheaterByID(ctx context.Context, req *user_admin.GetTheaterByIDRequest) (*user_admin.GetTheaterByIDResponse, error) {
	theater, err := h.svc.GetTheaterByID(ctx, int(req.TheaterId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetTheaterByIDResponse{
		Theater: &user_admin.Theater{
			TheaterId:       int32(theater.ID),
			Name:            theater.Name,
			Place:           theater.Place,
			City:            theater.City,
			District:        theater.District,
			State:           theater.State,
			OwnerId:         uint32(theater.OwnerID),
			NumberOfScreens: int32(theater.NumberOfScreens),
			TheaterTypeId:   int32(theater.TheaterTypeID),
		},
	}, nil
}

func (h *GrpcHandler) GetTheaterByName(ctx context.Context, req *user_admin.GetTheaterByNameRequest) (*user_admin.GetTheaterByNameResponse, error) {
	theater, err := h.svc.GetTheaterByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	theaters := []*user_admin.Theater{}

	for _, data := range theater {
		theaterData := &user_admin.Theater{
			TheaterId:       int32(data.ID),
			Name:            data.Name,
			Place:           data.Place,
			City:            data.City,
			District:        data.District,
			State:           data.State,
			OwnerId:         uint32(data.OwnerID),
			NumberOfScreens: int32(data.NumberOfScreens),
			TheaterTypeId:   int32(data.TheaterTypeID),
		}
		theaters = append(theaters, theaterData)
	}
	return &user_admin.GetTheaterByNameResponse{
		Theater: theaters,
	}, err
}

func (h *GrpcHandler) UpdateTheater(ctx context.Context, req *user_admin.UpdateTheaterRequest) (*user_admin.UpdateTheaterResponse, error) {
	err := h.svc.UpdateTheater(ctx, int(req.TheaterId), Theater{
		Name:            req.Name,
		Place:           req.Place,
		City:            req.City,
		District:        req.District,
		State:           req.State,
		OwnerID:         uint(req.OwnerId),
		NumberOfScreens: int(req.NumberOfScreens),
		TheaterTypeID:   int(req.TheaterTypeId),
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateTheaterResponse{}, nil
}

func (h *GrpcHandler) ListTheaters(ctx context.Context, req *user_admin.ListTheatersRequest) (*user_admin.ListTheatersResponse, error) {
	response, err := h.svc.ListTheaters(ctx)
	if err != nil {
		return nil, err
	}

	var grpcTheaters []*user_admin.Theater
	for _, m := range response {
		grpcTheater := &user_admin.Theater{
			TheaterId:       int32(m.ID),
			Name:            m.Name,
			Place:           m.Place,
			City:            m.City,
			District:        m.District,
			State:           m.State,
			OwnerId:         uint32(m.OwnerID),
			NumberOfScreens: int32(m.NumberOfScreens),
			TheaterTypeId:   int32(m.TheaterTypeID),
		}
		grpcTheaters = append(grpcTheaters, grpcTheater)
	}

	return &user_admin.ListTheatersResponse{
		Theaters: grpcTheaters,
	}, nil
}

// Movies
func (h *GrpcHandler) ListMovies(ctx context.Context, req *user_admin.ListMoviesRequest) (*user_admin.ListMoviesResponse, error) {
	response, err := h.svc.ListMovies(ctx)
	if err != nil {
		return nil, err
	}

	var grpcMovies []*user_admin.Movie
	for _, m := range response {
		grpcMovie := &user_admin.Movie{
			Title:       m.Title,
			Description: m.Description,
			Duration:    int32(m.Duration),
			Genre:       m.Genre,
			ReleaseDate: m.ReleaseDate,
			Rating:      float32(m.Rating),
			Language:    m.Language,
		}
		grpcMovies = append(grpcMovies, grpcMovie)
	}

	return &user_admin.ListMoviesResponse{
		Movies: grpcMovies,
	}, nil
}

// Theater types
func (h *GrpcHandler) ListTheaterTypes(ctx context.Context, req *user_admin.ListTheaterTypesRequest) (*user_admin.ListTheaterTypeResponse, error) {
	response, err := h.svc.ListTheaterTypes(ctx)
	if err != nil {
		return nil, err
	}

	var grpcTheaterTypes []*user_admin.TheaterType
	for _, m := range response {
		grpcTheaterType := &user_admin.TheaterType{
			Id:              int32(m.ID),
			TheaterTypeName: m.TheaterTypeName,
		}
		grpcTheaterTypes = append(grpcTheaterTypes, grpcTheaterType)
	}

	return &user_admin.ListTheaterTypeResponse{
		TheaterTypes: grpcTheaterTypes,
	}, nil
}

// Seat categories
func (h *GrpcHandler) ListSeatCategories(ctx context.Context, req *user_admin.ListSeatCategoriesRequest) (*user_admin.ListSeatCategoriesResponse, error) {
	response, err := h.svc.ListSeatCategories(ctx)
	if err != nil {
		return nil, err
	}

	var grpcSeatCategories []*user_admin.SeatCategory
	for _, m := range response {
		grpcSeatCategory := &user_admin.SeatCategory{
			Id:               int32(m.ID),
			SeatCategoryName: m.SeatCategoryName,
		}
		grpcSeatCategories = append(grpcSeatCategories, grpcSeatCategory)
	}

	return &user_admin.ListSeatCategoriesResponse{
		SeatCategories: grpcSeatCategories,
	}, nil
}

// Screen type
func (h *GrpcHandler) ListScreenTypes(ctx context.Context, req *user_admin.ListScreenTypesRequest) (*user_admin.ListScreenTypesResponse, error) {
	response, err := h.svc.ListScreenTypes(ctx)
	if err != nil {
		return nil, err
	}

	var grpcScreenTypes []*user_admin.ScreenType
	for _, m := range response {
		grpcScreenType := &user_admin.ScreenType{
			Id:             int32(m.ID),
			ScreenTypeName: m.ScreenTypeName,
		}
		grpcScreenTypes = append(grpcScreenTypes, grpcScreenType)
	}

	return &user_admin.ListScreenTypesResponse{
		ScreenTypes: grpcScreenTypes,
	}, nil
}

// Theater screen
func (h *GrpcHandler) AddTheaterScreen(ctx context.Context, req *user_admin.AddTheaterScreenRequest) (*user_admin.AddTheaterScreenResponse, error) {
	if err := h.svc.AddTheaterScreen(ctx, TheaterScreen{
		TheaterID:    int(req.TheaterScreen.TheaterID),
		ScreenNumber: int(req.TheaterScreen.ScreenNumber),
		SeatCapacity: int(req.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(req.TheaterScreen.ScreenTypeID),
	}, int(req.OwnerId)); err != nil {
		return &user_admin.AddTheaterScreenResponse{}, err
	}
	return &user_admin.AddTheaterScreenResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterScreenByID(ctx context.Context, req *user_admin.DeleteTheaterScreenRequest) (*user_admin.DeleteTheaterScreenResponse, error) {
	if err := h.svc.DeleteTheaterScreenByID(ctx, int(req.TheaterScreenId)); err != nil {
		return &user_admin.DeleteTheaterScreenResponse{}, err
	}
	return &user_admin.DeleteTheaterScreenResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterScreenByNumber(ctx context.Context, req *user_admin.DeleteTheaterScreenByNumberRequest) (*user_admin.DeleteTheaterScreenByNumberResponse, error) {
	if err := h.svc.DeleteTheaterScreenByNumber(ctx, int(req.TheaterID), int(req.ScreenNumber)); err != nil {
		return &user_admin.DeleteTheaterScreenByNumberResponse{}, err
	}
	return &user_admin.DeleteTheaterScreenByNumberResponse{}, nil
}

func (h *GrpcHandler) GetTheaterScreenByID(ctx context.Context, req *user_admin.GetTheaterScreenByIDRequest) (*user_admin.GetTheaterScreenByIDResponse, error) {
	theaterScreen, err := h.svc.GetTheaterScreenByID(ctx, int(req.TheaterScreenId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetTheaterScreenByIDResponse{
		TheaterScreen: &user_admin.TheaterScreen{
			ID:           uint32(theaterScreen.ID),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	}, nil
}

func (h *GrpcHandler) GetTheaterScreenByNumber(ctx context.Context, req *user_admin.GetTheaterScreenByNumberRequest) (*user_admin.GetTheaterScreenByNumberResponse, error) {
	theaterScreen, err := h.svc.GetTheaterScreenByNumber(ctx, int(req.TheaterID), int(req.ScreenNumber))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetTheaterScreenByNumberResponse{
		TheaterScreen: &user_admin.TheaterScreen{
			ID:           uint32(theaterScreen.ID),
			TheaterID:    int32(theaterScreen.TheaterID),
			ScreenNumber: int32(theaterScreen.ScreenNumber),
			SeatCapacity: int32(theaterScreen.SeatCapacity),
			ScreenTypeID: int32(theaterScreen.ScreenTypeID),
		},
	}, nil
}

func (h *GrpcHandler) UpdateTheaterScreen(ctx context.Context, req *user_admin.UpdateTheaterScreenRequest) (*user_admin.UpdateTheaterScreenResponse, error) {
	err := h.svc.UpdateTheaterScreen(ctx, int(req.TheaterScreen.ID), TheaterScreen{
		TheaterID:    int(req.TheaterScreen.TheaterID),
		ScreenNumber: int(req.TheaterScreen.ScreenNumber),
		SeatCapacity: int(req.TheaterScreen.SeatCapacity),
		ScreenTypeID: int(req.TheaterScreen.ScreenTypeID),
	}, int(req.OwnerId))
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateTheaterScreenResponse{}, nil
}

func (h *GrpcHandler) ListTheaterScreens(ctx context.Context, req *user_admin.ListTheaterScreensRequest) (*user_admin.ListTheaterScreensResponse, error) {
	response, err := h.svc.ListTheaterScreens(ctx, int(req.TheaterID))
	if err != nil {
		return nil, err
	}

	var grpcTheaterScreens []*user_admin.TheaterScreen
	for _, m := range response {
		grpcTheaterScreen := &user_admin.TheaterScreen{
			ID:           uint32(m.ID),
			TheaterID:    int32(m.TheaterID),
			ScreenNumber: int32(m.ScreenNumber),
			SeatCapacity: int32(m.SeatCapacity),
			ScreenTypeID: int32(m.ScreenTypeID),
		}
		grpcTheaterScreens = append(grpcTheaterScreens, grpcTheaterScreen)
	}

	return &user_admin.ListTheaterScreensResponse{
		TheaterScreens: grpcTheaterScreens,
	}, nil
}

// Show time
func (h *GrpcHandler) AddShowtime(ctx context.Context, req *user_admin.AddShowtimeRequest) (*user_admin.AddShowtimeResponse, error) {
	if err := h.svc.AddShowtime(ctx, Showtime{
		ID:       uint(req.Showtime.Id),
		MovieID:  int(req.Showtime.MovieId),
		ScreenID: int(req.Showtime.ScreenId),
		ShowDate: req.Showtime.ShowDate.AsTime(),
		ShowTime: req.Showtime.ShowTime.AsTime(),
	}, int(req.OwnerId)); err != nil {
		return &user_admin.AddShowtimeResponse{}, err
	}
	return &user_admin.AddShowtimeResponse{}, nil
}

func (h *GrpcHandler) DeleteShowtimeByID(ctx context.Context, req *user_admin.DeleteShowtimeRequest) (*user_admin.DeleteShowtimeResponse, error) {
	if err := h.svc.DeleteShowtimeByID(ctx, int(req.ShowtimeId)); err != nil {
		return &user_admin.DeleteShowtimeResponse{}, err
	}
	return &user_admin.DeleteShowtimeResponse{}, nil
}

func (h *GrpcHandler) DeleteShowtimeByDetails(ctx context.Context, req *user_admin.DeleteShowtimeByDetailsRequest) (*user_admin.DeleteShowtimeByDetailsResponse, error) {
	if err := h.svc.DeleteShowtimeByDetails(ctx, int(req.MovieId), int(req.ScreenId), req.ShowDate.AsTime(), req.ShowTime.AsTime()); err != nil {
		return &user_admin.DeleteShowtimeByDetailsResponse{}, err
	}
	return &user_admin.DeleteShowtimeByDetailsResponse{}, nil
}

func (h *GrpcHandler) GetShowtimeByID(ctx context.Context, req *user_admin.GetShowtimeByIDRequest) (*user_admin.GetShowtimeByIDResponse, error) {
	showtime, err := h.svc.GetShowtimeByID(ctx, int(req.ShowtimeId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetShowtimeByIDResponse{
		Showtime: &user_admin.Showtime{
			Id:       uint32(showtime.ID),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	}, nil
}

func (h *GrpcHandler) GetShowtimeByDetails(ctx context.Context, req *user_admin.GetShowtimeByDetailsRequest) (*user_admin.GetShowtimeByDetailsResponse, error) {
	showtime, err := h.svc.GetShowtimeByDetails(ctx, int(req.MovieId), int(req.ScreenId), req.ShowDate.AsTime(), req.ShowTime.AsTime())
	if err != nil {
		return nil, err
	}
	return &user_admin.GetShowtimeByDetailsResponse{
		Showtime: &user_admin.Showtime{
			Id:       uint32(showtime.ID),
			MovieId:  int32(showtime.MovieID),
			ScreenId: int32(showtime.ScreenID),
			ShowDate: timestamppb.New(showtime.ShowDate),
			ShowTime: timestamppb.New(showtime.ShowTime),
		},
	}, nil
}

func (h *GrpcHandler) UpdateShowtime(ctx context.Context, req *user_admin.UpdateShowtimeRequest) (*user_admin.UpdateShowtimeResponse, error) {
	err := h.svc.UpdateShowtime(ctx, int(req.Showtime.Id), Showtime{
		ID:       uint(req.Showtime.Id),
		MovieID:  int(req.Showtime.MovieId),
		ScreenID: int(req.Showtime.ScreenId),
		ShowDate: req.Showtime.ShowDate.AsTime(),
		ShowTime: req.Showtime.ShowTime.AsTime(),
	}, int(req.OwnerId))
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateShowtimeResponse{}, nil
}

func (h *GrpcHandler) ListShowtimes(ctx context.Context, req *user_admin.ListShowtimesRequest) (*user_admin.ListShowtimesResponse, error) {
	response, err := h.svc.ListShowtimes(ctx, int(req.MovieId))
	if err != nil {
		return nil, err
	}

	var grpcShowtimes []*user_admin.Showtime
	for _, m := range response {
		grpcShowtime := &user_admin.Showtime{
			Id:       uint32(m.ID),
			MovieId:  int32(m.MovieID),
			ScreenId: int32(m.ScreenID),
			ShowDate: timestamppb.New(m.ShowDate),
			ShowTime: timestamppb.New(m.ShowTime),
		}
		grpcShowtimes = append(grpcShowtimes, grpcShowtime)
	}

	return &user_admin.ListShowtimesResponse{
		Showtimes: grpcShowtimes,
	}, nil
}

// Movie Schedule
func (h *GrpcHandler) AddMovieSchedule(ctx context.Context, req *user_admin.AddMovieScheduleRequest) (*user_admin.AddMovieScheduleResponse, error) {
	err := h.svc.AddMovieSchedule(ctx, MovieSchedule{
		MovieID:    int(req.MovieSchedule.MovieId),
		TheaterID:  int(req.MovieSchedule.TheaterId),
		ShowtimeID: int(req.MovieSchedule.ShowtimeId),
	}, int(req.OwnerId))
	if err != nil {
		return nil, err
	}
	return &user_admin.AddMovieScheduleResponse{}, nil
}

func (h *GrpcHandler) UpdateMovieSchedule(ctx context.Context, req *user_admin.UpdateMovieScheduleRequest) (*user_admin.UpdateMovieScheduleResponse, error) {
	err := h.svc.UpdateMovieSchedule(ctx, int(req.MovieSchedule.Id), MovieSchedule{
		MovieID:    int(req.MovieSchedule.MovieId),
		TheaterID:  int(req.MovieSchedule.TheaterId),
		ShowtimeID: int(req.MovieSchedule.ShowtimeId),
	}, int(req.OwnerId))
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateMovieScheduleResponse{}, nil
}

func (h *GrpcHandler) GetAllMovieSchedules(ctx context.Context, req *user_admin.GetAllMovieScheduleRequest) (*user_admin.GetAllMovieScheduleResponse, error) {
	movieSchedules, err := h.svc.GetAllMovieSchedules(ctx)
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetAllMovieScheduleResponse{MovieSchedules: scheduleList}, nil
}

func (h *GrpcHandler) GetMovieScheduleByMovieID(ctx context.Context, req *user_admin.GetMovieScheduleByMovieIdRequest) (*user_admin.GetMovieScheduleByMovieIdResponse, error) {
	movieSchedules, err := h.svc.GetMovieScheduleByMovieID(ctx, int(req.MovieId))
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetMovieScheduleByMovieIdResponse{MovieSchedules: scheduleList}, nil
}

func (h *GrpcHandler) GetMovieScheduleByTheaterID(ctx context.Context, req *user_admin.GetMovieScheduleByTheaterIdRequest) (*user_admin.GetMovieScheduleByTheaterIdResponse, error) {
	movieSchedules, err := h.svc.GetMovieScheduleByTheaterID(ctx, int(req.TheaterId))
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetMovieScheduleByTheaterIdResponse{MovieSchedules: scheduleList}, nil
}

func (h *GrpcHandler) GetMovieScheduleByMovieIdAndTheaterId(ctx context.Context, req *user_admin.GetMovieScheduleByMovieIdAndTheaterIdRequest) (*user_admin.GetMovieScheduleByMovieIdAndTheaterIdResponse, error) {
	movieSchedules, err := h.svc.GetMovieScheduleByMovieIdAndTheaterId(ctx, int(req.MovieId), int(req.TheaterId))
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetMovieScheduleByMovieIdAndTheaterIdResponse{MovieSchedules: scheduleList}, nil
}

func (h *GrpcHandler) GetMovieScheduleByMovieIdAndShowTimeId(ctx context.Context, req *user_admin.GetMovieScheduleByMovieIdAndShowTimeIdRequest) (*user_admin.GetMovieScheduleByMovieIdAndShowTimeIdResponse, error) {
	movieSchedules, err := h.svc.GetMovieScheduleByMovieIdAndShowTimeId(ctx, int(req.MovieId), int(req.ShowtimeId))
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetMovieScheduleByMovieIdAndShowTimeIdResponse{
		MovieSchedules: scheduleList,
		Message:        "",
	}, nil
}

func (h *GrpcHandler) GetMovieScheduleByTheaterIdAndShowTimeId(ctx context.Context, req *user_admin.GetGetMovieScheduleByTheaterIdAndShowTimeIdRequest) (*user_admin.GetGetMovieScheduleByTheaterIdAndShowTimeIdResponse, error) {
	movieSchedules, err := h.svc.GetMovieScheduleByTheaterIdAndShowTimeId(ctx, int(req.TheaterId), int(req.ShowtimeId))
	if err != nil {
		return nil, err
	}

	var scheduleList []*user_admin.MovieSchedule
	for _, schedule := range movieSchedules {
		scheduleList = append(scheduleList, &user_admin.MovieSchedule{
			Id:         int32(schedule.ID),
			MovieId:    int32(schedule.MovieID),
			TheaterId:  int32(schedule.TheaterID),
			ShowtimeId: int32(schedule.ShowtimeID),
		})
	}
	return &user_admin.GetGetMovieScheduleByTheaterIdAndShowTimeIdResponse{
		MovieSchedules: scheduleList,
		Message:        "",
	}, nil
}

func (h *GrpcHandler) GetMovieScheduleByID(ctx context.Context, req *user_admin.GetMovieScheduleByIDRequest) (*user_admin.GetMovieScheduleByIDResponse, error) {
	movieSchedule, err := h.svc.GetMovieScheduleByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &user_admin.GetMovieScheduleByIDResponse{
		MovieSchedule: &user_admin.MovieSchedule{
			Id:         int32(movieSchedule.ID),
			MovieId:    int32(movieSchedule.MovieID),
			TheaterId:  int32(movieSchedule.TheaterID),
			ShowtimeId: int32(movieSchedule.ShowtimeID),
		},
	}, nil
}

func (h *GrpcHandler) DeleteMovieScheduleById(ctx context.Context, req *user_admin.DeleteMovieScheduleByIdRequest) (*user_admin.DeleteMovieScheduleByIdResponse, error) {
	err := h.svc.DeleteMovieScheduleById(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &user_admin.DeleteMovieScheduleByIdResponse{}, nil
}

func (h *GrpcHandler) DeleteMovieScheduleByMovieIdAndTheaterId(ctx context.Context, req *user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdRequest) (*user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdResponse, error) {
	err := h.svc.DeleteMovieScheduleByMovieIdAndTheaterId(ctx, int(req.MovieId), int(req.TheaterId))
	if err != nil {
		return nil, err
	}
	return &user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdResponse{}, nil
}

func (h *GrpcHandler) DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx context.Context, req *user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeIdRequest) (*user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeIdResponse, error) {
	err := h.svc.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeId(ctx, int(req.MovieId), int(req.TheaterId), int(req.ShowtimeId))
	if err != nil {
		return nil, err
	}
	return &user_admin.DeleteMovieScheduleByMovieIdAndTheaterIdAndShowTimeIdResponse{}, nil
}
