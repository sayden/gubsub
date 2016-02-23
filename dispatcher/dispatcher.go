package dispatcher

import (
	"fmt"
	"sync"

	"github.com/sayden/gubsub/types"
)

var mutex = &sync.Mutex{}

type Dispatcher struct {
	topics        map[string][]types.Listener
	listeners     []types.Listener
	msgDispatcher chan *types.Message
	dispatch      chan *[]byte
}

var d *Dispatcher

func init() {
	d = &Dispatcher{
		topics:        make(map[string][]types.Listener),
		listeners:     make([]types.Listener, 1),
		msgDispatcher: make(chan *types.Message, 20),
		dispatch:      make(chan *[]byte),
	}

	//Add default topic
	d.AddTopic("default")

	// go dispatcherLoop()
	go d.topicDispatcherLoop()
}

func (d *Dispatcher) AddTopic(name string) error {
	mutex.Lock()
	d.topics[name] = make([]types.Listener, 0)
	mutex.Unlock()
	println(len(d.topics))
	return nil
}

func DispatchMessage(m *types.Message) {
	d.msgDispatcher <- m
}

//Dispatch takes a message and distributes it among registered listeners
func Dispatch(m *[]byte) {
	d.dispatch <- m
}

func (d *Dispatcher) topicDispatcherLoop() {
	for {
		m := <-d.msgDispatcher
		ls := d.topics[*m.Topic]
		for _, l := range ls {
			*l.Ch <- m.Data
		}
	}
}

func (d *Dispatcher) dispatcherLoop() {
	for {
		m := <-d.dispatch
		for _, l := range d.listeners {
			*l.Ch <- m
		}
	}
}

//AddListener will make a Listener to receive all incoming messages
func AddListener(l types.Listener) {
	println("New client")

	mutex.Lock()
	d.listeners = append(d.listeners, l)
	mutex.Unlock()
	// for {
	// 	m := <-l.Ch
	// 	ws.Write(*m)
	// }
}

func AddListenerToTopic(l types.Listener, topic string) {
	fmt.Printf("New listener for topic %s\n", topic)

	mutex.Lock()
	d.topics[topic] = append(d.topics[topic], l)
	mutex.Unlock()
	// for {
	// 	m := <-l.Ch
	// 	ws.Write(*m)
	// }

	ls := d.topics[topic]
	for _, l := range ls {
		fmt.Printf("%s listener in topic %s\n", l.ID, l.Topic)
	}
}

func RemoveListener(l types.Listener) error {
	for i, registeredL := range d.listeners {
		if l == registeredL {
			mutex.Lock()
			d.listeners = append(d.listeners[:1], d.listeners[i+1:]...)
			mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("Listener %d not found in pool", l.ID)
}
