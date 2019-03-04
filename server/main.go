package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.SimpleServer.
type server struct{}

// SayHello implements helloworld.SimpleServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSimpleServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}