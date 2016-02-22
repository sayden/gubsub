package servers

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"

	"golang.org/x/net/websocket"
)

//StartSocketServer will launch the web socket server on the specified endpoint
func StartSocketServer(port int) {
	socket := http.NewServeMux()
	socket.Handle("/default", websocket.Handler(func(ws *websocket.Conn) {
		socketHandler(ws, "default")
	}))
	go func(socket *http.ServeMux) { http.ListenAndServe(fmt.Sprintf(":%d", port), socket) }(socket)
	fmt.Printf("Listening Websocket on port %d \n", port)
}

func socketHandler(ws *websocket.Conn, endpoint string) {

	c := make(chan *[]byte)
	q := make(chan bool)
	l := types.Listener{
		ID:    uuid.NewV4(),
		Ch:    c,
		Quit:  q,
		Topic: endpoint,
	}

	dispatcher.AddListener(ws, l)
}
