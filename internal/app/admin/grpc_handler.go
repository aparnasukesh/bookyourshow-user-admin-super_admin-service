package admin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
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

// Theater
func (h *GrpcHandler) AddTheater(ctx context.Context, req *user_admin.AddTheaterRequest) (*user_admin.AddTheaterResponse, error) {
	if err := h.svc.AddTheater(ctx, Theater{
		Name:            req.Name,
		Location:        req.Location,
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
			Location:        theater.Location,
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
	return &user_admin.GetTheaterByNameResponse{
		Theater: &user_admin.Theater{
			TheaterId:       int32(theater.ID),
			Name:            theater.Name,
			Location:        theater.Location,
			OwnerId:         uint32(theater.OwnerID),
			NumberOfScreens: int32(theater.NumberOfScreens),
			TheaterTypeId:   int32(theater.TheaterTypeID),
		},
	}, nil
}

func (h *GrpcHandler) UpdateTheater(ctx context.Context, req *user_admin.UpdateTheaterRequest) (*user_admin.UpdateTheaterResponse, error) {
	err := h.svc.UpdateTheater(ctx, int(req.TheaterId), Theater{
		Name:            req.Name,
		Location:        req.Location,
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
			Location:        m.Location,
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
			Id:                int32(m.ID),
			SeatCategoryName:  m.SeatCategoryName,
			SeatCategoryPrice: m.SeatCategoryPrice,
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
	}); err != nil {
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
	})
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
