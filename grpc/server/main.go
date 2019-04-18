package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"yuan.wang/mgpg"
	pb "yuan.wang/mgpg/pb"
)

const (
	port = ":8005"
)

func main() {

	errChan := make(chan error)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		errChan <- err
		return
	}

	mgpgSrv := mgpg.NewMGpgService()
	mgpgEndpoint := mgpg.MakeGenerateKeyEndpoint(mgpgSrv)
	srv := mgpg.MakeGRPCServer(context.Background(), mgpgEndpoint)

	grpcSrv := grpc.NewServer()
	pb.RegisterMGpgServer(grpcSrv, srv)

	errChan <- grpcSrv.Serve(lis)
}
