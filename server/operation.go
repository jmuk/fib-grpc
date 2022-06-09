package server

import (
	"sync"

	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type operation struct {
	sync.Mutex
	done   bool
	donec  chan struct{}
	result int64
}

func newOperation(task func() int64) *operation {
	donec := make(chan struct{})
	op := &operation{done: false, donec: donec}
	go func() {
		result := task()
		op.Lock()
		defer op.Unlock()
		op.result = result
		op.done = true
		close(op.donec)
	}()
	return op
}

func (op *operation) ToOperation(name string) *longrunning.Operation {
	op.Lock()
	defer op.Unlock()
	result := &longrunning.Operation{Name: name, Done: op.done}
	if op.done {
		val, _ := anypb.New(structpb.NewNumberValue(float64(op.result)))
		result.Result = &longrunning.Operation_Response{
			Response: val,
		}
	}
	return result
}
