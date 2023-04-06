package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	Schema        = "schema"
	helloEndpoint = "test.com"
)

var helloAddrs = []string{"127.0.0.1:50051", "127.0.0.1:50053"}

type helloResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *helloResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*helloResolver) Close() {}

type helloResolverBuilder struct{}

func (*helloResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &helloResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			helloEndpoint: helloAddrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*helloResolverBuilder) Scheme() string { return Schema }

func init() {
	resolver.Register(&helloResolverBuilder{})
}
