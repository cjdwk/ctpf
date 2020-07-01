package main

import (
	"flag"
	log "github.com/golang/glog"
	"github.com/oofpgDLD/ctpf/imp-logic/conf"
	"github.com/oofpgDLD/ctpf/imp-logic/server/grpc"
	"github.com/oofpgDLD/ctpf/imp-logic/service"
	"os"
	"os/signal"
	"syscall"
)

const (
	ver   = "2.0.0"
	appid = "goim.logic"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("goim-logic [version: %s env: %+v] start", ver, conf.Conf.Env)

	//get server name
	//serverName := discovery.ServerName2(conf.Conf.Env)

	// logic
	srv := service.New(conf.Conf)
	//httpSrv := http.New(conf.Conf.HTTPServer, srv)
	go grpc.New(conf.Conf, srv)

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-logic get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			srv.Close()
			//httpSrv.Close()
			//rpcSrv.GracefulStop()
			log.Infof("goim-logic [version: %s] exit", ver)
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}