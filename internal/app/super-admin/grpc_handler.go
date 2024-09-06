package superadmin

import (
	"context"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

type GrpcHandler struct {
	svc Service
	user_admin.UnimplementedSuperAdminServiceServer
}

func NewGrpcHandler(service Service) GrpcHandler {
	return GrpcHandler{
		svc: service,
	}
}

func (h *GrpcHandler) LoginSuperAdmin(ctx context.Context, req *user_admin.LoginSuperAdminRequest) (*user_admin.LoginSuperAdminResponse, error) {
	token, err := h.svc.LoginSuperAdmin(ctx, Admin{
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.LoginSuperAdminResponse{
		Status:     "login successfull",
		StatusCode: 200,
		Token:      token,
	}, nil
}

// admin
func (h *GrpcHandler) ListAdminRequests(ctx context.Context, req *user_admin.ListAdminRequestsRequest) (*user_admin.ListAdminRequestsResponse, error) {
	adminLists, err := h.svc.ListAdminRequests(ctx)
	if err != nil {
		return nil, err
	}

	emails := make([]*user_admin.Email, len(adminLists))
	for i, admin := range adminLists {
		emails[i] = &user_admin.Email{
			Email: admin.Email,
		}
	}

	return &user_admin.ListAdminRequestsResponse{
		Email: emails,
	}, nil
}

func (h *GrpcHandler) AdminApproval(ctx context.Context, req *user_admin.AdminApprovalRequest) (*user_admin.AdminApprovalResponse, error) {
	err := h.svc.AdminApproval(ctx, req.Email, req.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user_admin.AdminApprovalResponse{}, nil

}

func (h *GrpcHandler) ListAllAdmin(ctx context.Context, req *user_admin.ListAllAdminRequest) (*user_admin.ListAllAdminResponse, error) {
	admins, err := h.svc.ListAllAdmin(ctx)
	if err != nil {
		return nil, err
	}
	response := []*user_admin.Admin{}
	for _, res := range admins {
		admin := &user_admin.Admin{
			Id:        int32(res.ID),
			Username:  res.Username,
			Phone:     res.PhoneNumber,
			Email:     res.Email,
			FirstName: res.FirstName,
			LastName:  res.LastName,
			Gender:    res.Gender,
		}
		response = append(response, admin)
	}
	return &user_admin.ListAllAdminResponse{
		Admin: response,
	}, nil
}

func (h *GrpcHandler) GetAdminByID(ctx context.Context, req *user_admin.GetAdminByIdRequest) (*user_admin.GetAdminByIdResponse, error) {
	admin, err := h.svc.GetAdminByID(ctx, int(req.AdminId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetAdminByIdResponse{
		Admin: &user_admin.Admin{
			Id:        int32(admin.ID),
			Username:  admin.Username,
			Phone:     admin.PhoneNumber,
			Email:     admin.Email,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			Gender:    admin.Gender,
		},
	}, nil
}
func (h *GrpcHandler) DeleteMovie(ctx context.Context, req *user_admin.DeleteMovieRequest) (*user_admin.DeleteMovieResponse, error) {
	err := h.svc.DeleteMovie(ctx, int(req.MovieId))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) GetMovieDetails(ctx context.Context, req *user_admin.GetMovieDetailsRequest) (*user_admin.GetMovieDetailsResponse, error) {
	movie, err := h.svc.GetMovieDetails(ctx, int(req.MovieId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetMovieDetailsResponse{
		Movie: &user_admin.Movie{
			MovieId:     req.MovieId,
			Title:       movie.Title,
			Description: movie.Description,
			Duration:    int32(movie.Duration),
			Genre:       movie.Genre,
			ReleaseDate: movie.ReleaseDate,
			Rating:      float32(movie.Rating),
			Language:    movie.Language,
		},
	}, nil
}

func (h *GrpcHandler) UpdateMovie(ctx context.Context, req *user_admin.UpdateMovieRequest) (*user_admin.UpdateMovieResponse, error) {
	err := h.svc.UpdateMovie(ctx, Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    int(req.Duration),
		Genre:       req.Genre,
		ReleaseDate: req.ReleaseDate,
		Rating:      float64(req.Rating),
		Language:    req.Language,
	}, int(req.MovieId))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

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

// theater type
func (h *GrpcHandler) AddTheaterType(ctx context.Context, req *user_admin.AddTheaterTypeRequest) (*user_admin.AddTheaterTypeResponse, error) {
	if err := h.svc.AddTheaterType(ctx, TheaterType{
		TheaterTypeName: req.TheaterTypeName,
	}); err != nil {
		return nil, err
	}
	return &user_admin.AddTheaterTypeResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterTypeByID(ctx context.Context, req *user_admin.DeleteTheaterTypeRequest) (*user_admin.DeleteTheaterTypeResponse, error) {
	if err := h.svc.DeleteTheaterTypeByID(ctx, int(req.TheaterTypeId)); err != nil {
		return nil, err
	}
	return &user_admin.DeleteTheaterTypeResponse{}, nil
}

func (h *GrpcHandler) DeleteTheaterTypeByName(ctx context.Context, req *user_admin.DeleteTheaterTypeByNameRequest) (*user_admin.DeleteTheaterTypeByNameResponse, error) {
	if err := h.svc.DeleteTheaterTypeByName(ctx, req.Name); err != nil {
		return nil, err
	}
	return &user_admin.DeleteTheaterTypeByNameResponse{}, nil
}

func (h *GrpcHandler) GetTheaterTypeByID(ctx context.Context, req *user_admin.GetTheaterTypeByIDRequest) (*user_admin.GetTheaterTypeByIDResponse, error) {
	theaterType, err := h.svc.GetTheaterTypeByID(ctx, int(req.TheaterTypeId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetTheaterTypeByIDResponse{
		TheaterType: &user_admin.TheaterType{
			Id:              int32(theaterType.ID),
			TheaterTypeName: theaterType.TheaterTypeName,
		},
	}, nil
}

func (h *GrpcHandler) GetTheaterTypeByName(ctx context.Context, req *user_admin.GetTheaterTypeByNameRequest) (*user_admin.GetTheaterTypeBynameResponse, error) {
	theaterType, err := h.svc.GetTheaterTypeByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &user_admin.GetTheaterTypeBynameResponse{
		TheaterType: &user_admin.TheaterType{
			Id:              int32(theaterType.ID),
			TheaterTypeName: theaterType.TheaterTypeName,
		},
	}, nil
}

func (h *GrpcHandler) UpdateTheaterType(ctx context.Context, req *user_admin.UpdateTheaterTypeRequest) (*user_admin.UpdateTheaterTypeResponse, error) {
	err := h.svc.UpdateTheaterType(ctx, int(req.Id), TheaterType{
		TheaterTypeName: req.TheaterTypeName,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateTheaterTypeResponse{}, nil
}

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

// screen -types
func (h *GrpcHandler) AddScreenType(ctx context.Context, req *user_admin.AddScreenTypeRequest) (*user_admin.AddScreenTypeResponse, error) {
	if err := h.svc.AddScreenType(ctx, ScreenType{
		ScreenTypeName: req.ScreenTypeName,
	}); err != nil {
		return nil, err
	}
	return &user_admin.AddScreenTypeResponse{}, nil
}

func (h *GrpcHandler) DeleteScreenTypeByID(ctx context.Context, req *user_admin.DeleteScreenTypeRequest) (*user_admin.DeleteScreenTypeResponse, error) {
	if err := h.svc.DeleteScreenTypeById(ctx, int(req.ScreenTypeId)); err != nil {
		return nil, err
	}
	return &user_admin.DeleteScreenTypeResponse{}, nil
}

func (h *GrpcHandler) DeleteScreenTypeByName(ctx context.Context, req *user_admin.DeleteScreenTypeByNameRequest) (*user_admin.DeleteScreenTypeByNameResponse, error) {
	if err := h.svc.DeleteScreenTypeByName(ctx, req.Name); err != nil {
		return nil, err
	}
	return &user_admin.DeleteScreenTypeByNameResponse{}, nil
}

func (h *GrpcHandler) GetScreenTypeByID(ctx context.Context, req *user_admin.GetScreenTypeByIDRequest) (*user_admin.GetScreenTypeByIDResponse, error) {
	screenType, err := h.svc.GetScreenTypeByID(ctx, int(req.ScreenTypeId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetScreenTypeByIDResponse{
		ScreenType: &user_admin.ScreenType{
			Id:             int32(screenType.ID),
			ScreenTypeName: screenType.ScreenTypeName,
		},
	}, nil
}

func (h *GrpcHandler) GetScreenTypeByName(ctx context.Context, req *user_admin.GetScreenTypeByNameRequest) (*user_admin.GetScreenTypeByNameResponse, error) {
	screenType, err := h.svc.GetScreenTypeByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &user_admin.GetScreenTypeByNameResponse{
		ScreenType: &user_admin.ScreenType{
			Id:             int32(screenType.ID),
			ScreenTypeName: screenType.ScreenTypeName,
		},
	}, nil
}

func (h *GrpcHandler) UpdateScreenType(ctx context.Context, req *user_admin.UpdateScreenTypeRequest) (*user_admin.UpdateScreenTypeResponse, error) {
	err := h.svc.UpdateScreenType(ctx, int(req.Id), ScreenType{
		ScreenTypeName: req.ScreenTypeName,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateScreenTypeResponse{}, nil
}

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

// seat category
func (h *GrpcHandler) AddSeatCategory(ctx context.Context, req *user_admin.AddSeatCategoryRequest) (*user_admin.AddSeatCategoryResponse, error) {
	if err := h.svc.AddSeatCategory(ctx, SeatCategory{
		SeatCategoryName: req.SeatCategory.SeatCategoryName,
	}); err != nil {
		return &user_admin.AddSeatCategoryResponse{
			Message: "failed to add seat category",
		}, err
	}
	return &user_admin.AddSeatCategoryResponse{
		Message: "successfully added seat category",
	}, nil
}

func (h *GrpcHandler) DeleteSeatCategoryByID(ctx context.Context, req *user_admin.DeleteSeatCategoryRequest) (*user_admin.DeleteSeatCategoryResponse, error) {
	if err := h.svc.DeleteSeatCategoryByID(ctx, int(req.SeatCategoryId)); err != nil {
		return &user_admin.DeleteSeatCategoryResponse{
			Message: "failed to delete seat category",
		}, err
	}
	return &user_admin.DeleteSeatCategoryResponse{
		Message: "successfully deleted seat category",
	}, nil
}

func (h *GrpcHandler) DeleteSeatCategoryByName(ctx context.Context, req *user_admin.DeleteSeatCategoryByNameRequest) (*user_admin.DeleteSeatCategoryByNameResponse, error) {
	if err := h.svc.DeleteSeatCategoryByName(ctx, req.Name); err != nil {
		return &user_admin.DeleteSeatCategoryByNameResponse{
			Message: "failed to delete seat category",
		}, err
	}
	return &user_admin.DeleteSeatCategoryByNameResponse{
		Message: "successfully deleted seat category",
	}, nil
}

func (h *GrpcHandler) GetSeatCategoryByID(ctx context.Context, req *user_admin.GetSeatCategoryByIDRequest) (*user_admin.GetSeatCategoryByIDResponse, error) {
	seatCategory, err := h.svc.GetSeatCategoryByID(ctx, int(req.SeatCategoryId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetSeatCategoryByIDResponse{
		SeatCategory: &user_admin.SeatCategory{
			Id:               int32(seatCategory.ID),
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	}, nil
}

func (h *GrpcHandler) GetSeatCategoryByName(ctx context.Context, req *user_admin.GetSeatCategoryByNameRequest) (*user_admin.GetSeatCategoryByNameResponse, error) {
	seatCategory, err := h.svc.GetSeatCategoryByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &user_admin.GetSeatCategoryByNameResponse{
		SeatCategory: &user_admin.SeatCategory{
			Id:               int32(seatCategory.ID),
			SeatCategoryName: seatCategory.SeatCategoryName,
		},
	}, nil
}

func (h *GrpcHandler) UpdateSeatCategory(ctx context.Context, req *user_admin.UpdateSeatCategoryRequest) (*user_admin.UpdateSeatCategoryResponse, error) {
	err := h.svc.UpdateSeatCategory(ctx, int(req.Id), SeatCategory{
		SeatCategoryName: req.SeatCategory.SeatCategoryName,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.UpdateSeatCategoryResponse{
		Message: "successfully updated seat category",
	}, nil
}

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

// User
func (h *GrpcHandler) ListAllUser(ctx context.Context, req *user_admin.ListAllUserRequest) (*user_admin.ListAllUserResponse, error) {
	users, err := h.svc.ListAllUser(ctx)
	if err != nil {
		return nil, err
	}
	response := []*user_admin.User{}

	for _, res := range users {
		user := user_admin.User{
			Id:          int32(res.ID),
			Username:    res.Username,
			Phone:       res.PhoneNumber,
			Email:       res.Email,
			FirstName:   res.FirstName,
			LastName:    res.LastName,
			Gender:      res.Gender,
			DateOfBirth: res.DateOfBirth,
			IsVerified:  res.IsVerified,
		}
		response = append(response, &user)
	}
	return &user_admin.ListAllUserResponse{
		User: response,
	}, nil
}

func (h *GrpcHandler) GetUserByID(ctx context.Context, req *user_admin.GetUserByIdRequest) (*user_admin.GetUserByIdResponse, error) {
	user, err := h.svc.GetUserByID(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &user_admin.GetUserByIdResponse{
		User: &user_admin.User{
			Id:          int32(user.ID),
			Username:    user.Username,
			Phone:       user.PhoneNumber,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Gender:      user.Gender,
			DateOfBirth: user.DateOfBirth,
			IsVerified:  user.IsVerified,
		},
	}, nil
}

func (h *GrpcHandler) BlockUser(ctx context.Context, req *user_admin.BlockUserRequest) (*user_admin.BlockUserResponse, error) {
	err := h.svc.BlockUser(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *GrpcHandler) UnBlockUser(ctx context.Context, req *user_admin.UnBlockUserRequest) (*user_admin.UnBlockUserResponse, error) {
	err := h.svc.UnBlockUser(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
