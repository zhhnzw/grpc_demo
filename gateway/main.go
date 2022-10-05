package main

import (
	"flag"
	"github.com/golang/glog"
	gw "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	echoEndpoint = flag.String("echo_endpoint", "localhost:50051", "endpoint of YourService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterSimpleServiceHandlerFromEndpoint(ctx, mux, *echoEndpoint, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":9090", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
