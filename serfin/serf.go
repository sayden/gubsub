package serfin

import (
	"os"

	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	serfclient "github.com/hashicorp/serf/client"
	log "github.com/Sirupsen/logrus"
)

//// Serfer is the common serf client for gubsub.
var serf *serfclient.RPCClient

func StartSerf() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	serf := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}

	serf.Run(nil)
}

func JoinSerfin() error {
	var err error
	serf, err = serfclient.NewRPCClient("localhost:7946")
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	server := make([]string, 1)
	server[0] = "localhost:7946"
	_, err = serf.Join(server, false)

	return err
}
