package transport

import (
	"github.com/alexandria-oss/core/config"
	"github.com/alexandria-oss/core/transport/proxy"
	"google.golang.org/grpc"
)

type Transport struct {
	RPCProxy   *grpc.Server
	HTTPProxy  *proxy.HTTP
	EventProxy *proxy.Event
	Config     *config.Kernel
}

func NewTransport(rpcProxy *grpc.Server, httpProxy *proxy.HTTP, eventProxy *proxy.Event, cfg *config.Kernel) *Transport {
	return &Transport{rpcProxy, httpProxy, eventProxy, cfg}
}
