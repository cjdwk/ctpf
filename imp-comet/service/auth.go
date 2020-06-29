package service

import (
	"context"
	"time"

	"github.com/oofpgDLD/ctpf/imp-comet/api"

	"github.com/Terry-Mao/goim/pkg/bufio"
	log "github.com/golang/glog"
)

func (s *Service) Auth(ctx context.Context, rr *bufio.Reader, wr *bufio.Writer, p *api.Proto) (mid int64,key string, gids []int64, hb time.Duration, err error) {
	for {
		if err = p.ReadTCP(rr); err != nil {
			return
		}
		if p.Op == api.OpAuth {
			break
		} else {
			log.Errorf("tcp request operation(%d) not auth", p.Op)
		}
	}
	if mid, key, gids, hb, err = s.Connect(ctx, p, ""); err != nil {
		log.Errorf("authTCP.Connect(key:%v).err(%v)", key, err)
		return
	}
	p.Op = api.OpAuthReply
	p.Body = nil
	if err = p.WriteTCP(wr); err != nil {
		log.Errorf("authTCP.WriteTCP(key:%v).err(%v)", key, err)
		return
	}
	err = wr.Flush()
	return
}