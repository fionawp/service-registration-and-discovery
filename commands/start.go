package commands

import (
	"fmt"
	"github.com/fionawp/service-registration-and-discovery/context"
	"github.com/fionawp/service-registration-and-discovery/server"
	"github.com/urfave/cli"
	"log"
)

// Starts web server (user interface)
var StartCommand = cli.Command{
	Name:   "start",
	Usage:  "Starts web server",
	Flags:  startFlags,
	Action: startAction,
}

var startFlags = []cli.Flag{
	cli.IntFlag{
		Name:   "http-port, p",
		Usage:  "HTTP server port",
		Value:  8089,
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_PORT",
	},
	cli.StringFlag{
		Name:   "http-host, i",
		Usage:  "HTTP server host",
		Value:  "",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_HOST",
	},
	cli.StringFlag{
		Name:   "http-mode, m",
		Usage:  "debug, release or test",
		Value:  "",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_MODE",
	},
}

func startAction(ctx *cli.Context) error {
	conf := context.NewConfig(ctx)

	if conf.HttpServerPort() < 1 {
		log.Fatal("Server port must be a positive integer")
	}

	fmt.Printf("Starting web server at %s:%d...\n", ctx.String("http-host"), ctx.Int("http-port"))

	server.Start(conf)

	fmt.Println("Done.")

	return nil
}
