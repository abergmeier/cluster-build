package server

import (
	"context"

	"github.com/abergmeier/cluster-build/operation"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/longrunning"
)

type OperationsServer struct {
	ops *operation.Operations
}

func NewOperations(ops *operation.Operations) (*OperationsServer, error) {
	o := &OperationsServer{
		ops: ops,
	}
	return o, nil
}

func (s *OperationsServer) Close() {
	// TODO: Implement
}

func (s *OperationsServer) ListOperations(ctx context.Context, r *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error) {
	c := make(chan operation.ListOperationsResponse)
	s.ops.List <- operation.ListOperationsRequest{
		R: r,
		C: c,
	}
	resp := <-c
	return &resp.Value, resp.Err
}

func (s *OperationsServer) GetOperation(ctx context.Context, r *longrunning.GetOperationRequest) (*longrunning.Operation, error) {
	c := make(chan operation.GetOperationResponse)
	s.ops.Get <- operation.GetOperationRequest{
		R: r,
		C: c,
	}
	resp := <-c
	return &resp.Value, resp.Err
}

func (s *OperationsServer) DeleteOperation(ctx context.Context, r *longrunning.DeleteOperationRequest) (*empty.Empty, error) {
	c := make(chan error)
	s.ops.Delete <- operation.DeleteOperationRequest{
		R: r,
		C: c,
	}
	err := <-c
	return nil, err
}

// On successful cancellation,
// the operation is not deleted; instead, it becomes an operation with
// an [Operation.error][google.longrunning.Operation.error] value with a [google.rpc.Status.code][google.rpc.Status.code] of 1,
// corresponding to `Code.CANCELLED`.
func (s *OperationsServer) CancelOperation(ctx context.Context, r *longrunning.CancelOperationRequest) (*empty.Empty, error) {
	c := make(chan error)
	s.ops.Cancel <- operation.CancelOperationRequest{
		R: r,
		C: c,
	}
	err := <-c
	return nil, err
}
