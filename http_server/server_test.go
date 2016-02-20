package server

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/sayden/gubsub/queue"
)

type MockStacker struct{}

func (m *MockStacker) AddToStack(msg *queue.Message, topic string) error {
	return nil
}

func TestServerCanReceivePostRequests(t *testing.T) {
	m := MockStacker{}
	go Start(&m)
	time.Sleep(1 * time.Second)

	//Create a POST request, should return a 200
	r, err := http.PostForm("http://localhost:12345/receive",
		url.Values{"key": {"Value"}, "id": {"123"}})
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 200 {
		t.Fatal("Status code of server response was not 200")
	}

	//Create a GET request, should return a "not allowed"
	r, err = http.Get("http://localhost:12345/receive")
	if err != nil {
		t.Fatal(err)
	}
	if r.StatusCode != 405 {
		t.Fatal("Status code of server response was not 405")
	}
}
