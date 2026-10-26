package transport

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/williabk198/timeclock/internal/services/auth/endpoints"
	authpbv1 "github.com/williabk198/timeclock/pkg/authpb/v1"
)

type introspectionServer struct {
	authpbv1.UnimplementedTokenIntrospectionServer
	authEndpoints endpoints.Endpoints
}

func (s introspectionServer) Introspect(ctx context.Context, req *authpbv1.IntrospectRequest) (*authpbv1.IntrospectResponse, error) {
	resp, err := s.authEndpoints.Introspect(ctx, endpoints.IntrospectRequestData{
		Token: req.GetToken(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to introspect token: %v", err)
	}
	if !resp.Active {
		return &authpbv1.IntrospectResponse{Active: false}, nil
	}

	return &authpbv1.IntrospectResponse{
		Active:    true,
		Scope:     resp.Scope,
		TokenType: resp.TokenType,
		Sub:       resp.Subject,
		Iat:       resp.IssuedAt,
		Exp:       resp.ExpiresAt,
	}, nil
}

// NewGRPCServer builds the gRPC server that exposes token introspection to
// other services.
//
// It's created with grpc.NewServer()'s defaults, which means it accepts
// plaintext (insecure) connections. That's fine when every service runs on a
// trusted local network, but any deployment where traffic between services
// could be intercepted should pass grpc.Creds(...) with TLS-based transport
// credentials instead.
func NewGRPCServer(authEndpoints endpoints.Endpoints) *grpc.Server {
	srv := grpc.NewServer()
	authpbv1.RegisterTokenIntrospectionServer(srv, introspectionServer{authEndpoints: authEndpoints})
	return srv
}
