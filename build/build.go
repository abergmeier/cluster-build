package build

import (
	"fmt"
	"github.com/abergmeier/cluster-build/operation"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
)

type Builds struct {
	Create   chan CreateBuildRequest
	Cancel   chan CancelBuildRequest
	Get      chan GetBuildRequest
	List     chan ListBuildsRequest
	builds   map[string]*cloudbuild.Build
	ops      *operation.Operations
	latestId uint64
}

type CancelBuildRequest struct {
	R *cloudbuild.CancelBuildRequest
	C chan cloudbuild.Build
}

type CreateBuildRequest struct {
	R *cloudbuild.CreateBuildRequest
	C chan longrunning.Operation
}

type GetBuildRequest struct {
	R *cloudbuild.GetBuildRequest
	C chan cloudbuild.Build
}

type ListBuildsRequest struct {
	R *cloudbuild.ListBuildsRequest
	C chan cloudbuild.ListBuildsResponse
}

func NewBuilds(ops *operation.Operations) (*Builds, error) {
	b := &Builds{
		builds: make(map[string]*cloudbuild.Build),
		ops:    ops,
	}
	go b.actor()
	return b, nil
}

func (b *Builds) newBuildId(build *cloudbuild.Build) string {
	b.latestId++
	return fmt.Sprintf("%s/%u", build.ProjectId, b.latestId)
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
		case g := <-b.Get:
			b.get(&g)
		case l := <-b.List:
			b.list(&l)
		}
	}

}

func (b *Builds) cancel(c *CancelBuildRequest) {
	gb, ok := b.builds[c.R.Id]
	if !ok {
		panic("Not ok")
	}

	panic("Implement to operator")
	c.C <- *gb
}

func (b *Builds) create(r *CreateBuildRequest) {

	bId := b.newBuildId(r.R.Build)
	b.builds[bId] = r.R.Build
	opName := fmt.Sprintf("CreateBuild%s", bId)
	b.ops.Create <- operation.CreateOperation{
		Name: opName,
		Call: func(cancel <-chan bool) {

		},
	}
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
	gb, ok := b.builds[g.R.Id]
	if !ok {
		panic("Not ok")
	}
	g.C <- *gb
}

func (b *Builds) list(l *ListBuildsRequest) {

	builds := make([]*cloudbuild.Build, len(b.builds), len(b.builds))

	for _, v := range b.builds {
		builds = append(builds, v)
	}

	l.C <- cloudbuild.ListBuildsResponse{
		Builds: builds,
	}
}
