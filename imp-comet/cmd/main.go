package main

import (
	"flag"
	log "github.com/golang/glog"
	proto "github.com/oofpgDLD/ctpf/imp-comet/api"
	"github.com/oofpgDLD/ctpf/imp-comet/conf"
	"github.com/oofpgDLD/ctpf/imp-comet/server"
	"github.com/oofpgDLD/ctpf/imp-comet/server/grpc"
	"github.com/oofpgDLD/ctpf/imp-comet/service"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const(
	ver   = "1.0.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Infof("environment=%v",conf.Conf.Debug)
	log.Infof("ctpf-comet [version=%v] start", ver)

	//get server name
	//serverName := discovery.ServerName(conf.Conf.Env)

	// new comet server
	srv := service.New(proto.ServerName, conf.Conf)

	//tcp serve
	if err := server.InitTCP(srv); err != nil {
		panic(err)
	}

	go grpc.New(conf.Conf, srv)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-comet get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//rpcSrv.GracefulStop()
			srv.Close()
			log.Infof("ctpf-comet [version: %s] exit", ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}