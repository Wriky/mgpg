package mgpg

import (
	"context"
	"time"

	"earlydata.com/xxg/microservice/generics"
	"github.com/alokmenghrajani/gpgeez"
	"github.com/go-kit/kit/endpoint"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	pb "yuan.wang/mgpg/pb"
)

/* mgpg */

// Service : mpgp service
type Service interface {
	GenerateKey(ctx context.Context, req interface{}) (interface{}, error)
}

type mgpgService struct{}

// NewMGpgService : init
func NewMGpgService() Service {
	return mgpgService{}
}

func (s mgpgService) GenerateKey(ctx context.Context, req interface{}) (interface{}, error) {
	request := req.(*pb.GenerateKeyRequest)
	config := &gpgeez.Config{Expiry: time.Duration(request.GetExpiry()) * time.Hour}
	key, err := gpgeez.CreateKey(request.Name, request.Comment, request.Email, config)
	if err != nil {
		return nil, generics.Error{ErrorCode: 1001, ErrorMessage: err.Error()}
	}

	reply := pb.GenerateKeyReply{}

	if request.Armor {
		pub, err := key.Armor()
		if err != nil {
			return nil, generics.Error{ErrorCode: 1002, ErrorMessage: err.Error()}
		}
		reply.Pub = []byte(pub)

		sec, err := key.ArmorPrivate(config)
		if err != nil {
			return nil, generics.Error{ErrorCode: 1003, ErrorMessage: err.Error()}
		}
		reply.Sec = []byte(sec)
	} else {
		reply.Pub = key.Keyring()
		reply.Sec = key.Secring(config)
	}
	return &pb.GenerateKeyResponse{Response: &reply}, nil
}

// MakeGenerateKeyEndpoint : transform Endpoint
func MakeGenerateKeyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GenerateKeyRequest)
		return s.GenerateKey(ctx, req)
	}
}

/* grpc */

type grpcServer struct {
	generateKey grpctransport.Handler
}

func (s *grpcServer) GenerateKey(ctx context.Context, req *pb.GenerateKeyRequest) (*pb.GenerateKeyResponse, error) {
	_, rep, err := s.generateKey.ServeGRPC(ctx, req)
	if err != nil {
		if genericErr, ok := err.(generics.Error); ok {
			return &pb.GenerateKeyResponse{
				ErrorCode:    int64(genericErr.ErrorCode),
				ErrorMessage: genericErr.ErrorMessage,
			}, nil
		}
	}
	return rep.(*pb.GenerateKeyResponse), nil
}

// MakeGRPCServer : grpc server
func MakeGRPCServer(ctx context.Context, grpcEndpoint endpoint.Endpoint) pb.MGpgServer {
	return &grpcServer{
		generateKey: grpctransport.NewServer(
			grpcEndpoint,
			DecodeGRPCGenerateKeyRequest,
			EncodeGRPCGenerateKeyResponse,
		),
	}
}

// DecodeGRPCGenerateKeyRequest : decode request
func DecodeGRPCGenerateKeyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

// EncodeGRPCGenerateKeyResponse : encode response
func EncodeGRPCGenerateKeyResponse(_ context.Context, grpcRep interface{}) (interface{}, error) {
	return grpcRep, nil
}

// EncodeGRPCGenerateKeyRequest : encode request
func EncodeGRPCGenerateKeyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

// DecodeGRPCGenerateKeyResponse : decode response
func DecodeGRPCGenerateKeyResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenerateKeyResponse)
	if reply.ErrorCode != 0 {
		genericErr := generics.Error{
			ErrorCode:    int(reply.ErrorCode),
			ErrorMessage: reply.ErrorMessage,
		}
		return genericErr, nil
	}
	return reply, nil
}
