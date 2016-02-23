package types

import "github.com/satori/go.uuid"

type Listener struct {
	ID    uuid.UUID     `json:"id"`
	Ch    chan *Message `json:"-"`
	Quit  chan bool     `json:"-"`
	Topic string        `json:"topic"`
}
