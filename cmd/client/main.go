package main

import (
	"context"
	"flag"
	"log"
	"time"

	v1 "github.com/jmuk/fib-grpc/gen/fib/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	n := flag.Int64("n", 0, "")

	flag.Parse()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	fc := v1.NewFibServiceClient(conn)
	lc := longrunning.NewOperationsClient(conn)

	resp, err := fc.Fib(context.Background(), &v1.FibRequest{N: *n})
	if err != nil {
		panic(err)
	}

	tick := time.NewTicker(time.Second)
	for range tick.C {
		log.Printf("Checking the status...")
		op, err := lc.GetOperation(context.Background(), &longrunning.GetOperationRequest{
			Name: resp.Name,
		})
		if err != nil {
			panic(err)
		}
		if op.Done {
			nv := &structpb.Value{}
			anypb.UnmarshalTo(op.GetResponse(), nv, proto.UnmarshalOptions{})
			log.Printf("Finished!  The answer is %f", nv.GetNumberValue())
			break
		}
	}
}
