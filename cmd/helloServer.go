package main

import (
	"fmt"
	"github.com/sniperCore/app/uploader/grpc/pb"
	"github.com/sniperCore/app/uploader/grpc/server"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                              // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server.HelloServer{}) // 在gRPC服务端注册hello服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
