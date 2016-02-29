package cli

import (
	"os"
	//"time"

	log "github.com/sayden/gubsub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/servers"
	//"github.com/sayden/gubsub/types"
	"time"
"github.com/sayden/gubsub/types"
	"fmt"
)

func StartCli() {
	app := cli.NewApp()
	app.Name = "Gubsub"
	app.Usage = "A very simple yet powerful real-time message broker"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "topic, t",
			Value: "default",
			Usage: "Sets the name of the default topic",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 8002,
			Usage: "Sets the broker port",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "topics",
			Aliases: []string{"t"},
			Usage:   "Gets the current registered topics",
			Action: func(c *cli.Context) {
				ts := dispatcher.GetAllTopics()
				log.Info(ts)
			},
		},
		{
			Name:    "listeners",
			Aliases: []string{"l"},
			Usage:   "Returns all connected listeners",
			Action: func(c *cli.Context) {
				ls := dispatcher.GetAllListeners()
				log.Info(ls)
			},
		},
		{
			Name:    "dispatch",
			Aliases: []string{"d"},
			Usage:   "dispatch [topic] [message]",
			Action: func(c *cli.Context) {
				data := []byte(c.Args()[1])
				topic := c.Args()[0]
				dispatcher.DispatchMessage(&types.Message{
					Data:      &data,
					Topic:     &topic,
					Timestamp: time.Now(),
				})
			},
		},
		{
			Name:    "server",
			Usage:   "Start the publish subscribing server",
			Action: func(c *cli.Context) {
				port := c.Int("port")
				topic := c.String("topic")
				if port == 0 {
					port = 8300
				}

				if topic == "" {
					topic = "default"
				}

				servers.StartHTTPServer(port, topic)
			},
		},
	}

	app.Run(os.Args)
}
