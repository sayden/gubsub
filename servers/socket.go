package servers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
	"golang.org/x/net/websocket"
)

func clientConnected(ws *websocket.Conn, endpoint string) {
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
		_, err := ws.Write(*m.Data)
		if err != nil {
			log.Error("Error trying to write on socket: ", err)
			dispatcher.RemoveListener(l)
		}
	}
}

//AddClient must be called from HTTP endpoints with the new clients connected
//through websocket
func AddClient(r *http.Request, w http.ResponseWriter, endpoint string) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		clientConnected(ws, endpoint)
	})
	handler.ServeHTTP(w, r)
}
