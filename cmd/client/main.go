package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	v1 "github.com/jmuk/fib-grpc/gen/fib/v1"
	"github.com/jmuk/fib-grpc/gen/fib/v1/fibv1connect"
	"github.com/jmuk/fib-grpc/gen/google/longrunning/longrunningconnect"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	n := flag.Int64("n", 0, "")

	flag.Parse()
	fc := fibv1connect.NewFibServiceClient(http.DefaultClient, "http://localhost:8080/")
	lc := longrunningconnect.NewOperationsClient(http.DefaultClient, "http://localhost:8080/")

	resp, err := fc.Fib(context.Background(), connect.NewRequest(&v1.FibRequest{N: *n}))
	if err != nil {
		panic(err)
	}

	tick := time.NewTicker(time.Second)
	for range tick.C {
		log.Printf("Checking the status...")
		op, err := lc.GetOperation(context.Background(), connect.NewRequest(&longrunning.GetOperationRequest{
			Name: resp.Msg.Name,
		}))
		if err != nil {
			panic(err)
		}
		if op.Msg.Done {
			nv := &structpb.Value{}
			anypb.UnmarshalTo(op.Msg.GetResponse(), nv, proto.UnmarshalOptions{})
			log.Printf("Finished!  The answer is %f", nv.GetNumberValue())
			break
		}
	}
}
