package operation

import (
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"log"
)

type CancelOperationRequest struct {
	R *longrunning.CancelOperationRequest
	C chan error
}

type CancelledOperation struct {
	Name         string
	ErrorMessage string
}

type CreateOperation struct {
	Name string
	Call func(<-chan bool)
}

type DeleteOperationRequest struct {
	R *longrunning.DeleteOperationRequest
	C chan error
}

type GetOperationRequest struct {
	R *longrunning.GetOperationRequest
	C chan GetOperationResponse
}

type GetOperationResponse struct {
	Value longrunning.Operation
	Err   error
}

type ListOperationsRequest struct {
	R *longrunning.ListOperationsRequest
	C chan ListOperationsResponse
}

type ListOperationsResponse struct {
	Value longrunning.ListOperationsResponse
	Err   error
}

type Operation struct {
	longrunning.Operation
	cancel chan bool
}

type Operations struct {
	Create chan CreateOperation
	Cancel chan CancelOperationRequest
	Delete chan DeleteOperationRequest
	Get    chan GetOperationRequest
	List   chan ListOperationsRequest

	cancelled chan CancelledOperation

	ops map[string]*Operation
}

func NewOperationWorker() (*Operations, error) {
	o := &Operations{}
	go o.worker()
	return o, nil
}

func (o *Operations) Close() {
	// TODO: Implement
}

func (o *Operations) cancel(c *CancelOperationRequest) {
	op, ok := o.ops[c.R.Name]
	if !ok {
		log.Fatalf("Could not find Operation %s", c.R.Name)
	}
	op.cancel <- true
}

func (o *Operations) setCancel(c *CancelledOperation) {
	op, ok := o.ops[c.Name]
	if !ok {
		log.Fatalf("Could not find Operation %s", c.Name)
	}
	// TODO: check Done
	op.Done = true
	op.Result = &longrunning.Operation_Error{
		Error: &status.Status{
			Code:    int32(code.Code_CANCELLED),
			Message: c.ErrorMessage,
		},
	}
}

func (o *Operations) create(c *CreateOperation) {
	cancel := make(chan bool)
	o.ops[c.Name] = &Operation{
		cancel: cancel,
	}
	go c.Call(cancel)
}

func (o *Operations) delete(d *DeleteOperationRequest) {
	// TODO: Perhaps cancel first
	delete(o.ops, d.R.Name)
	d.C <- nil
}

func (o *Operations) get(g *GetOperationRequest) {
	op, ok := o.ops[g.R.Name]
	if !ok {
		log.Fatalf("Could not find Operation %s", g.R.Name)
	}
	g.C <- GetOperationResponse{
		Value: op.Operation,
		Err:   nil,
	}
}

func (o *Operations) list(l *ListOperationsRequest) {
	// TODO: Handle Names and Filters
	opList := make([]*longrunning.Operation, len(o.ops), len(o.ops))
	for _, v := range o.ops {
		opList = append(opList, &v.Operation)
	}
	l.C <- ListOperationsResponse{
		Value: longrunning.ListOperationsResponse{
			Operations: opList,
		},
		Err: nil,
	}
}

func (o *Operations) worker() {

	for {
		select {
		case c := <-o.Cancel:
			o.cancel(&c)
		case c := <-o.Create:
			o.create(&c)
		case c := <-o.cancelled:
			o.setCancel(&c)
		case d := <-o.Delete:
			o.delete(&d)
		case g := <-o.Get:
			o.get(&g)
		case l := <-o.List:
			o.list(&l)
		}
	}
}
