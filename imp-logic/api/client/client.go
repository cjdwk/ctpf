package client

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/micro/go-micro/client"
	xgrpc "github.com/micro/go-micro/client/grpc"
	registry2 "github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/oofpgDLD/ctpf/library/discovery"
	"github.com/oofpgDLD/ctpf/imp-comet/conf"
	proto "github.com/oofpgDLD/ctpf/imp-logic/api"
)

var (
	log = log15.New("client", "grpc")
)

func New(serverName string, dCfg *discovery.Discovery, tCfg *conf.Trace) *Client {
	if dCfg == nil {
		err := errors.New("discovery config not find")
		log.Error("init grpc client failed", "err", err)
		panic(err)
	}

	if tCfg == nil {
		err := errors.New("trace config not find")
		log.Error("init grpc client failed", "err", err)
		panic(err)
	}

	if serverName == "" {
		err := errors.New("undefined server name")
		log.Error("init grpc client failed", "err", err)
		panic(err)
	}
	//etcd插件
	registry := etcdv3.NewRegistry(
		registry2.Addrs(dCfg.Address),
		etcdv3.Auth(dCfg.Name, dCfg.Password),
	)

	cfg := config.Configuration{
		ServiceName: serverName, //自定义服务名称
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  tCfg.LocalAgentHostPort, //jaeger agent
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Error("init tracer failed", "err", err, "server", serverName)
		return nil
	}

	c := xgrpc.NewClient(
		client.Registry(registry),
		client.Wrap(opentracing.NewClientWrapper(tracer)),
	)
	// Create new greeter client
	return &Client{
		LogicService: proto.NewLogicService(serverName, c),
		closer:            closer,
	}
}

type Client struct {
	proto.LogicService
	closer io.Closer
}

func (c *Client) Close(ctx context.Context, in *proto.CloseReq, opts ...client.CallOption) (*proto.CloseReply, error){
	c.closer.Close()
	return c.LogicService.Close(ctx, in, opts...)
}