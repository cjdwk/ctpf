package grpc

import (
	"context"
	xerrors "errors"
	"fmt"
	"github.com/micro/go-micro"
	registry2 "github.com/micro/go-micro/registry"
	xserver "github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	"time"

	"github.com/inconshreveable/log15"

	pb "github.com/oofpgDLD/ctpf/imp-comet/api"
	comet "github.com/oofpgDLD/ctpf/imp-comet/service"
	"github.com/oofpgDLD/ctpf/imp-comet/conf"
	"github.com/oofpgDLD/ctpf/imp-comet/errors"
)

const ServerName = "imp-logic"

var (
	log = log15.New("server", "grpc")
)

// New comet grpc server.
func New(c *conf.Config, s *comet.Service) {
	if c.Discovery == nil {
		err := xerrors.New("discovery config not find")
		log.Error("init grpc api failed", "err", err)
		panic(err)
	}
	if c.Trace == nil {
		err := xerrors.New("trace config not find")
		log.Error("init grpc api failed", "err", err)
		panic(err)
	}

	//etcdv3插件
	registry := etcdv3.NewRegistry(
		registry2.Addrs(c.Discovery.Address), //etch 服务器地址
		etcdv3.Auth(c.Discovery.Name, c.Discovery.Password),
	)

	cfg := config.Configuration{
		ServiceName: ServerName, //自定义服务名称
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  c.Trace.LocalAgentHostPort, //jaeger agent
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Error("init tracer failed", "err", err, "server", ServerName)
		return
	}
	defer closer.Close()

	service := micro.NewService(
		// Set service registry
		micro.Registry(registry),
		// Set service name
		micro.Name(ServerName),
		// Set trace
		micro.WrapHandler(opentracing.NewHandlerWrapper(tracer)),
		// Set log wrapper
		micro.WrapHandler(logWrapper),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	err = pb.RegisterCometHandler(service.Server(), &server{s})
	if err != nil {
		log.Error("register server failed", "err", err, "server", ServerName)
	}
	// Run the server
	if err = service.Run(); err != nil {
		log.Error("run grpc server failed", "err", err, "server", ServerName)
	}
}

// logWrapper is a handler wrapper
func logWrapper(fn xserver.HandlerFunc) xserver.HandlerFunc {
	return func(ctx context.Context, req xserver.Request, rsp interface{}) error {
		log.Info(fmt.Sprintf("[wrapper] server request: %v", req.Endpoint()))
		err := fn(ctx, req, rsp)
		return err
	}
}

type server struct {
	srv *comet.Service
}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.Empty, reply *pb.Empty) error{
	reply = &pb.Empty{}
	return nil
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.Empty, reply *pb.Empty) error{
	// TODO: some graceful close
	reply = &pb.Empty{}
	return nil
}

// PushMsg push a message to specified sub keys.
func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq, reply *pb.PushMsgReply) error{
	if len(req.Keys) == 0 || req.Proto == nil {
		return errors.ErrPushMsgArg
	}
	for _, key := range req.Keys {
		if channel := s.srv.Bucket(key).Channel(key); channel != nil {
			/*if !channel.NeedPush(req.ProtoOp) {
				continue
			}*/
			if err := channel.Push(req.Proto); err != nil {
				return err
			}
		}
	}
	reply = &pb.PushMsgReply{}
	return nil
}

// Broadcast broadcast msg to all user.
func (s *server) Broadcast(ctx context.Context, req *pb.BroadcastReq, reply *pb.BroadcastReply) error{
	if req.Proto == nil || req.GroupId == "" {
		return errors.ErrBroadCastRoomArg
	}
	for _, bucket := range s.srv.Buckets() {
		bucket.BroadcastGroup(req)
	}
	reply = &pb.BroadcastReply{}
	return nil
}