package rpc

import (
	"errors"
	"sync"

	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/types"
)

type Member struct {
	Address string `json:"address"`
}

type ClusterActions interface {
	AddMember(m string) error
	RemoveMember(m string) error
	ListMembers() []string
	NewMessage(m types.Message) error
}

type Cluster struct {
	members map[string]string
}

var mutex = &sync.Mutex{}

func (c *Cluster) NewMessage(m types.Message) error {
	dispatcher.DispatchMessage(&m)
	return nil
}

func (c *Cluster) AddMember(m string) error {
	mutex.Lock()
	c.members[m] = m
	mutex.Unlock()
	return nil
}

func (c *Cluster) RemoveMember(m string) error {
	mutex.Lock()
	found := false
	delete(c.members, m)
	mutex.Unlock()

	if !found {
		return errors.New("Member " + m + "not found")
	}

	return nil
}

func (c *Cluster) ListMembers() map[string]string {
	return c.members
}
