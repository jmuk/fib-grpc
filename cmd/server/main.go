package main

import (
	"net/http"

	"github.com/jmuk/fib-grpc/gen/fib/v1/fibv1connect"
	"github.com/jmuk/fib-grpc/gen/google/longrunning/longrunningconnect"
	"github.com/jmuk/fib-grpc/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	// The generated constructors return a path and a plain net/http
	// handler.
	s := server.NewServer()
	mux.Handle(longrunningconnect.NewOperationsHandler(s))
	mux.Handle(fibv1connect.NewFibServiceHandler(s))
	http.ListenAndServe(
		"localhost:8080",
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
