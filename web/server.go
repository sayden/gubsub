package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sayden/gubsub/queue"
)

//Start will run the web server
func Start(q queue.Stacker) {
	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, q)
	})
	http.ListenAndServe(":12345", nil)
}

//Parses a request to return a queue.Message object
func getMessageFromRequest(r *http.Request) (*queue.Message, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}

	msg := queue.Message{Data: &body}
	return &msg, nil
}

func handler(w http.ResponseWriter, r *http.Request, q queue.Stacker) {
	if r.Method == "POST" {
		msg, err := getMessageFromRequest(r)
		if err != nil {
			fmt.Println("Error trying to parse message from body", err)
			w.WriteHeader(422) // unprocessable entity
		} else {
			err = q.AddToStack(msg, "default")

			if err != nil {
				w.WriteHeader(422) // unprocessable entity
				fmt.Println("Error:", err)
			}
			w.WriteHeader(200)
		}
	} else {
		w.WriteHeader(405)
	}
}
