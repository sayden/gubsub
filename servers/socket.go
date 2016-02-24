package servers

import (
	"fmt"
	"net/http"

	"github.com/olahol/melody"
	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
)

func ClientConnected(endpoint string, s *melody.Session) {
	fmt.Printf("Listening topic %s \n", endpoint)

	//Creates new listener
	l := types.Listener{
		ID:    uuid.NewV4(),
		Ch:    make(chan *types.Message),
		Quit:  make(chan bool),
		Topic: endpoint,
	}

	dispatcher.AddListenerToTopic(l, endpoint)

	for {
		m := <-l.Ch
		s.Write(*m.Data)
	}
}

func AddClient(r *http.Request, w http.ResponseWriter) {

}
