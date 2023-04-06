package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

// server is used to implement helloworld.SimpleServer.
type server struct{}

// SayHello implements helloworld.SimpleServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s", in.Name, in.Names, in.Embed.GetParam())}, nil
}

func (s *server) SayHello1(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	claims := ctx.Value("claims").(Claims)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s, names: %v, embed.param: %s claims:%+v", in.Name, in.Names, in.Embed.GetParam(), claims)}, nil
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("before handling. Info: %+v %+v", info, req)
	resp, err := handler(ctx, req)
	log.Printf("after handling. resp: %+v", resp)
	return resp, err
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/helloworld.SimpleService/SayHello1" {
		claims, err := CheckAuth(ctx)
		if err != nil {
			return nil, err
		}
		log.Printf("auth claims: %+v", claims)
		ctx = context.WithValue(ctx, "claims", claims)
	}
	resp, err := handler(ctx, req)
	return resp, err
}

func main() {
	helloServicePort := flag.Int("helloServicePort", 50051, "listen helloService port")
	userServicePort := flag.Int("userServicePort", 50052, "listen userService port")
	flag.Parse()
	go func() {
		runUserService(*userServicePort)
	}()
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*helloServicePort))
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
