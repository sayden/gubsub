package dispatcher

import (
	"fmt"
	"sync"

	"github.com/sayden/gubsub/types"
	"golang.org/x/net/websocket"
)

var dispatch = make(chan *[]byte)
var msgDispatcher = make(chan *types.Message)
var listeners = make([]types.Listener, 0)
var topics = make(map[string][]types.Listener)
var mutex = &sync.Mutex{}

func init() {
	go dispatcherLoop()
	//Add default topic
	AddTopic("default")
}

func AddTopic(name string) error {
	mutex.Lock()
	topics[name] = make([]types.Listener, 0)
	mutex.Unlock()
	return nil
}

//Dispatch takes a message and distributes it among registered listeners
func Dispatch(m *[]byte) {
	dispatch <- m
}

func topicDispatcherLoop() {
	for {
		m := <-msgDispatcher
		ls := topics[m.Topic]
		for _, l := range ls {
			l.Ch <- m.Data
		}
	}
}

func dispatcherLoop() {
	for {
		m := <-dispatch
		for _, l := range listeners {
			l.Ch <- m
		}
	}
}

//AddListener will make a Listener to receive all incoming messages
func AddListener(ws *websocket.Conn, l types.Listener) {
	println("New client")

	mutex.Lock()
	listeners = append(listeners, l)
	mutex.Unlock()
	for {
		m := <-l.Ch
		ws.Write(*m)
	}
}

func AddListenerToTopic(ws *websocket.Conn, l types.Listener, topic string) {
	fmt.Printf("New listener for topic %s", topic)

	mutex.Lock()
	topics[topic] = append(topics[topic], l)
	mutex.Unlock()
	for {
		m := <-l.Ch
		ws.Write(*m)
	}
}

func RemoveListener(l types.Listener) error {
	for i, registeredL := range listeners {
		if l == registeredL {
			mutex.Lock()
			listeners = append(listeners[:1], listeners[i+1:]...)
			mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("Listener %d not found in pool", l.ID)
}
