package main

import (
	"context"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	address            = "localhost:50051"
	userServiceAddress = "localhost:50052"
	defaultName        = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSimpleServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("call response message: %s", r.Message)

	// 服务端抛出了ErrNoAuthorizationInMetadata的错误
	r, err = c.SayHello1(ctx, &pb.HelloRequest{Name: name})
	log.Printf("%v", err)

	conn1, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn1.Close()
	userCli := pb.NewUserServiceClient(conn1)
	loginReply, err := userCli.Login(ctx, &pb.LoginRequest{
		Username: "zw",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("Error when calling Login: %s", err)
	}
	log.Println("Login Reply:", loginReply)
	requestToken := AuthToken{Token: loginReply.Token}
	//requestToken := new(AuthToken)
	//requestToken.Token = loginReply.Token
	log.Println("??:", requestToken)

	// 携带token创建新的客户端连接
	conn2, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&requestToken))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn2.Close()
	c = pb.NewSimpleServiceClient(conn2)
	r, err = c.SayHello1(ctx, &pb.HelloRequest{Name: name})
	log.Printf("%v %v", r, err)
}
