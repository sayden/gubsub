package servers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//StartHTTPServer is the starting point of the application. It is called from
//the main package and will configure all needed endpoints as well as the
//"default" endpoint to receive messages
func StartHTTPServer(port int, endpoint string) {
	r := gin.Default()

	topic := r.Group("/topic")
	topic.GET("/", GetAllTopics)
	topic.POST("/:name/message", PostMessageOnTopic)
	topic.GET("/:name", func(c *gin.Context) {
		endpoint := c.Param("name")
		AddClient(c.Request, c.Writer, endpoint)
	})

	listener := r.Group("/listener")
	listener.GET("/", GetAllListeners)
	listener.POST("/http/topic/:endpoint", CreateNewHTTPListener)
	listener.POST("/file/topic/:endpoint", CreateNewFileListener)

	r.Run(fmt.Sprintf(":%d", port))
}
