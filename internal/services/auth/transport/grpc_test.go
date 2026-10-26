package transport

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/williabk198/timeclock/internal/services/auth/endpoints"
	authpbv1 "github.com/williabk198/timeclock/pkg/authpb/v1"
)

const bufSize = 1024 * 1024

type mockAuthEndpoints struct {
	mock.Mock
}

func (m *mockAuthEndpoints) Register(ctx context.Context, reqData endpoints.RegisterUserData) (endpoints.RegisterUserResponse, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.RegisterUserResponse), args.Error(1)
}

func (m *mockAuthEndpoints) Login(ctx context.Context, reqData endpoints.LoginRequestData) (endpoints.LoginResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.LoginResponseData), args.Error(1)
}

func (m *mockAuthEndpoints) Introspect(ctx context.Context, reqData endpoints.IntrospectRequestData) (endpoints.IntrospectResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.IntrospectResponseData), args.Error(1)
}

func (m *mockAuthEndpoints) Refresh(ctx context.Context, reqData endpoints.RefreshRequestData) (endpoints.LoginResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.LoginResponseData), args.Error(1)
}

func (m *mockAuthEndpoints) Revoke(ctx context.Context, reqData endpoints.RevokeRequestData) (endpoints.RevokeResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.RevokeResponseData), args.Error(1)
}

func (m *mockAuthEndpoints) GrantScope(ctx context.Context, reqData endpoints.GrantScopeRequestData) (endpoints.ScopeResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.ScopeResponseData), args.Error(1)
}

func (m *mockAuthEndpoints) RevokeScope(ctx context.Context, reqData endpoints.RevokeScopeRequestData) (endpoints.ScopeResponseData, error) {
	args := m.Called(ctx, reqData)
	return args.Get(0).(endpoints.ScopeResponseData), args.Error(1)
}

// startTestGRPCServer spins up the real gRPC server from NewGRPCServer over an
// in-memory bufconn listener, backed by the given mock, and returns a client
// connected to it.
func startTestGRPCServer(t *testing.T, authEndpoints endpoints.Endpoints) authpbv1.TokenIntrospectionClient {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := NewGRPCServer(authEndpoints)
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

	return authpbv1.NewTokenIntrospectionClient(conn)
}

func Test_introspectionServer_Introspect(t *testing.T) {
	type args struct {
		token string
	}

	tests := []struct {
		name          string
		setupMock     func(*mockAuthEndpoints)
		args          args
		wantActive    bool
		wantScope     string
		wantTokenType string
		wantSub       string
		wantIat       int64
		wantExp       int64
		wantErr       bool
	}{
		{
			name: "Active",
			setupMock: func(m *mockAuthEndpoints) {
				m.On("Introspect", mock.Anything, endpoints.IntrospectRequestData{Token: "valid-token"}).Return(endpoints.IntrospectResponseData{
					Active:    true,
					Scope:     "default user:admin",
					TokenType: "access",
					Subject:   "user-1",
					IssuedAt:  100,
					ExpiresAt: 200,
				}, error(nil))
			},
			args:          args{token: "valid-token"},
			wantActive:    true,
			wantScope:     "default user:admin",
			wantTokenType: "access",
			wantSub:       "user-1",
			wantIat:       100,
			wantExp:       200,
		},
		{
			name: "Inactive",
			setupMock: func(m *mockAuthEndpoints) {
				m.On("Introspect", mock.Anything, endpoints.IntrospectRequestData{Token: "bad-token"}).Return(endpoints.IntrospectResponseData{Active: false}, error(nil))
			},
			args:       args{token: "bad-token"},
			wantActive: false,
		},
		{
			name: "EndpointError",
			setupMock: func(m *mockAuthEndpoints) {
				m.On("Introspect", mock.Anything, endpoints.IntrospectRequestData{Token: "boom"}).Return(endpoints.IntrospectResponseData{}, assert.AnError)
			},
			args:    args{token: "boom"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEndpoints := &mockAuthEndpoints{}
			tt.setupMock(mockEndpoints)

			client := startTestGRPCServer(t, mockEndpoints)
			got, err := client.Introspect(context.Background(), &authpbv1.IntrospectRequest{Token: tt.args.token})

			if tt.wantErr {
				assert.Error(t, err)
				mockEndpoints.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantActive, got.GetActive())
			assert.Equal(t, tt.wantScope, got.GetScope())
			assert.Equal(t, tt.wantTokenType, got.GetTokenType())
			assert.Equal(t, tt.wantSub, got.GetSub())
			assert.Equal(t, tt.wantIat, got.GetIat())
			assert.Equal(t, tt.wantExp, got.GetExp())
			mockEndpoints.AssertExpectations(t)
		})
	}
}
