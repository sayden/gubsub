package servers

import (
	"fmt"
	"net/http"
)

func StartHTTPServer(port int, messageHandler func(w http.ResponseWriter, r *http.Request)) {
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/message", messageHandler)
	println("Listening HTTP on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), httpServer)
}
