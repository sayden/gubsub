package servers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sayden/gubsub/dispatcher"
)

//StartHTTPServer will launch the http server
func StartHTTPServer(port int) {
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/topic/default", httpMessageHandler)
	httpServer.HandleFunc("/topic", httpTopicHandler)
	println("Listening HTTP on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), httpServer)
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

		var j map[string]interface{}
		err = json.Unmarshal(*msg, &j)
		if err != nil {
			println("Error trying to parse json", err)
		}

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
