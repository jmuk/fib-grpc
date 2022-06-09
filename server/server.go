package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/bufbuild/connect-go"
	v1 "github.com/jmuk/fib-grpc/gen/fib/v1"
	"github.com/jmuk/fib-grpc/gen/fib/v1/fibv1connect"
	"github.com/jmuk/fib-grpc/gen/google/longrunning/longrunningconnect"
	"google.golang.org/genproto/googleapis/longrunning"
)

type Server struct {
	fibv1connect.UnimplementedFibServiceHandler
	longrunningconnect.UnimplementedOperationsHandler

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

func (s *Server) Fib(ctx context.Context, req *connect.Request[v1.FibRequest]) (*connect.Response[v1.FibResponse], error) {
	newName := fmt.Sprintf("%x", rand.Uint64())
	s.ops.Store(newName, newOperation(func() int64 {
		return fib(req.Msg.N)
	}))
	return connect.NewResponse(&v1.FibResponse{Name: newName}), nil
}

func (s *Server) GetOperation(ctx context.Context, req *connect.Request[longrunning.GetOperationRequest]) (*connect.Response[longrunning.Operation], error) {
	got, ok := s.ops.Load(req.Msg.Name)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("not found"))
	}
	return connect.NewResponse(got.(*operation).ToOperation(req.Msg.Name)), nil
}
