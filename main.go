package main

import "github.com/sayden/gubsub/servers"

func main() {
	servers.StartSocketServer(12345)
	servers.StartHTTPServer(8002)
}
