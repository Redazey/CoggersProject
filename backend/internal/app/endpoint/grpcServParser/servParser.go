package grpcServParser

import (
	"CoggersProject/config"
	pb "CoggersProject/pkg/protos/servParser"
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServParser interface {
	GetServersInfo() ([]config.Servers, error)
}

type Endpoint struct {
	Parser ServParser
	pb.UnimplementedServParserServiceServer
}

func (e *Endpoint) GetServersInfo(ctx context.Context, _ *empty.Empty) (*pb.ServParserResponse, error) {
	responseList := []*pb.ServerInfo{}
	serversInfo, err := e.Parser.GetServersInfo()

	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка получения информации о серверах")
	}
	for _, value := range serversInfo {
		responseList = append(responseList, &pb.ServerInfo{
			Adress:    value.Adress,
			Name:      value.Name,
			Version:   value.Version,
			MaxOnline: int64(value.MaxOnline),
			Online:    int64(value.Online),
		})
	}

	return &pb.ServParserResponse{ServersInfo: responseList}, nil
}
