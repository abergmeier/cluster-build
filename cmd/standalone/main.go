package main

import (
	"fmt"
	"github.com/abergmeier/cluster-build/operation"
	"github.com/abergmeier/cluster-build/server"
	"google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ops, err := operation.NewOperationWorker()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	cbs, err := server.NewCloudBuild(ops)
	if err != nil {
		panic(err)
	}
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
