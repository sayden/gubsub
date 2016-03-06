package dispatcher

import (
	"sync"

	"time"

	log "github.com/Sirupsen/logrus"
	serfclient "github.com/hashicorp/serf/client"
	"github.com/sayden/gubsub/config"
	"github.com/sayden/gubsub/grpc"
	"github.com/sayden/gubsub/serf"
	"github.com/sayden/gubsub/types"
	"github.com/spf13/viper"
)

var mutex = &sync.Mutex{}

type dispatcher struct {
	topics            map[string][]types.Listener
	listeners         []types.Listener
	msgDispatcher     chan *types.Message
	clusterDispatcher chan *types.Message
	dispatch          chan *[]byte
	servers           []serfclient.Member
}

var Dispatcher *dispatcher

func init() {

	Dispatcher = &dispatcher{
		topics:            make(map[string][]types.Listener),
		listeners:         make([]types.Listener, 1),
		msgDispatcher:     make(chan *types.Message, viper.GetInt(config.MESSAGE_SIZE)),
		clusterDispatcher: make(chan *types.Message, viper.GetInt(config.MESSAGE_CLUSTER_SIZE)),
		//dispatch:          make(chan *[]byte),
		servers: []serfclient.Member{},
	}

	Dispatcher.AddTopic("default")

	go Dispatcher.topicDispatcherLoop()
	go Dispatcher.refreshMemberListLoop()
	go Dispatcher.clusterDispatcherLoop()
}

func (d *dispatcher) refreshMemberListLoop() {
	for {
		time.Sleep(time.Duration(viper.GetInt(config.MEMBER_LIST_REFRESH_SECONDS)) * time.Second)
		_, ms, err := serf.GetIPs()
		if err != nil {
			log.Error(err)
		}

		mutex.Lock()
		d.servers = ms
		mutex.Unlock()
	}
}

func (d *dispatcher) clusterDispatcherLoop() {
	for {
		m := <-d.clusterDispatcher
		grpcservice.RPC.SendMessageInCluster(m, d.servers)
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
	Dispatcher.DispatchMessage(m)
}

//DispatchMessage takes a message and inserts it into the generic messages channel
//that will distribute it to the registered servers
func (d *dispatcher) DispatchMessage(m *types.Message) {
	d.clusterDispatcher <- m
	go d.DispatchMessageLocal(m)
}

//DispatchMessageLocal takes a message and inserts it into the generic messages
// channel that will distribute it to the registered listeners
func (d *dispatcher) DispatchMessageLocal(m *types.Message) {
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
func AddListenerToTopic(l *types.Listener) {
	log.WithFields(log.Fields{
		"ID":    l.ID,
		"topic": l.Topic,
	}).Info("New listener")

	mutex.Lock()
	Dispatcher.topics[l.Topic] = append(Dispatcher.topics[l.Topic], *l)
	mutex.Unlock()
}

//GetAllTopics will return an array of strings with the registered topic names
func GetAllTopics() []string {
	var ts []string
	for k := range Dispatcher.topics {
		ts = append(ts, k)
	}
	return ts
}

//GetAllListeners will return an array with all listeners for each topic
func GetAllListeners() []types.Listener {
	var ls []types.Listener
	for k := range Dispatcher.topics {
		for _, l := range Dispatcher.topics[k] {
			ls = append(ls, l)
		}
	}
	return ls
}

//RemoveListener takes a types.Listener and a topic and removes the listener from
//the the specified topic
func RemoveListener(l *types.Listener) error {
	mutex.Lock()
	ls := Dispatcher.topics[l.Topic]
	for i, _l := range ls {
		if *l == _l {
			ls = append(ls[:i], ls[i+1:]...)
		}
	}

	Dispatcher.topics[l.Topic] = ls
	mutex.Unlock()

	return nil
}
