package listener

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
	"golang.org/x/net/websocket"
)

//NewSocketListener must be called everytime that a new client is connected
//through websocket. It will configure and make it ready to receive messages
func NewSocketListener(ws *websocket.Conn, endpoint *string) {
	//Creates new listener
	l := types.NewListener(endpoint)

	dispatcher.AddListenerToTopic(l)

	for {
		m := <-l.Ch
		_, err := ws.Write(*m.Data)
		if err != nil {
			log.Error("Error trying to write on socket: ", err)
			dispatcher.RemoveListener(l)
		}
	}
}
