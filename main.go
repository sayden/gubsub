package main

import "github.com/sayden/gubsub/servers"

func main() {
	servers.StartHTTPServer(8002, "default")
}
