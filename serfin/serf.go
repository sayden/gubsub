package serfin

import (
	"os"

	log "github.com/Sirupsen/logrus"
	serfclient "github.com/hashicorp/serf/client"
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
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

func Join(targetServer string) error {
	var err error
	if serf == nil {
		serf, err = serfclient.NewRPCClient("localhost:7373")
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}
	}

	server := make([]string, 1)
	server[0] = targetServer
	_, err = serf.Join(server, false)

	return err
}

func ListMembers() ([]serfclient.Member, error) {
	serf, err := serfclient.NewRPCClient("127.0.0.1:7373")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return serf.Members()
}

func GetIP() (string, error){

}

func getSerfClient() (serfclient.RPCClient, error) {
	var err error
	if serf == nil {
		serf, err = serfclient.NewRPCClient("localhost:7373")
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}
	}
}