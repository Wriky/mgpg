package main

import (
	"fmt"
	"os"
	"time"

	"earlydata.com/waterdrop/microservice/gpg/pb"
	"github.com/go-kit/kit/endpoint"

	"yuan.wang/mgpg"

	"google.golang.org/grpc"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

const (
	port = ":8005"
)

// NewGRPCClientService : client service
func NewGRPCClientService(conn *grpc.ClientConn) mgpg.Service {
	var generateKeyEndpoint endpoint.Endpoint
	{
		generateKeyEndpoint = grpctransport.NewClient(conn, "MGpg.MGpg", "GenerateKey",
			mgpg.EncodeGRPCGenerateKeyRequest, mgpg.DecodeGRPCGenerateKeyResponse, pb.GenerateKeyResponse{}).Endpoint()
	}
	return mgpg.MakeGRPCServer(context, generateKeyEndpoint)

}

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	defer conn.Close()

}
