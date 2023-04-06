package main

import (
	"flag"
	"github.com/golang/glog"
	gw "github.com/zhhnzw/grpc_demo/helloworld"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"regexp"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	endpoint            = flag.String("endpoint", "localhost:50051", "endpoint of YourService")
	userServiceEndpoint = flag.String("userServiceEndpoint", "localhost:50052", "endpoint of YourService")
)

func allowedOrigin(origin string) bool {
	if viper.GetString("cors") == "*" {
		return true
	}
	if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
		return true
	}
	return false
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allowedOrigin(r.Header.Get("Origin")) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	err := gw.RegisterSimpleServiceHandlerFromEndpoint(
		ctx,
		mux,
		"schema:///test.com",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithResolvers(&helloResolverBuilder{}),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		})
	if err != nil {
		return err
	}
	err = gw.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		mux,
		"schema:///test.user.com",
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithResolvers(&userResolverBuilder{}),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		})
	if err != nil {
		return err
	}
	srv := http.Server{
		Addr:    ":9090",
		Handler: cors(mux),
	}
	return srv.ListenAndServe()
}

func main() {
	flag.Parse()
	defer glog.Flush()
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
