package servers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
)

func GetAllTopics(c *gin.Context) {
	ts := dispatcher.GetAllTopics()
	c.JSON(http.StatusOK, gin.H{"result": ts})
}

func PostMessageOnTopic(c *gin.Context) {
	name := c.Param("name")

	var msg types.Message
	err := c.BindWith(&msg, &types.MessageBinding{})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
	} else {
		msg.Topic = &name
		msg.Timestamp = time.Now()
		dispatcher.DispatchMessage(&msg)
		c.JSON(http.StatusOK, gin.H{"result": "Message received"})
	}
}

func GetAllListeners(c *gin.Context) {
	ls := dispatcher.GetAllListeners()
	c.JSON(http.StatusOK, gin.H{"result": ls})
}
