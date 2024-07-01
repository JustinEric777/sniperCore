package server

import (
	"context"
	"github.com/sniperCore/app/uploader/grpc/pb"
)

type HelloServer struct {
	pb.UnimplementedGreeterServer
}

func (s *HelloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
