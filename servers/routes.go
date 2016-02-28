package servers

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/listener"
	"github.com/sayden/gubsub/types"
)

//GetAllTopics is the "GET /topic" REST handler that will return all available
//topics in the session
func GetAllTopics(c *gin.Context) {
	ts := dispatcher.GetAllTopics()
	c.JSON(http.StatusOK, gin.H{"result": ts})
}

//PostMessageOnTopic is the "POST /topic/[a topic]/message" REST handler to insert
//messages in the queue
func PostMessageOnTopic(c *gin.Context) {
	name := c.Param("name")

	var msg types.Message
	err := c.BindWith(&msg, &types.MessageBinding{})
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"result": err.Error()})
	} else {
		msg.Topic = &name
		msg.Timestamp = time.Now()
		dispatcher.DispatchMessage(&msg)
		c.JSON(http.StatusOK, gin.H{"result": "Message received"})
	}
}

//GetAllListeners is the "GET /listener" REST handler to return all connected
//listeners
func GetAllListeners(c *gin.Context) {
	ls := dispatcher.GetAllListeners()
	c.JSON(http.StatusOK, gin.H{"total": len(ls), "result": ls})
}

func CreateNewHTTPListener(c *gin.Context) {
	endpoint := c.Param("endpoint")

	var msg types.HTTPListenerData
	err := c.BindJSON(&msg)
	if err != nil {
		log.Error("Couldn't parse json", err)
		c.JSON(http.StatusNotAcceptable, gin.H{"result": err.Error()})
	} else {
		var id uuid.UUID
		listener.NewHTTPListener(msg, &endpoint, &id)
		if &id != nil {
			c.JSON(http.StatusOK, gin.H{"result": id})
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{"result": err.Error()})
		}
	}
}
