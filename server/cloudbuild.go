package server

import (
	"context"
	"fmt"
	"github.com/abergmeier/cluster-build/build"
	"github.com/abergmeier/cluster-build/operation"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
)

type CloudBuildServer struct {
	ops    *operation.Operations
	builds *build.Builds

	operationC chan string
}

func NewCloudBuild(ops *operation.Operations) (*CloudBuildServer, error) {
	bs, err := build.NewBuilds(ops)
	if err != nil {
		panic(err)
	}
	return &CloudBuildServer{
		ops:    ops,
		builds: bs,
	}, nil
}

func (s *CloudBuildServer) Close() {
	s.ops.Close()
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
	gottenBuild := make(chan cloudbuild.Build)
	s.builds.Get <- build.GetBuildRequest{
		R: r,
		C: gottenBuild,
	}

	b := <-gottenBuild
	return &b, nil
}

func (s *CloudBuildServer) ListBuilds(ctx context.Context, r *cloudbuild.ListBuildsRequest) (*cloudbuild.ListBuildsResponse, error) {
	listedBuilds := make(chan cloudbuild.ListBuildsResponse)

	s.builds.List <- build.ListBuildsRequest{
		R: r,
		C: listedBuilds,
	}

	list := <-listedBuilds
	return &list, nil
}

func (s *CloudBuildServer) CancelBuild(ctx context.Context, r *cloudbuild.CancelBuildRequest) (*cloudbuild.Build, error) {
	canceledBuild := make(chan cloudbuild.Build)

	s.builds.Cancel <- build.CancelBuildRequest{
		R: r,
		C: canceledBuild,
	}

	build := <-canceledBuild
	return &build, nil
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
