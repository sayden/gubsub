package servers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
)

func StartHTTPServer(port int) {
	r := gin.Default()

	topic := r.Group("/topic")
	topic.POST("/:name", func(c *gin.Context) {
		name := c.Param("name")
		AddTopic(name)
	})

	topic.POST("/:name/message", func(c *gin.Context) {
		name := c.Param("name")

		var msg types.Message
		err := c.BindWith(&msg, &types.MessageBinding{})
		if err != nil {
			msg.Topic = &name
			dispatcher.DispatchMessage(&msg)
			c.JSON(http.StatusOK, gin.H{"result": "Message received"})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		}
	})

	r.Run(fmt.Sprintf(":%d", port))
}

//StartHTTPServer will launch the http server
func StartHTTPServer2(port int) {
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/topic/default", httpMessageHandler)
	httpServer.HandleFunc("/topic", httpTopicHandler)
	println("Listening HTTP on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), httpServer)
}

func handleRouting(w http.ResponseWriter, r *http.Request) {
	//Get first endpoint

	// TOPIC endpoint handles topic creating, messages
	// POST /topic Creates a new topic
	// GET /topic Returns all registered topics
	// DELETE /topic/{name} Delete that specific topic
	// POST /topic/{name}/message Adds a message to the topic dispatcher

	// LISTENERS endpoint provides information about current listeners
	// GET /listener returns all registered listeners and their topics

	// CONFIG endpoint provides information about configuration
	// GET /config returns information about the app
}

//Parses a request to return a []byte object with the body content
func getMessageFromRequest(r *http.Request) (*[]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}

	return &body, nil
}

func httpTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//Create a new topic with the body contents

		msg, err := getMessageFromRequest(r)
		if err != nil {
			fmt.Println("Error trying to create topic", err)
		}

		var j map[string]string
		err = json.Unmarshal(*msg, &j)
		if err != nil {
			println("Error trying to parse json", err)
		}

		AddTopic(j["topicName"])
		w.WriteHeader(200)
	} else {
		w.WriteHeader(405)
	}
}

func httpMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		msg, err := getMessageFromRequest(r)
		if err != nil {
			fmt.Println("Error trying to parse message from body", err)
			w.WriteHeader(422) // unprocessable entity
		} else {
			if err != nil {
				w.WriteHeader(422) // unprocessable entity
				fmt.Println("Error:", err)
			}
			dispatcher.Dispatch(msg)
			w.WriteHeader(200)
		}
	} else {
		w.WriteHeader(405)
	}
}
