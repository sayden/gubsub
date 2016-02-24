package servers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func StartHTTPServer(port int, endpoint string) {
	r := gin.Default()
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	topic := r.Group("/topic")
	topic.GET("/", GetAllTopics)
	topic.POST("/:name/message", PostMessageOnTopic)
	topic.GET("/:name", func(c *gin.Context) {
		AddClient(c.Request, c.Writer)
		// m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		endpoint := getEndpoint(s.Request.URL.Path)
		ClientConnected(endpoint, s)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		println("Client disconnected in topic", getEndpoint(s.Request.URL.Path))
	})

	listener := r.Group("/listener")
	listener.GET("/", GetAllListeners)

	r.Run(fmt.Sprintf(":%d", port))
}

//takes an url of topic "/topic/default" and returns only "default"
func getEndpoint(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}
