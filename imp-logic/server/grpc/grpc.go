package grpc

import (
	"context"
	xerrors "errors"
	"fmt"
	"time"

	"github.com/micro/go-micro"
	registry2 "github.com/micro/go-micro/registry"
	xserver "github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/inconshreveable/log15"
	pb "github.com/oofpgDLD/ctpf/imp-logic/api"
	logic "github.com/oofpgDLD/ctpf/imp-logic/service"
	"github.com/oofpgDLD/ctpf/imp-logic/conf"
)

var (
	log = log15.New("server", "grpc")
)

// New comet grpc server.
func New(c *conf.Config, s *logic.Service) {
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
		ServiceName: pb.ServerName, //自定义服务名称
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
		log.Error("init tracer failed", "err", err, "server", pb.ServerName)
		return
	}
	defer closer.Close()

	service := micro.NewService(
		// Set service registry
		micro.Registry(registry),
		// Set service name
		micro.Name(pb.ServerName),
		// Set trace
		micro.WrapHandler(opentracing.NewHandlerWrapper(tracer)),
		// Set log wrapper
		micro.WrapHandler(logWrapper),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	err = pb.RegisterLogicHandler(service.Server(), &server{s})
	if err != nil {
		log.Error("register server failed", "err", err, "server", pb.ServerName)
	}
	// Run the server
	if err = service.Run(); err != nil {
		log.Error("run grpc server failed", "err", err, "server", pb.ServerName)
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
	srv *logic.Service
}

var _ pb.LogicHandler = &server{}

// Ping Service
func (s *server) Ping(ctx context.Context, req *pb.PingReq, reply *pb.PingReply) error {
	reply = &pb.PingReply{}
	return s.srv.Ping(ctx)
}

// Close Service
func (s *server) Close(ctx context.Context, req *pb.CloseReq, reply *pb.CloseReply) error {
	reply = &pb.CloseReply{}
	return nil
}

// Connect connect a conn.
func (s *server) Connect(ctx context.Context, req *pb.ConnectReq, reply *pb.ConnectReply) error {
	mid, key, room, accepts, hb, err := s.srv.Connect(ctx, req.Server, req.Cookie, req.Token)
	if err != nil {
		reply = &pb.ConnectReply{}
		return err
	}
	reply = &pb.ConnectReply{Mid: mid, Key: key, RoomID: room, Accepts: accepts, Heartbeat: hb}
	return nil
}

// Disconnect disconnect a conn.
func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq, reply *pb.DisconnectReply) error {
	has, err := s.srv.Disconnect(ctx, req.Mid, req.Key, req.Server)
	if err != nil {
		reply = &pb.DisconnectReply{}
		return err
	}
	reply = &pb.DisconnectReply{Has: has}
	return nil
}

// Heartbeat beartbeat a conn.
func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq, reply *pb.HeartbeatReply) error {
	if err := s.srv.Heartbeat(ctx, req.Mid, req.Key, req.Server); err != nil {
		reply = &pb.HeartbeatReply{}
		return err
	}
	reply = &pb.HeartbeatReply{}
	return nil
}

// RenewOnline renew server online.
func (s *server) RenewOnline(ctx context.Context, req *pb.OnlineReq, reply *pb.OnlineReply) error {
	//TODO
	//allRoomCount, err := s.srv.RenewOnline(ctx, req.Server, req.RoomCount)
	var err error
	var allRoomCount = make(map[string]int32)
	if err != nil {
		reply = &pb.OnlineReply{}
		return err
	}
	reply = &pb.OnlineReply{AllRoomCount: allRoomCount}
	return nil
}

// Receive receive a message.
func (s *server) Receive(ctx context.Context, req *pb.ReceiveReq, reply *pb.ReceiveReply) error {
	//TODO
	//if err := s.srv.Receive(ctx, req.Mid, req.Proto); err != nil {
	if err := s.srv.Receive(ctx, req.Mid, nil); err != nil {
		reply = &pb.ReceiveReply{}
		return err
	}
	reply = &pb.ReceiveReply{}
	return nil
}

// nodes return nodes.
func (s *server) Nodes(ctx context.Context, req *pb.NodesReq, reply *pb.NodesReply) error {
	//TODO
	//return s.srv.NodesWeighted(ctx, req.Platform, req.ClientIP), nil
	return nil
}
