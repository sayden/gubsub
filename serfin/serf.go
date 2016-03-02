package serfin

import (
	"os"

	"errors"
	"net"

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

func GetIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var ips []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Error("Error getting address", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ips = append(ips, ip.String())
		}
	}

	serf, err := getSerfClient()
	members, err := serf.Members()
	if err != nil {
		return "", err
	}
	for _, m := range members {
		for _, ip := range ips {
			if m.Addr.String() == ip {
				return ip, nil
			}
		}
	}

	return "", errors.New("Could not find the local ip")
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
