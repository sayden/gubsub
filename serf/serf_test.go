package serf
import (
	"testing"
	serfclient "github.com/hashicorp/serf/client"
	"net"
)


func TestGetLocalNetworksIPs(t *testing.T){
	ips := getLocalNetworksIPs()
	if len(ips) == 0 {
		t.Error("No networks have been found")
	}

	found := false
	for _, v := range ips {
		if v == "127.0.0.1"{
			found = true
		}
	}

	if !found {
		t.Error("127.0.0.1 was not found")
	}
}

func TestGetMatchingMembers(t *testing.T){
	//Add some members first with IP's 127.0.0.1 and 192.168.1.10
	ip := net.ParseIP("127.0.0.1")
	m1 := serfclient.Member{ Addr:ip }
	ms := make([]serfclient.Member, 0)
	ms = append(ms, m1)

	ip = net.ParseIP("192.168.1.10")
	m2 := serfclient.Member{ Addr:ip }
	ms = append(ms, m2)

	ip = net.ParseIP("192.168.1.100")
	m3 := serfclient.Member{ Addr:ip }
	ms = append(ms, m3)

	if len(ms) != 3 {
		t.Errorf("Members must be 3 but have %d", len(ms))
	}

	ips := make([]string,1)
	ips[0] = "192.168.1.10"
	member, rest, err := getMatchingMembers(ips, ms)
	if err != nil {
		t.Error(err)
	}

	if len(rest) != 2 {
		t.Error("Rest of interfaces must be two")
		for _, v := range rest {
			t.Logf("I: %s\n", v.Addr.String())
		}
	}

	if member.Addr.String() != "192.168.1.10" {
		t.Error("Function didn't extract correct ip")
	}
}