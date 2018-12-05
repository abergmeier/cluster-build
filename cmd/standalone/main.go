package main

import (
	"fmt"
	"github.com/abergmeier/cluster-build/build"
	"github.com/abergmeier/cluster-build/operation"
	"github.com/abergmeier/cluster-build/server"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		panic(errors.Wrap(err, "Failed to listen"))
	}

	bs, err := build.NewBuilds()
	if err != nil {
		panic(errors.Wrap(err, "Cannot create Builds"))
	}

	cbs, err := server.NewCloudBuild(bs)
	if err != nil {
		panic(errors.Wrap(err, "Cannot create CloudBuild instance"))
	}

	ops, err := operation.NewOperationWorker(bs)
	if err != nil {
		panic(errors.Wrap(err, "Cannot create Operation instance"))
	}
	defer ops.Close()

	grpcServer := grpc.NewServer()
	defer cbs.Close()
	cloudbuild.RegisterCloudBuildServer(grpcServer, cbs)
	os, err := server.NewOperations(ops)
	if err != nil {
		panic(err)
	}
	defer os.Close()
	longrunning.RegisterOperationsServer(grpcServer, os)
	grpcServer.Serve(lis)
}
