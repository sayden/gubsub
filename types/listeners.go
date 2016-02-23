package types

import "github.com/satori/go.uuid"

type Listener struct {
	ID    uuid.UUID
	Ch    *chan *[]byte
	Quit  *chan bool
	Topic string
}
