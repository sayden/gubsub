package servers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

func StartHTTPServer(port int, endpoint string) {
	r := gin.Default()
	m := melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// TOPIC endpoint handles topic creating, messages
	// POST /topic Creates a new topic
	// DELETE /topic/{name} Delete that specific topic
	// POST /topic/{name}/message Adds a message to the topic dispatcher
	topic := r.Group("/topic")

	// GET /topic Returns all registered topics
	topic.GET("/", GetAllTopics)

	topic.GET("/default", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		ClientConnected(endpoint, s)
	})

	topic.POST("/:name/message", PostMessageOnTopic)

	listener := r.Group("/listener")
	listener.GET("/", GetAllListeners)

	// LISTENERS endpoint provides information about current listeners
	// GET /listener returns all registered listeners and their topics

	// CONFIG endpoint provides information about configuration
	// GET /config returns information about the app

	r.Run(fmt.Sprintf(":%d", port))
}
