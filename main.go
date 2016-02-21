package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/sayden/gubsub/servers"
	"github.com/sayden/gubsub/types"

	"golang.org/x/net/websocket"
)

var dispatch = make(chan *[]byte)

var mutex = &sync.Mutex{}
var Listeners = make([]types.Listener, 0)

func main() {
	servers.StartSocketServer(12345, sockerHandler)

	go dispatcherLoop()

	servers.StartHTTPServer(8002, messageHandler)
}

func sockerHandler(ws *websocket.Conn) {
	c := make(chan *[]byte)
	q := make(chan bool)
	l := types.Listener{c, q, "default"}

	mutex.Lock()
	Listeners = append(Listeners, l)
	mutex.Unlock()

	startListener(ws, l)
}

//Parses a request to return a message.Message object
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

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		msg, err := getMessageFromRequest(r)
		if err != nil {
			fmt.Println("Error trying to parse message from body", err)
			w.WriteHeader(422) // unprocessable entity
		} else {
			dispatch <- msg

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

func dispatcherLoop() {
	for {
		m := <-dispatch
		for _, l := range Listeners {
			l.Ch <- m
		}
	}
}

func startListener(ws *websocket.Conn, l types.Listener) {
	println("New client")
	for {
		m := <-l.Ch
		ws.Write(*m)
	}
}
