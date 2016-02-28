package listener

import (
	log "github.com/sayden/gubsub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/Godeps/_workspace/src/golang.org/x/net/websocket"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
)

//NewSocketListener must be called everytime that a new client is connected
//through websocket. It will configure and make it ready to receive messages
func NewSocketListener(ws *websocket.Conn, endpoint *string) {
	//Creates new listener
	l := types.NewListener(endpoint)

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
