package servers

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"

	"golang.org/x/net/websocket"
)

var socket *http.ServeMux

//StartSocketServer will launch the web socket server on the specified endpoint
func StartSocketServer(port int, endpoint string) {
	socket = http.NewServeMux()
	socket.Handle(fmt.Sprintf("/%s", endpoint), websocket.Handler(
		func(ws *websocket.Conn) {
			socketHandler(ws, endpoint)
		}))

	go startServer(socket, port)

	fmt.Printf("Listening Websocket on port %d. ", port)
}

func startServer(socket *http.ServeMux, port int) {
	http.ListenAndServe(fmt.Sprintf(":%d", port), socket)
}

//AddTopic register a new topic as an endpoint. Since this moment, you can post
//messages to this new endpoint
func AddTopic(endpoint string) error {
	println("Adding topic", endpoint)
	socket.Handle(fmt.Sprintf("/%s", endpoint), websocket.Handler(
		func(ws *websocket.Conn) {
			socketHandler(ws, endpoint)
		}))
	return nil
}

func socketHandler(ws *websocket.Conn, endpoint string) {
	fmt.Printf("Listening topic %s \n", endpoint)

	//Creates new listener
	c := make(chan *[]byte)
	q := make(chan bool)
	l := types.Listener{
		ID:    uuid.NewV4(),
		Ch:    c,
		Quit:  q,
		Topic: endpoint,
	}

	dispatcher.AddListener(l)

	for {
		m := <-l.Ch
		ws.Write(*m)
	}
}
