package operation

import (
	"github.com/abergmeier/cluster-build/build"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/longrunning"
)

type CancelOperationRequest struct {
	Name string
	C    chan error
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
	Data interface{}
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

	builds *build.Builds
}

func NewOperationWorker(builds *build.Builds) (*Operations, error) {
	o := &Operations{
		builds: builds,
	}
	go o.actor()
	return o, nil
}

func (o *Operations) Close() {
	// TODO: Implement
}

func (o *Operations) cancel(c *CancelOperationRequest) {
	defer close(c.C)
	buildId, err := build.ExtractIdFromOperationName(c.Name)
	if err != nil {
		c.C <- errors.Wrapf(err, "Operation name %s not recognized", c.Name)
		return
	}
	lastState := make(chan build.CancelBuildResponse)
	o.builds.Cancel <- build.CancelBuildRequest{
		Id: buildId,
		C:  lastState,
	}
	resp := <-lastState
	c.C <- resp.Err
}

func (o *Operations) delete(d *DeleteOperationRequest) {
	defer close(d.C)
	buildId, err := build.ExtractIdFromOperationName(d.R.Name)
	if err != nil {
		d.C <- errors.Wrapf(err, "Operation name %s not recognized", d.R.Name)
		return
	}
	result := make(chan error)
	o.builds.Delete <- build.DeleteBuildRequest{
		Id: buildId,
		C:  result,
	}

	err = <-result
	d.C <- err
}

func (o *Operations) get(g *GetOperationRequest) {
	defer close(g.C)
	buildId, err := build.ExtractIdFromOperationName(g.R.Name)
	if err != nil {
		g.C <- GetOperationResponse{
			Err: errors.Wrapf(err, "Operation name %s not recognized", g.R.Name),
		}
		return
	}
	gotten := make(chan build.GetBuildResponse)
	o.builds.Get <- build.GetBuildRequest{
		Id: buildId,
		C:  gotten,
	}
	resp := <-gotten
	g.C <- GetOperationResponse{
		Value: resp.Build.Operation,
		Err:   resp.Err,
	}
}

func (o *Operations) list(l *ListOperationsRequest) {
	defer close(l.C)
	// TODO: Handle Names and Filters
	opList := []*longrunning.Operation{}
	listed := make(chan build.Build)
	go func() {
		for li := range listed {
			opList = append(opList, &li.Operation)
		}

		l.C <- ListOperationsResponse{
			Value: longrunning.ListOperationsResponse{
				Operations: opList,
			},
			Err: nil,
		}
	}()
	o.builds.List <- build.ListBuildsRequest{
		C: listed,
	}
}

func (o *Operations) actor() {

	for {
		select {
		case c := <-o.Cancel:
			o.cancel(&c)
		case d := <-o.Delete:
			o.delete(&d)
		case g := <-o.Get:
			o.get(&g)
		case l := <-o.List:
			o.list(&l)
		}
	}
}
