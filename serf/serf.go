package serf

import (
	"os"

	"errors"
	"net"

	"fmt"
	log "github.com/Sirupsen/logrus"
	serfclient "github.com/hashicorp/serf/client"
	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
	"github.com/sayden/gubsub/config"
	"github.com/spf13/viper"
)

var serfClient *serfclient.RPCClient
var serfServer *agent.Command

func StartSerf() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	serfServer := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}

	serfServer.Run(nil)
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
	//Local ips
	ips := getLocalNetworksIPs()
	if len(ips) == 0 {
		return nil, nil, errors.New("No networks found")
	}

	serf, err := getSerfClient()
	members, err := serf.Members()
	if err != nil {
		return nil, nil, err
	}

	return getMatchingMembers(ips, members)
}

func getMatchingMembers(ips []string, members []serfclient.Member) (*serfclient.Member, []serfclient.Member, error) {
	if len(members) > 1 { //two or more
		for k, m := range members {
			for _, ip := range ips {
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

func getLocalNetworksIPs() (ips []string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	//Local ips
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

	return
}

func getSerfClient() (*serfclient.RPCClient, error) {
	var err error
	if serfClient == nil {
		serfClient, err = serfclient.NewRPCClient(
			fmt.Sprintf("localhost:%s", viper.GetInt(config.SERF_RPC)))
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(1)
		}
	}

	return serfClient, nil
}
