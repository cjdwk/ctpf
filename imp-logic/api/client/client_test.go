package client

import (
	"context"
	"os"
	"testing"

	_ "github.com/micro/go-plugins"
	pb "github.com/oofpgDLD/ctpf/imp-logic/api"
	"github.com/oofpgDLD/ctpf/library/discovery"
	"github.com/oofpgDLD/ctpf/library/trace"
	"github.com/stretchr/testify/assert"
)

var cli *Client

func TestMain(m *testing.M) {
	d := &discovery.Discovery{
		Address:  "172.16.103.31:2379",
		Name:     "root",
		Password: "admin",
	}
	tc := &trace.Trace{
		LocalAgentHostPort: "172.16.103.31:5775",
	}
	cli = New("test-logic-gclient", d, tc)
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	in := &pb.PingReq{}
	reply, err := cli.Ping(context.Background(), in)
	assert.Nil(t, err)
	t.Log(reply)
}

func TestConnect(t *testing.T) {
	//in :=
	//cli.Connect(context.Background(),)
}

//
//func Test_ClientAuth(t *testing.T) {
//	os.Args = os.Args[:1]
//
//	c := &grpc.Discovery{
//		Address:  "172.16.103.31:2379",
//		Name:     "root",
//		Password: "admin",
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			err, ok := r.(error)
//			if ok {
//				t.Error(err)
//			} else {
//				t.Error(r)
//			}
//		}
//	}()
//	es := New(c)
//	ret, err := es.Auth(context.TODO(), &proto.AuthRequest{AppId: "1001", Token: "session-login=MTU4OTc5MzA2NXxEdi1CQkFFQ180SUFBUkFCRUFBQV9nRUJfNElBQmdaemRISnBibWNNQ1FBSFpHVjJkSGx3WlFaemRISnBibWNNQ1FBSFFXNWtjbTlwWkFaemRISnBibWNNQmdBRWRYVnBaQVp6ZEhKcGJtY01JZ0FnTmtSR05USXpOak5HTmtORlEwSkRNRVEwTlRjNU5EUXlPRE15TlVFM016RUdjM1J5YVc1bkRBY0FCV0Z3Y0Vsa0JuTjBjbWx1Wnd3R0FBUXhNREF4Qm5OMGNtbHVad3dHQUFSMGFXMWxCV2x1ZERZMEJBZ0EtZ0xrVGhvWFJnWnpkSEpwYm1jTUNRQUhkWE5sY2w5cFpBWnpkSEpwYm1jTUF3QUJOQVp6ZEhKcGJtY01Cd0FGZEc5clpXNEdjM1J5YVc1bkRDb0FLR1EyTVdZek5EVmhPVFJpWWprME1HSTBORGMwTlRVd01qY3hPVEExWlRVd1pHTTFNR1psTWpVPXwcqpFI8g8aQmOJzn0QjMJips-2pD37C4rWTJ0DMpJSoQ=="})
//	if err != nil {
//		t.Error("call auth failed", err)
//		return
//	}
//
//	t.Log("call auth success", ret)
//	return
//}
//
//func Test_ClientUser(t *testing.T) {
//	os.Args = os.Args[:1]
//
//	c := &grpc.Discovery{
//		Address:  "172.16.103.31:2379",
//		Name:     "root",
//		Password: "admin",
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			err, ok := r.(error)
//			if ok {
//				t.Error(err)
//			} else {
//				t.Error(r)
//			}
//		}
//	}()
//	es := New(c)
//	ret, err := es.User(context.TODO(), &proto.UserRequest{EpCode: "FZM0001", UserId: "1"})
//	if err != nil {
//		t.Error("call user failed", err)
//		return
//	}
//
//	t.Log("call user success", ret)
//	return
//}
//
//func Test_ClientUsers(t *testing.T) {
//	os.Args = os.Args[:1]
//
//	c := &grpc.Discovery{
//		Address:  "172.16.103.31:2379",
//		Name:     "root",
//		Password: "admin",
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			err, ok := r.(error)
//			if ok {
//				t.Error(err)
//			} else {
//				t.Error(r)
//			}
//		}
//	}()
//	es := New(c)
//	ret, err := es.Users(context.TODO(), &proto.UsersRequest{EpCode: "FZM0001"})
//	if err != nil {
//		t.Error("call users failed", err)
//		return
//	}
//
//	t.Log("call users success", ret)
//	return
//}
