package serf

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

//GetIPs will give you your own Member object in the cluster, the rest of the
//cluster members or an erro
func GetIPs() (*serfclient.Member, []serfclient.Member, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	//Local ips
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
		return nil, nil, err
	}
	if len(members) > 1 {
		for _, m := range members {
			for k, ip := range ips {
				if m.Addr.String() == ip {
					return &m, append(members[:k], members[k+1:]...), nil
				}
			}
		}
		return nil, nil, errors.New("Could not find the local ip")
	} else if len(members) == 0 {
		return nil, nil, errors.New("No servers found")
	}

	return &members[0], nil, nil
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
