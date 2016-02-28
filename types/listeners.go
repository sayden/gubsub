package types

import "github.com/satori/go.uuid"

//Listener struct must hold the information needed by any type of listener
type Listener struct {
	ID    uuid.UUID     `json:"id"`
	Ch    chan *Message `json:"-"`
	Quit  chan bool     `json:"-"`
	Topic string        `json:"topic"`
}

//HTTPListenerData is to store the information needed by a HTTP listener type
//So if we want to POST to "http://localhost:8000" we should store this in
//TargetURL and "POST" in Method parameter.
type HTTPListenerData struct {
	TargetURL string `json:"targetURL"`
	Method    string `json:"method"`
}
