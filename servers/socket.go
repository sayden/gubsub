package servers

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func StartSocketServer(port int, socketHandler func(ws *websocket.Conn)) {
	socket := http.NewServeMux()
	socket.Handle("/websocket", websocket.Handler(socketHandler))
	go func(socket *http.ServeMux) { http.ListenAndServe(fmt.Sprintf(":%d", port), socket) }(socket)
	println("Listening Websocket on port 12345")
}
