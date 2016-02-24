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

	// TOPIC endpoint handles topic creating, messages
	// DELETE /topic/{name} Delete that specific topic
	// POST /topic/{name}/message Adds a message to the topic dispatcher
	topic := r.Group("/topic")

	// GET /topic Returns all registered topics
	topic.GET("/", GetAllTopics)

	// POST /topic Creates a new topic
	topic.POST("/:name/message", PostMessageOnTopic)

	topic.GET("/:name", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		endpoint := getEndpoint(s.Request.URL.Path)
		ClientConnected(endpoint, s)
	})

	// LISTENERS endpoint provides information about current listeners
	// GET /listener returns all registered listeners and their topics
	listener := r.Group("/listener")
	listener.GET("/", GetAllListeners)

	// CONFIG endpoint provides information about configuration
	// GET /config returns information about the app

	r.Run(fmt.Sprintf(":%d", port))
}

//takes an url of topic "/topic/default" and returns only "default"
func getEndpoint(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}
