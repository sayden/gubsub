package types

type Listener struct {
	Ch    chan *[]byte
	Quit  chan bool
	Topic string
}
