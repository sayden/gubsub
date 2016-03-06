package dispatcher

import (
	"testing"
	"github.com/sayden/gubsub/types"
"github.com/satori/go.uuid"
	"time"
)

func TestAddRemoveListener(t *testing.T){
	l := &types.Listener{
		ID: uuid.NewV4(),
		Ch: make(chan *types.Message),
		Quit: make(chan bool),
		Topic: "default",
	}

	ls := GetAllListeners()
	curLength := len(ls)

	AddListenerToTopic(l)
	l.ID = uuid.NewV4()
	AddListenerToTopic(l)
	l.ID = uuid.NewV4()
	AddListenerToTopic(l)

	//Check current listeners list is 3
	ls = GetAllListeners()
	if len(ls) != curLength + 3 {
		t.Error("Length of listeners has not raised in 3")
	}
	curLength = len(ls)
	t.Logf("Current length %d\n", curLength)
	err := RemoveListener(l)
	if err != nil {
		t.Error("Error removing listener")
	}
	ls = GetAllListeners()
	t.Logf("'After' length %d\n", len(ls))
	if len(ls) != curLength - 1 {
		t.Error("Length of listeners has not dropped in one")
	}
}

func TestGetAllTopics(t *testing.T){
	l := &types.Listener{
		ID: uuid.NewV4(),
		Ch: make(chan *types.Message),
		Quit: make(chan bool),
		Topic: "default",
	}

	AddListenerToTopic(l)
	l.ID = uuid.NewV4()
	l.Topic = "another"
	AddListenerToTopic(l)

	ts := GetAllTopics()
	if len(ts) != 2 {
		t.Error("Topic must be 2 after adding to listeners of different topics")
	}

	if ts[0] != "default" {
		t.Error("First topic should be 'default'")
	}
}

func TestChannelRoutes(t *testing.T){
	l := &types.Listener{
		ID: uuid.NewV4(),
		Ch: make(chan *types.Message),
		Quit: make(chan bool),
		Topic: "default",
	}

	AddListenerToTopic(l)
	go func(t *testing.T, l *types.Listener){
		for {
			select{
			case m  := <-l.Ch:
				if string(*m.Data) != "hello" {
					t.Error("Incorrect message arrived")
				}
			case <-l.Quit:
				return
			}
		}
	}(t, l)

	d := []byte("hello")
	to := "default"
	m := types.Message{
		Data:&d,
		Topic:&to,
		Timestamp:time.Now(),
	}

	disp.DispatchMessageLocal(&m)

	DispatchMessage(&m)

	l.Quit <- true
}