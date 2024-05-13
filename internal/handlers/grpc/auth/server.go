package authgrpc

import (
	"context"

	authv1 "github.com/k3k13a1/rm_opts/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	AddPhone(
		ctx context.Context,
		id int64,
		phone string,
	) (err error)
	EditPhone(
		ctx context.Context,
		id int64,
		phone string,
	) (err error)
	AddEmail(
		ctx context.Context,
		id int64,
		email string,
	) (err error)
	EditEmail(
		ctx context.Context,
		id int64,
		email string,
	) (err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServiceServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) AddPhone(ctx context.Context, req *authv1.AddPhoneRequest) (*authv1.AddPhoneResponse, error) {
	if req.GetPhone() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.auth.AddPhone(ctx, int64(req.GetId()), req.GetPhone())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add phone")
	}

	return &authv1.AddPhoneResponse{}, nil
}

func (s *serverAPI) EditPhone(ctx context.Context, req *authv1.EditPhoneRequest) (*authv1.EditPhoneResponse, error) {
	if req.GetPhone() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.auth.EditPhone(ctx, int64(req.GetId()), req.GetPhone())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to edit phone")
	}

	return &authv1.EditPhoneResponse{}, nil
}

func (s *serverAPI) AddEmail(ctx context.Context, req *authv1.AddEmailRequest) (*authv1.AddEmailResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.auth.AddEmail(ctx, int64(req.GetId()), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to add email")
	}

	return &authv1.AddEmailResponse{}, nil
}

func (s *serverAPI) EditEmail(ctx context.Context, req *authv1.EditEmailRequest) (*authv1.EditEmailResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.auth.EditEmail(ctx, int64(req.GetId()), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to edit email")
	}

	return &authv1.EditEmailResponse{}, nil
}
