package dispatcher

import (
	"sync"

	serfclient "github.com/hashicorp/serf/client"
	log "github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/serf"
	"github.com/sayden/gubsub/types"
	"time"
)

var mutex = &sync.Mutex{}

type dispatcher struct {
	topics            map[string][]types.Listener
	listeners         []types.Listener
	msgDispatcher     chan *types.Message
	clusterDispatcher chan *types.Message
	dispatch          chan *[]byte
}

type Servers []serfclient.Member

var d *dispatcher
var servers Servers

func init() {
	d = &dispatcher{
		topics:        make(map[string][]types.Listener),
		listeners:     make([]types.Listener, 1),
		msgDispatcher: make(chan *types.Message, 20),
		clusterDispatcher: make(chan *types.Message, 100),
		dispatch:      make(chan *[]byte),
	}

	d.AddTopic("default")

	go d.topicDispatcherLoop()
	go refreshMemberListLoop()
}

func refreshMemberListLoop() {
	for {
		time.Sleep(5 * time.Second)
		_, ms, err := serf.GetIPs()
		if err != nil {
			log.Error(err)
		}

		mutex.Lock()
		servers = ms
		mutex.Unlock()
	}
}

//AddTopic adds a new topic to be available to listeners. It will expose a new
//endpoint to be connected too
func (d *dispatcher) AddTopic(name string) error {
	mutex.Lock()
	d.topics[name] = make([]types.Listener, 0)
	mutex.Unlock()

	return nil
}

//DispatchMessage takes a message and inserts it into the generic messages channel
//that will distribute it to the registered servers
func DispatchMessage(m *types.Message) {
	d.clusterDispatcher <- m
	go DispatchMessageLocal(m)
}

//DispatchMessageLocal takes a message and inserts it into the generic messages
// channel that will distribute it to the registered listeners
func DispatchMessageLocal(m *types.Message) {
	d.msgDispatcher <- m
}

func (d *dispatcher) topicDispatcherLoop() {
	for {
		m := <-d.msgDispatcher
		ls := d.topics[*m.Topic]
		for _, l := range ls {
			l.Ch <- m
		}
	}
}

//AddListenerToTopic will add a listener to one of the topic's arrays so it can
//be notified since that moment of new messages
func AddListenerToTopic(l *types.Listener, topic *string) {
	log.WithFields(log.Fields{
		"ID":    l.ID,
		"topic": *topic,
	}).Info("New listener")

	mutex.Lock()
	d.topics[*topic] = append(d.topics[*topic], *l)
	mutex.Unlock()
}

//GetAllTopics will return an array of strings with the registered topic names
func GetAllTopics() []string {
	var ts []string
	for k := range d.topics {
		ts = append(ts, k)
	}
	return ts
}

//GetAllListeners will return an array with all listeners for each topic
func GetAllListeners() []types.Listener {
	var ls []types.Listener
	for k := range d.topics {
		for _, l := range d.topics[k] {
			ls = append(ls, l)
		}
	}
	return ls
}

//RemoveListener takes a types.Listener and a topic and removes the listener from
//the the specified topic
func RemoveListener(l *types.Listener) error {
	mutex.Lock()
	ls := d.topics[l.Topic]
	for i, _l := range ls {
		if *l == _l {
			ls = append(ls[:i], ls[i+1:]...)
		}
	}

	d.topics[l.Topic] = ls
	mutex.Unlock()

	return nil
}
