package main

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"
)

const (
	userServicePort = ":50052"
)

type userServer struct{}

func (s *userServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	if in.Username == "zw" && in.Password == "password" {
		tokenString := CreateToken(in.Username)
		return &pb.LoginReply{Status: "200", Token: tokenString}, nil
	} else {
		return &pb.LoginReply{Status: "403", Token: ""}, nil
	}
}

func CreateToken(userName string) (tokenString string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "lora-app-server",
		"aud":      "lora-app-server",
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
		"sub":      "user",
		"username": userName,
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		panic(err)
	}
	return tokenString
}

// AuthToken 自定义认证，需要实现 credentials.PerRPCCredentials 接口
type AuthToken struct {
	Token string
}

func (auth *AuthToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": auth.Token,
	}, nil
}

func (auth *AuthToken) RequireTransportSecurity() bool {
	return false
}

type Claims struct {
	jwt.StandardClaims
	UserName string `json:"username"`
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "ErrNoMetadataInContext")
	}
	// md 的类型是 type MD map[string][]string
	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "ErrNoAuthorizationInMetadata")
	}
	// 因此，token 是一个字符串数组，我们只用了 token[0]
	return token[0], nil
}

func CheckAuth(ctx context.Context) (string, error) {
	tokenStr, err := getTokenFromContext(ctx)
	if err != nil {
		return "", err
	}
	var clientClaims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &clientClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != "HS256" {
			return nil, errors.New("ParseWithClaims error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("ErrInvalidToken")
	}
	return clientClaims.UserName, nil
}

func runUserService() {
	lis, err := net.Listen("tcp", userServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
