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

// movies
func (h *GrpcHandler) RegisterMovie(ctx context.Context, req *user_admin.RegisterMovieRequest) (*user_admin.RegisterMovieResponse, error) {
	movieId, err := h.svc.RegisterMovie(ctx, Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    int(req.Duration),
		Genre:       req.Genre,
		ReleaseDate: req.ReleaseDate,
		Rating:      float64(req.Rating),
		Language:    req.Language,
	})
	if err != nil {
		return nil, err
	}
	return &user_admin.RegisterMovieResponse{
		MovieId: uint32(movieId),
		Message: "create movie successfull",
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
