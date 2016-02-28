package main

import (
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/olahol/melody"
	"net/http"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}
