package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.SimpleServer.
type server struct{}

// SayHello implements helloworld.SimpleServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s", in.Name, in.Names, in.Embed.GetParam())}, nil
}

func (s *server) SayHello1(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s", in.Name, in.Names, in.Embed.GetParam())}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 可以传多个拦截器
	s := grpc.NewServer([]grpc.ServerOption{grpc.UnaryInterceptor(UnaryServerInterceptor)}...)
	pb.RegisterSimpleServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling. Info: %+v", info)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)
	return resp, err
}
