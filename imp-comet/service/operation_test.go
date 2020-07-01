package service

import (
	"context"
	"encoding/json"
	"github.com/oofpgDLD/ctpf/imp-comet/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	opHeartbeat      = int32(2)
	opHeartbeatReply = int32(3)
	opAuth           = int32(7)
	opAuthReply      = int32(8)
)

// AuthToken auth token.
type AuthToken struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

func TestConnect(t *testing.T) {
	mid := int64(1)
	seq := int32(0)
	authToken := &AuthToken{
		mid,
		"",
		"test://1",
		"ios",
		[]int32{1000, 1001, 1002},
	}
	body, err := json.Marshal(authToken)
	if err != nil {
		t.Error(err)
		return
	}
	p := &api.Proto{
		Ver: 1,
		Op: opAuth,
		Seq: seq,
		Body:body,
	}
	mid, key, gids, heartbeat, err := s.Connect(context.Background(), p, "")
	assert.Nil(t, err)
	t.Log(mid, key, gids, heartbeat)
}