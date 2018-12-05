package server

import (
	"context"
	"fmt"
	"github.com/abergmeier/cluster-build/build"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
)

type CloudBuildServer struct {
	builds *build.Builds

	operationC chan string
}

func NewCloudBuild(bs *build.Builds) (*CloudBuildServer, error) {
	return &CloudBuildServer{
		builds: bs,
	}, nil
}

func (s *CloudBuildServer) Close() {
	s.builds.Close()
}

func (s *CloudBuildServer) CreateBuild(ctx context.Context, r *cloudbuild.CreateBuildRequest) (*longrunning.Operation, error) {
	fmt.Println(r.ProjectId)

	createdBuild := make(chan longrunning.Operation)

	s.builds.Create <- build.CreateBuildRequest{
		R: r,
		C: createdBuild,
	}

	op := <-createdBuild
	return &op, nil
}

func (s *CloudBuildServer) GetBuild(ctx context.Context, r *cloudbuild.GetBuildRequest) (*cloudbuild.Build, error) {
	gottenBuild := make(chan build.GetBuildResponse)
	s.builds.Get <- build.GetBuildRequest{
		Id: r.Id,
		C:  gottenBuild,
	}

	b := <-gottenBuild
	return &b.Build.Build, b.Err
}

func (s *CloudBuildServer) ListBuilds(ctx context.Context, r *cloudbuild.ListBuildsRequest) (*cloudbuild.ListBuildsResponse, error) {
	listedBuilds := make(chan build.Build)

	s.builds.List <- build.ListBuildsRequest{
		R: r,
		C: listedBuilds,
	}

	builds := []*cloudbuild.Build{}
	go func() {
		for lb := range listedBuilds {
			builds = append(builds, &lb.Build)
		}
	}()

	resp := &cloudbuild.ListBuildsResponse{
		Builds: builds,
	}

	return resp, nil
}

func (s *CloudBuildServer) CancelBuild(ctx context.Context, r *cloudbuild.CancelBuildRequest) (*cloudbuild.Build, error) {
	recentState := make(chan build.CancelBuildResponse)

	s.builds.Cancel <- build.CancelBuildRequest{
		Id: r.Id,
		C:  recentState,
	}

	resp := <-recentState
	return &resp.Build.Build, resp.Err
}

func (s *CloudBuildServer) RetryBuild(ctx context.Context, r *cloudbuild.RetryBuildRequest) (*longrunning.Operation, error) {
	panic("RETRY")
}

func (s *CloudBuildServer) CreateBuildTrigger(ctx context.Context, r *cloudbuild.CreateBuildTriggerRequest) (*cloudbuild.BuildTrigger, error) {
	panic("BUILDTRIGGER")
}

func (s *CloudBuildServer) GetBuildTrigger(ctx context.Context, r *cloudbuild.GetBuildTriggerRequest) (*cloudbuild.BuildTrigger, error) {
	panic("GETBUILDTRIGGER")
}

func (s *CloudBuildServer) ListBuildTriggers(ctx context.Context, r *cloudbuild.ListBuildTriggersRequest) (*cloudbuild.ListBuildTriggersResponse, error) {
	panic("LISTBUILDTRIGGER")
}

func (s *CloudBuildServer) DeleteBuildTrigger(ctx context.Context, r *cloudbuild.DeleteBuildTriggerRequest) (*empty.Empty, error) {
	panic("DELETEBUILDTRIGGER")
}

func (s *CloudBuildServer) UpdateBuildTrigger(ctx context.Context, r *cloudbuild.UpdateBuildTriggerRequest) (*cloudbuild.BuildTrigger, error) {
	panic("UPDPATEBUILDTRIGGER")
}

func (s *CloudBuildServer) RunBuildTrigger(ctx context.Context, r *cloudbuild.RunBuildTriggerRequest) (*longrunning.Operation, error) {
	panic("RUNBUILDTRIGGER")
}
