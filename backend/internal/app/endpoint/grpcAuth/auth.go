package grpcAuth

import (
	"CoggersProject/internal/app/errorz"
	pb "CoggersProject/pkg/protos/auth"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	UserLogin(string, string) (string, error)
	NewUserRegistration(string, string) (string, error)
}

type Endpoint struct {
	Auth      Auth
	SecretKey string
	pb.UnimplementedAuthServiceServer
}

func (e *Endpoint) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username пустое значение")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password пустое значение")
	}

	token, err := e.Auth.UserLogin(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errorz.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "неверный username или password")
		}

		return nil, status.Error(codes.Internal, "ошибка аутентификации")
	}

	return &pb.AuthResponse{Key: token}, nil
}

func (e *Endpoint) Registration(ctx context.Context, req *pb.RegistrationRequest) (*pb.AuthResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "username пустое значение")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password пустое значение")
	}

	token, err := e.Auth.NewUserRegistration(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, errorz.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, "пользователь с таким именем уже существует")
		}

		return nil, status.Error(codes.Internal, "ошибка регистрации")
	}

	return &pb.AuthResponse{Key: token}, nil
}
