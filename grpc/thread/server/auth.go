package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ClientNameKey    = "client_name"
	ClientSecrectKey = "client_secret"
)

func Authentication(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "no data in context")
	}

	cnk := md.Get(ClientNameKey)
	if len(cnk) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: empty name")
	}

	csk := md.Get(ClientSecrectKey)
	if len(csk) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: empty secret")
	}

	if cnk[0] != "Barry" || csk[0] != "Barry" {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed")
	}

	return handler(ctx, req)
}
