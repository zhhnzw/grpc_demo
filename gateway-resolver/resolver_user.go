package main

import (
	"google.golang.org/grpc/resolver"
)

const (
	userEndpoint = "test.user.com"
)

var userAddrs = []string{"127.0.0.1:50052", "127.0.0.1:50054"}

type userResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *userResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*userResolver) Close() {}

type userResolverBuilder struct{}

func (*userResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &userResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			userEndpoint: userAddrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*userResolverBuilder) Scheme() string { return Schema }

func init() {
	resolver.Register(&userResolverBuilder{})
}
