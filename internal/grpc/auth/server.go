package grpcauth

import (
	"context"
	"errors"

	"github.com/PolyAbit/auth/internal/services/auth"
	authv1 "github.com/PolyAbit/protos/gen/go/auth"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateCredentials(in); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateCredentials(in); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.InvalidArgument, "user already exists")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, in *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if in.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userId is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, in.GetUserId())
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get permission")
	}

	return &authv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

type validableCredentials interface {
	GetEmail() string
	GetPassword() string
}

func validateCredentials(req validableCredentials) error {
	validate := validator.New()

	err := validate.Var(req.GetEmail(), "required")
	if err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	err = validate.Var(req.GetEmail(), "email")
	if err != nil {
		return status.Error(codes.InvalidArgument, "email not valid")
	}

	err = validate.Var(req.GetPassword(), "required")
	if err != nil {
		return status.Error(codes.InvalidArgument, "password not valid")
	}

	return nil
}
