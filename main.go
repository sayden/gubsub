package main

import "github.com/sayden/gubsub/servers"

func main() {
	servers.StartSocketServer(12345, "default")
	servers.StartHTTPServer(8002)
}
