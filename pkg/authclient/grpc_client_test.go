package authclient

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	authpbv1 "github.com/williabk198/timeclock/pkg/authpb/v1"
)

const bufSize = 1024 * 1024

// fakeIntrospectionServer is a minimal stand-in for the auth service's gRPC
// TokenIntrospection service, used to test GRPCClient without a real server.
type fakeIntrospectionServer struct {
	authpbv1.UnimplementedTokenIntrospectionServer
	resp *authpbv1.IntrospectResponse
	err  error
}

func (f fakeIntrospectionServer) Introspect(ctx context.Context, req *authpbv1.IntrospectRequest) (*authpbv1.IntrospectResponse, error) {
	return f.resp, f.err
}

// newTestGRPCClient starts the given fake server over an in-memory bufconn
// listener and returns a GRPCClient connected to it.
func newTestGRPCClient(t *testing.T, fake authpbv1.TokenIntrospectionServer) *GRPCClient {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	authpbv1.RegisterTokenIntrospectionServer(srv, fake)
	go func() {
		_ = srv.Serve(lis)
	}()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	return &GRPCClient{conn: conn, client: authpbv1.NewTokenIntrospectionClient(conn)}
}

func TestGRPCClient_Introspect(t *testing.T) {
	type args struct {
		token string
	}

	tests := []struct {
		name      string
		fake      authpbv1.TokenIntrospectionServer
		args      args
		want      Introspection
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Active",
			fake: fakeIntrospectionServer{resp: &authpbv1.IntrospectResponse{
				Active:    true,
				Scope:     "default user:admin",
				TokenType: "access",
				Sub:       "user-1",
				Iat:       100,
				Exp:       200,
			}},
			args: args{token: "valid-token"},
			want: Introspection{
				Active:    true,
				Scope:     "default user:admin",
				TokenType: "access",
				Subject:   "user-1",
				IssuedAt:  100,
				ExpiresAt: 200,
			},
			assertion: assert.NoError,
		},
		{
			name:      "Inactive",
			fake:      fakeIntrospectionServer{resp: &authpbv1.IntrospectResponse{Active: false}},
			args:      args{token: "bad-token"},
			want:      Introspection{Active: false},
			assertion: assert.NoError,
		},
		{
			name:      "ServerError",
			fake:      fakeIntrospectionServer{err: errors.New("boom")},
			args:      args{token: "whatever"},
			want:      Introspection{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestGRPCClient(t, tt.fake)

			got, err := client.Introspect(context.Background(), tt.args.token)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewGRPCClient(t *testing.T) {
	client, err := NewGRPCClient("localhost:0")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NoError(t, client.Close())
}
