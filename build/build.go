package build

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"strings"
)

const (
	operationPrefix = "Build/"
)

type Build struct {
	cloudbuild.Build
	longrunning.Operation
	cancel chan bool
}

type Builds struct {
	Create   chan CreateBuildRequest
	Cancel   chan CancelBuildRequest
	Delete   chan DeleteBuildRequest
	Get      chan GetBuildRequest
	List     chan ListBuildsRequest
	builds   map[string]Build
	latestId uint64
}

type CancelBuildRequest struct {
	Id string
	C  chan CancelBuildResponse
}

type CancelBuildResponse struct {
	Build Build
	Err   error
}

type CreateBuildRequest struct {
	R *cloudbuild.CreateBuildRequest
	C chan longrunning.Operation
}

type DeleteBuildRequest struct {
	Id string
	C  chan error
}

type GetBuildRequest struct {
	Id string
	C  chan GetBuildResponse
}

type GetBuildResponse struct {
	Build Build
	Err   error
}

type ListBuildsRequest struct {
	R *cloudbuild.ListBuildsRequest
	C chan Build
}

func NewBuilds() (*Builds, error) {
	b := &Builds{
		builds: make(map[string]Build),
	}
	go b.actor()
	return b, nil
}

func ExtractIdFromOperationName(operationName string) (string, error) {
	if !strings.HasPrefix(operationName, "Build/") {
		return "", errors.Errorf("%s missing Build/ prefix", operationName)
	}

	buildId := operationName[len("Build/"):]
	return buildId, nil
}

func (b *Builds) newBuildId(build *cloudbuild.Build) string {
	b.latestId++
	return fmt.Sprintf("%s/%u", build.ProjectId, b.latestId)
}

func getOperationName(build *cloudbuild.Build) string {
	return getOperationNameFromId(build.Id)
}

func getOperationNameFromId(buildId string) string {
	return fmt.Sprintf("Build/%s", buildId)
}

func (b *Builds) Close() {
	// TODO: Implement
}

func (b *Builds) actor() {
	for {
		select {
		case c := <-b.Cancel:
			b.cancel(&c)
		case c := <-b.Create:
			b.create(&c)
		case d := <-b.Delete:
			b.delete(&d)
		case g := <-b.Get:
			b.get(&g)
		case l := <-b.List:
			b.list(&l)
		}
	}

}

func (b *Builds) cancel(c *CancelBuildRequest) {
	defer close(c.C)
	cb, ok := b.builds[c.Id]
	if !ok {
		c.C <- CancelBuildResponse{
			Err: errors.Errorf("Build %s not found", c.Id),
		}
		return
	}
	cb.cancel <- true
	// For now we set this as a high level value. If later cancel processing
	// fails, it should update this field to internal error
	cb.Status = cloudbuild.Build_CANCELLED
	cb.Operation.Result = &longrunning.Operation_Error{
		Error: &status.Status{
			Code: int32(code.Code_CANCELLED),
		},
	}
	c.C <- CancelBuildResponse{
		Build: cb,
	}
}

func (b *Builds) create(r *CreateBuildRequest) {
	defer close(r.C)

	bId := b.newBuildId(r.R.Build)
	r.R.Build.Id = bId
	b.builds[bId] = Build{
		Build: *r.R.Build,
	}
	opName := getOperationName(r.R.Build)
	r.C <- longrunning.Operation{
		Name: opName,
		Done: false,
		Result: &longrunning.Operation_Response{
			Response: &any.Any{
				TypeUrl: "",
				Value:   []byte{},
			},
		},
	}
}

func (b *Builds) get(g *GetBuildRequest) {
	defer close(g.C)
	gb, ok := b.builds[g.Id]
	if !ok {
		g.C <- GetBuildResponse{
			Err: errors.Errorf("Build %s not found", g.Id),
		}
		return
	}
	g.C <- GetBuildResponse{
		Build: gb,
	}
}

func (b *Builds) delete(d *DeleteBuildRequest) {
	defer close(d.C)
	delete(b.builds, d.Id)
	d.C <- nil
}

func (b *Builds) list(l *ListBuildsRequest) {
	defer close(l.C)

	for _, v := range b.builds {
		l.C <- v
	}
}
