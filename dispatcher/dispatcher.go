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

	d.AddTopic("default")

	go d.topicDispatcherLoop()
}

func (d *Dispatcher) AddTopic(name string) error {
	mutex.Lock()
	d.topics[name] = make([]types.Listener, 0)
	mutex.Unlock()

	return nil
}

func DispatchMessage(m *types.Message) {
	d.msgDispatcher <- m
}

func (d *Dispatcher) topicDispatcherLoop() {
	for {
		m := <-d.msgDispatcher
		ls := d.topics[*m.Topic]
		for _, l := range ls {
			l.Ch <- m
		}
	}
}

func AddListenerToTopic(l types.Listener, topic string) {
	fmt.Printf("New listener for topic %s\n", topic)

	mutex.Lock()
	d.topics[topic] = append(d.topics[topic], l)
	mutex.Unlock()

	ls := d.topics[topic]
	for _, l := range ls {
		fmt.Printf("%s listener in topic %s\n", l.ID, l.Topic)
	}
}

func GetAllTopics() []string {
	var ts []string
	for k := range d.topics {
		ts = append(ts, k)
	}
	return ts
}

func GetAllListeners() []types.Listener {
	var ls []types.Listener
	for k := range d.topics {
		for _, l := range d.topics[k] {
			ls = append(ls, l)
		}
	}
	return ls
}
