package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

var dispatch = make(chan *[]byte)
var listeners = make([]chan *[]byte, 0)
var mutex = &sync.Mutex{}

func main() {
	//Creates socket server
	socket := http.NewServeMux()
	socket.Handle("/websocket", websocket.Handler(func(ws *websocket.Conn) {
		c := make(chan *[]byte)
		mutex.Lock()
		listeners = append(listeners, c)
		mutex.Unlock()
		socketHandler(ws, c)
	}))
	go func(socket *http.ServeMux) { http.ListenAndServe(":12345", socket) }(socket)
	println("Listening Websocket on port 12345")

	go dispatcher()

	//Creates http server
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/message", messageHandler)
	println("Listening HTTP on port 8002")
	http.ListenAndServe(":8002", httpServer)
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

func dispatcher() {
	for {
		m := <-dispatch
		for _, l := range listeners {
			l <- m
		}
	}
}

func socketHandler(ws *websocket.Conn, c chan *[]byte) {
	println("New client")
	for {
		m := <-c
		ws.Write(*m)
	}
}
