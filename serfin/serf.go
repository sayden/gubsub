package serfin

import (
	"os"

	"github.com/hashicorp/serf/command/agent"
	"github.com/mitchellh/cli"
)

func StartSerf() {
	ui := &cli.BasicUi{Writer: os.Stdout}
	serf := &agent.Command{
		Ui:         ui,
		ShutdownCh: make(chan struct{}),
	}

	serf.Run(nil)
}
