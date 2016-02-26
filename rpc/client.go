package rpc

import (
	"errors"

	"github.com/sayden/gubsub/types"
	"github.com/valyala/gorpc"
)

//Client is an rpc client with functions to easily communicate to other servers
type Client struct {
	ServiceAddress string
	rpcClient      *gorpc.Client
	dClient        *gorpc.DispatcherClient
}

type rPCClientActions interface {
	JoinCluster(m string) error
	BroadcastMessage(m *types.Message) error
}

func (c *Client) getOrCreateClient() error {
	if c.rpcClient == nil {
		if c.ServiceAddress == "" {
			return errors.New("Service address was not set before creating client. Set it using 'rpc.Client{ServiceAddress:[address]}'")
		}

		d := gorpc.NewDispatcher()
		c.rpcClient = gorpc.NewTCPClient(c.ServiceAddress)
		c.rpcClient.Start()
		c.dClient = d.NewServiceClient("ClusterActions", c.rpcClient)
	}

	return nil
}

//JoinCluster takes an hostname:port string and send it to an already running
//gubsub server
func (c *Client) JoinCluster(m string) error {
	err := c.getOrCreateClient()
	if err != nil {
		return err
	}

	_, err = c.dClient.Call("AddMember", m)
	if err != nil {
		return err
	}

	return nil
}

//BroadcastMessage must be called when a new message arrives at the local server
//to spread (broadcast) the message to connected clients
func (c *Client) BroadcastMessage(m *types.Message) error {
	err := c.getOrCreateClient()
	if err != nil {
		return err
	}

	_, err = c.dClient.Call("NewMessage", m)
	if err != nil {
		return err
	}

	return nil
}
