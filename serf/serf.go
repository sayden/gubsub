package serf

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func StartSerf() int {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := make([]string, 2)
	args[1] = "agent"

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("serf"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
