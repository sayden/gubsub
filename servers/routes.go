package servers

import (
	"net/http"
	"time"

	log "github.com/sayden/gubsub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/satori/go.uuid"
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
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
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

//CretaeNewHTTPListener will add to the listeners queue a new listener that will
//execute an HTTP request to a specified point
func CreateNewHTTPListener(c *gin.Context) {
	endpoint := c.Param("endpoint")

	var msg types.HTTPListener
	err := c.BindJSON(&msg)
	if err != nil {
		log.Error("Couldn't parse json", err)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	} else {
		var id uuid.UUID
		listener.NewHTTPListener(msg, &endpoint, &id)
		if &id != nil {
			c.JSON(http.StatusOK, gin.H{"result": id.String()})
		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		}
	}
}

//CreateNewFileListener will add to the listener queue a new listener that will
// write to a file
func CreateNewFileListener(c *gin.Context) {
	endpoint := c.Param("endpoint")

	var msg types.FileListener
	err := c.BindJSON(&msg)
	if err != nil {
		log.Error("Couldn't parse json", err)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
	} else {
		id, err := listener.NewFileListener(msg, &endpoint)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Error creating file")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error creating file",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": id.String()})
		}
	}

}
