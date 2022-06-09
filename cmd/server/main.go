package main

import (
	"log"
	"net"

	v1 "github.com/jmuk/fib-grpc/gen/fib/v1"
	"github.com/jmuk/fib-grpc/server"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("starting a grpc server at port 8080")
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	gs := grpc.NewServer()
	s := server.NewServer()
	v1.RegisterFibServiceServer(gs, s)
	longrunning.RegisterOperationsServer(gs, s)
	gs.Serve(lis)
}
