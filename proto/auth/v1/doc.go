// package authv1 contains the protobuf source definitions for the auth service's gRPC
// API. It holds no Go code of its own; running `go generate` here (or `make generate`
// from the repo root) regenerates the Go bindings into pkg/authpb/v1 from
// introspection.proto.
//
// Requires the protoc compiler plus the protoc-gen-go and protoc-gen-go-grpc plugins.
// Run `make tools` to install the Go-based plugins (protoc itself must be installed
// separately, since it isn't a Go module).
package authv1

//go:generate protoc --proto_path=. --go_out=../../.. --go_opt=module=github.com/williabk198/timeclock --go-grpc_out=../../.. --go-grpc_opt=module=github.com/williabk198/timeclock introspection.proto
