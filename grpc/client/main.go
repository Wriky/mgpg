package main

import (
	"context"
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
	return mgpg.Endpoints{
		GenerateKeyEndpoint: generateKeyEndpoint,
	}
}

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	service := NewGRPCClientService(conn)
	var result interface{}
	result, err = service.GenerateKey(context.Background(), &pb.GenerateKeyRequest{
		Name:    "yuan.wang",
		Comment: "haha",
		Email:   "yuan.wang@analyticservice.net",
		Expiry:  50000000,
		Armor:   false,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp := result.(*pb.GenerateKeyResponse)
	fmt.Printf("%x\n\n\n", resp.Response.GetPub())
	fmt.Printf("%x", resp.Response.GetSec())

}
