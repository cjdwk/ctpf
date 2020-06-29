package service

import (
	"sync"

	"github.com/Terry-Mao/goim/pkg/bufio"
	"github.com/oofpgDLD/ctpf/imp-comet/api"
)

// Channel used by message pusher send msg to write goroutine.
type Channel struct {
	CliProto Ring
	signal   chan *api.Proto
	Writer   bufio.Writer
	Reader   bufio.Reader
	Next     *Channel
	Prev     *Channel

	Mid      int64
	Key      string
	IP       string
	mutex    sync.RWMutex
}

// NewChannel new a channel.
func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)
	c.signal = make(chan *api.Proto, svr)
	return c
}

// Push server push message.
func (c *Channel) Push(p *api.Proto) (err error) {
	select {
	case c.signal <- p:
	default:
	}
	return
}

// Ready check the channel ready or close?
func (c *Channel) Ready() *api.Proto {
	return <-c.signal
}

// Signal send signal to the channel, protocol ready.
func (c *Channel) Signal() {
	c.signal <- api.ProtoReady
}

// Close close the channel.
func (c *Channel) Close() {
	c.signal <- api.ProtoFinish
}
