package servers

import (
	"net/http"

	"github.com/sayden/gubsub/listener"
	"golang.org/x/net/websocket"
)

//AddClient must be called from HTTP endpoints with the new clients connected
//through websocket
func AddClient(r *http.Request, w http.ResponseWriter, endpoint string) {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		listener.NewSocketListener(ws, &endpoint)
	})
	handler.ServeHTTP(w, r)
}
