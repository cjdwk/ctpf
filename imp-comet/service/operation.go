package service

import (
	"context"
	"time"

	model "github.com/oofpgDLD/ctpf/imp-comet/api"
	logic "github.com/oofpgDLD/ctpf/imp-logic/api"

	_ "github.com/Terry-Mao/goim/pkg/strings"
	log "github.com/golang/glog"
)

// Connect connected a connection.
func (s *Service) Connect(c context.Context, p *model.Proto, cookie string) (mid int64, key string, gids []int64, heartbeat time.Duration, err error) {
	reply, err := s.rpcClient.Connect(c, &logic.ConnectReq{
		Server: s.serverID,
		Cookie: cookie,
		Token:  p.Body,
	})
	if err != nil {
		return
	}
	//TODO gids
	return reply.Mid, reply.Key, nil, time.Duration(reply.Heartbeat), nil
}

// Disconnect disconnected a connection.
func (s *Service) Disconnect(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Disconnect(context.Background(), &logic.DisconnectReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// Heartbeat heartbeat a connection session.
func (s *Service) Heartbeat(ctx context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Heartbeat(ctx, &logic.HeartbeatReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// RenewOnline renew room online.
func (s *Service) RenewOnline(ctx context.Context, serverID string, rommCount map[string]int32) (allRoom map[string]int32, err error) {
	reply, err := s.rpcClient.RenewOnline(ctx, &logic.OnlineReq{
		Server:    s.serverID,
		RoomCount: rommCount,
	})//grpc.UseCompressor(gzip.Name))
	if err != nil {
		return
	}
	return reply.AllRoomCount, nil
}

// Receive receive a message.
func (s *Service) Receive(ctx context.Context, mid int64, p *model.Proto) (err error) {
	_, err = s.rpcClient.Receive(ctx, &logic.ReceiveReq{Mid: mid, Proto: p})
	return
}

// Operate operate.
func (s *Service) Operate(ctx context.Context, p *model.Proto, ch *Channel, b *Bucket) error {
	switch p.Op {
	default:
		// TODO ack ok&failed
		if err := s.Receive(ctx, ch.Mid, p); err != nil {
			log.Errorf("s.Report(%d) op:%d error(%v)", ch.Mid, p.Op, err)
		}
		p.Body = nil
	}
	return nil
}