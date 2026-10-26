package authclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpbv1 "github.com/williabk198/timeclock/pkg/authpb/v1"
)

// GRPCClient calls the auth service's gRPC TokenIntrospection service.
// It satisfies the same Introspector interface as Client, so it's a
// drop-in replacement wherever an Introspector is accepted.
type GRPCClient struct {
	conn   *grpc.ClientConn
	client authpbv1.TokenIntrospectionClient
}

// NewGRPCClient dials the auth service's gRPC endpoint at target (e.g.
// "auth:9090"). The returned GRPCClient owns the connection; call Close
// when it's no longer needed (e.g. on service shutdown).
//
// If no dial options are provided, the connection defaults to plaintext
// (insecure) transport credentials, which is only appropriate when both
// services run on a trusted local network. For any deployment where traffic
// between services could be intercepted, pass TLS-based credentials via opts
// (e.g. grpc.WithTransportCredentials(credentials.NewTLS(...))).
func NewGRPCClient(target string, opts ...grpc.DialOption) (*GRPCClient, error) {
	if len(opts) == 0 {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth gRPC service: %w", err)
	}

	return &GRPCClient{
		conn:   conn,
		client: authpbv1.NewTokenIntrospectionClient(conn),
	}, nil
}

func (c *GRPCClient) Introspect(ctx context.Context, token string) (Introspection, error) {
	resp, err := c.client.Introspect(ctx, &authpbv1.IntrospectRequest{Token: token})
	if err != nil {
		return Introspection{}, fmt.Errorf("introspection rpc failed: %w", err)
	}

	return Introspection{
		Active:    resp.GetActive(),
		Scope:     resp.GetScope(),
		TokenType: resp.GetTokenType(),
		Subject:   resp.GetSub(),
		IssuedAt:  resp.GetIat(),
		ExpiresAt: resp.GetExp(),
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}
