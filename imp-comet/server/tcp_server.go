package server

import(
	log "github.com/golang/glog"
	"github.com/oofpgDLD/ctpf/imp-comet/service"
	"net"
)

const (
	maxInt = 1<<31 - 1
)

func InitTCP(svc *service.Service) error{
	var (
		r int
	)
	ls, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		return err
	}
	log.Info("serve tcp start listen [list:%v]")
	go func(lis net.Listener) {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Errorf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
				continue
			}
			go svc.ServeConn(conn, r)
			if r++; r == maxInt {
				r = 0
			}
		}
	}(ls)
	return nil
}