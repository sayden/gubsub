package serfin

import (
	"os"

	log "github.com/Sirupsen/logrus"
	serfclient "github.com/hashicorp/serf/client"
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	"fmt"
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
	serf, err := getSerfClient()
	if err != nil {
		return err
	}

	server := make([]string, 1)
	server[0] = targetServer
	_, err = serf.Join(server, false)

	return err
}

func ListMembers() ([]serfclient.Member, error) {
	serf, err := getSerfClient()
	if err != nil {
		return nil, err
	}

	return serf.Members()
}

func GetIP() (string, error){
	serf, err := getSerfClient()
	if err != nil {
		return "", err
	}

	stats, two, three, four := serf.ListKeys()

	if err != nil {
		return "", err
	}

	log.Info(stats)
	log.Info(two)
	log.Info(three)
	fmt.Printf("%+v\n", three)
	log.Info(four)
	return "", nil
}

func getSerfClient() (*serfclient.RPCClient, error) {
	var err error
	if serf == nil {
		serf, err = serfclient.NewRPCClient("localhost:7373")
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}
	}

	return serf, nil
}