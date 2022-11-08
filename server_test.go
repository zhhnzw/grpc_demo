package main

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	pb "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	address            = "localhost:50051"
	userServiceAddress = "localhost:50052"
	defaultName        = "world"
)

func TestGRPC(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	assert.Nil(t, err)
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
	assert.Nil(t, err)
	t.Logf("call response message: %s", r.Message)

	// 服务端抛出了ErrNoAuthorizationInMetadata的错误
	r, err = c.SayHello1(ctx, &pb.HelloRequest{Name: name})
	t.Logf("%v", err)

	conn1, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	assert.Nil(t, err)
	defer conn1.Close()
	userCli := pb.NewUserServiceClient(conn1)
	loginReply, err := userCli.Login(ctx, &pb.LoginRequest{
		Username: "zw",
		Password: "password",
	})
	assert.Nil(t, err)
	t.Log("Login Reply:", loginReply)
	requestToken := AuthToken{Token: loginReply.Token}

	// 携带token创建新的客户端连接
	conn2, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&requestToken))
	assert.Nil(t, err)
	defer conn2.Close()
	c = pb.NewSimpleServiceClient(conn2)
	r, err = c.SayHello1(ctx, &pb.HelloRequest{Name: name})
	assert.Nil(t, err)
	t.Logf("%v %v", r, err)
}

func TestGateWay(t *testing.T) {
	urlpfx := "http://localhost:9090"
	// 测试直接访问
	helloRequest := pb.HelloRequest{Name: "zhangsan"}
	helloRequestByte, _ := json.Marshal(helloRequest)
	request, _ := http.NewRequest("GET", urlpfx+"/helloworld", strings.NewReader(string(helloRequestByte)))
	request.Header.Set("Content-Type", "application/json")
	helloResponse, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)
	helloReplyBytes, err := ioutil.ReadAll(helloResponse.Body)
	assert.Nil(t, err)
	t.Logf(string(helloReplyBytes))

	// 登录
	loginRequest := pb.LoginRequest{}
	loginrequestByte, _ := json.Marshal(loginRequest)
	request, _ = http.NewRequest("POST", urlpfx+"/login", strings.NewReader(string(loginrequestByte)))
	request.Header.Set("Content-Type", "application/json")
	loginResponse, _ := http.DefaultClient.Do(request)
	loginReplyBytes, _ := ioutil.ReadAll(loginResponse.Body)
	defer loginResponse.Body.Close()
	var loginReply pb.LoginReply
	json.Unmarshal(loginReplyBytes, &loginReply)
	assert.Equal(t, "403", loginReply.Status)

	loginRequest.Username = "zw"
	loginRequest.Password = "password"
	loginrequestByte, _ = json.Marshal(loginRequest)
	request, _ = http.NewRequest("POST", urlpfx+"/login", strings.NewReader(string(loginrequestByte)))
	loginResponse, _ = http.DefaultClient.Do(request)
	loginReplyBytes, _ = ioutil.ReadAll(loginResponse.Body)
	defer loginResponse.Body.Close()
	json.Unmarshal(loginReplyBytes, &loginReply)
	assert.Equal(t, "200", loginReply.Status)
	t.Log(loginResponse.Status, loginReply)

	helloRequest = pb.HelloRequest{Name: "zhangsan"}
	helloRequestByte, _ = json.Marshal(helloRequest)
	request, _ = http.NewRequest("GET", urlpfx+"/helloworld", strings.NewReader(string(helloRequestByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", loginReply.Token)
	helloResponse, err = http.DefaultClient.Do(request)
	assert.Nil(t, err)
	helloReplyBytes, err = ioutil.ReadAll(helloResponse.Body)
	assert.Nil(t, err)
	t.Logf(string(helloReplyBytes))
}
