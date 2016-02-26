package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/cli"
	"github.com/sayden/gubsub/serf"
)

func main() {
	go startSerf()
	cli.StartCli()
}

func startSerf() {
	code := serf.StartSerf()
	log.WithFields(log.Fields{
		"code": code,
	}).Info("Serf result")
}
