package cli

import (
	"os"
	//"time"

	log "github.com/sayden/gubsub/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/sayden/gubsub/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/sayden/gubsub/dispatcher"
	"github.com/sayden/gubsub/serfin"
	"github.com/sayden/gubsub/servers"

	"time"

	"os/signal"
	"syscall"

	"fmt"

	"github.com/sayden/gubsub/types"
	"github.com/sayden/gubsub/grpc"
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
		cli.IntFlag{
			Name:  "serf-port, sp",
			Value: 7946,
			Usage: "Sets the Serf Bind address port",
		},
		cli.IntFlag{
			Name:  "serf-rpc, srp",
			Value: 7373,
			Usage: "Sets the Serf RPC port",
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
			Name:    "serf",
			Aliases: []string{"d"},
			Usage:   "serf [commands]",
			Subcommands: []cli.Command{
				{
					Name:  "event",
					Usage: "Send a custom event through the Serf cluster",
					Action: func(c *cli.Context) {
						//TODO Serf event command
						log.Info("Not implemented yet")
					},
				},
				{
					Name:        "join",
					Usage:       "joinTell Serf agent to join cluster",
					Description: "Pass a --server [server:port] as the server you want to connect to",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "server, s",
							Usage: "Sets the server url and port",
						},
					},
					Action: func(c *cli.Context) {
						targetServer := c.String("server")
						if targetServer == "" {
							log.Fatal("You have to provide a --server flag")
						}

						err := serfin.Join(targetServer)
						if err != nil {
							log.Error("Could not join cluster", err)
						} else {
							log.Debug("Joined to cluster successfully")
						}
					},
				},
				{
					Name:  "members",
					Usage: "Lists the members of a Serf cluster",
					Action: func(c *cli.Context) {
						members, err := serfin.ListMembers()
						if err != nil {
							log.Fatal("Error trying to get member list:")
						}
						for _, v := range members {
							log.Info(fmt.Sprintf("Name: %s, Addr: %s, Port: %d, Status: %s", v.Name, v.Addr, v.Port, v.Status))
						}
					},
				},
				{
					Name:  "query",
					Usage: "Send a query to the Serf cluster",
					Action: func(c *cli.Context) {
						//TODO Serf query command
						log.Info("Not implemented yet")
					},
				},
			},
		},
		{
			Name:  "server",
			Usage: "Start the publish subscribing server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "join, j",
					Usage: "Sets the server url and port",
				},
			},
			Action: func(c *cli.Context) {
				port := c.GlobalInt("port")
				topic := c.GlobalString("topic")
				join := c.String("join")

				if port == 0 {
					port = 8300
				}

				if topic == "" {
					topic = "default"
				}

				//Directly join to a different server. This could be improved
				if join != "" {
					go func(s string){
						log.Debug("Trying to connect to %s server", s)
						time.Sleep(5 * time.Second)
						serfin.Join(s)
					}(join)
				}

				go serfin.StartSerf()
				servers.StartHTTPServer(port, topic)
				//gRPC server is started automatically using init() function
			},
		},
	}

	go signalCapture()

	app.Run(os.Args)
}

//signalCapture is used to capture syscalls like Ctrl+C, to wait for some
//seconds after serf is shut down and panic the app (force the exit)
func signalCapture() {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	for {
		<-signalCh
		log.Info("Shutting down Gubsub. Waiting 5 seconds")
		time.Sleep(5 * time.Second)
		log.Info("Waited for 5 seconds. Bye!")
		os.Exit(0)
	}
}
