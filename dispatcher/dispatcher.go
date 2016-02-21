package dispatcher

import "github.com/sayden/gubsub/types"

var dispatch = make(chan *[]byte)
var Listeners = make([]types.Listener, 0)

func StartLoop() {
	for {
		m := <-dispatch
		for _, l := range Listeners {
			l.Ch <- m
		}
	}
}

func Dispatch(m *[]byte) {
	dispatch <- m
}
