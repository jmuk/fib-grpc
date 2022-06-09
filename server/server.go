package server

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	v1 "github.com/jmuk/fib-grpc/gen/fib/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	v1.UnimplementedFibServiceServer
	longrunning.UnimplementedOperationsServer

	ops sync.Map
}

func NewServer() *Server {
	return &Server{}
}

func fib(n int64) int64 {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func (s *Server) Fib(ctx context.Context, req *v1.FibRequest) (*v1.FibResponse, error) {
	newName := fmt.Sprintf("%x", rand.Uint64())
	s.ops.Store(newName, newOperation(func() int64 {
		return fib(req.N)
	}))
	return &v1.FibResponse{Name: newName}, nil
}

func (s *Server) GetOperation(ctx context.Context, req *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	got, ok := s.ops.Load(req.Name)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "not found")
	}

	return got.(*operation).ToOperation(req.Name), nil
}
