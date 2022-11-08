package main

import (
	"context"
	"fmt"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.SimpleServer.
type server struct{}

// SayHello implements helloworld.SimpleServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s", in.Name, in.Names, in.Embed.GetParam())}, nil
}

func (s *server) SayHello1(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s", in.Name, in.Names, in.Embed.GetParam())}, nil
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling. Info: %+v %+v", info, req)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)
	return resp, err
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/helloworld.SimpleService/SayHello1" {
		userName, err := CheckAuth(ctx)
		if err != nil {
			return nil, err
		}
		log.Printf("userName: %s", userName)
	}
	resp, err := handler(ctx, req)
	return resp, err
}

func main() {
	go func() {
		runUserService()
	}()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 可以传多个拦截器
	s := grpc.NewServer([]grpc.ServerOption{grpc.ChainUnaryInterceptor(UnaryServerInterceptor, AuthInterceptor)}...)
	pb.RegisterSimpleServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
