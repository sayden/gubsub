package grpcservice

import (
	serf "github.com/hashicorp/serf/client"
	"github.com/sayden/gubsub/types"
	"net"
	"testing"
	"time"
)

type mockDispatcher struct {
	t *testing.T
}

func (m mockDispatcher) DispatchMessageLocal(msg *types.Message) {
	if string(*msg.Data) != "hello" {
		m.t.Error("Message doesn't match")
	}
}

func TestStartServer(t *testing.T) {
	d := mockDispatcher{t: t}
	server := NewGRPCServer(d, 5124)

	m := serf.Member{
		Addr: net.ParseIP("127.0.0.1"),
		Port:5124,
	}

	to := "default"
	data := []byte("hello")
	msg := types.Message{
		Data:&data,
		Topic:&to,
	}

	time.Sleep(2 * time.Second)
	if server.Dispatcher == nil {
		t.Log("Dispatcher is nil in test")
	}
	_, err := server.SendMessage(&msg, m)
	if err != nil {
		t.Error(err)
	}
}
