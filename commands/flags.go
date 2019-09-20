package commands

import "github.com/urfave/cli"

// Global CLI flags
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:   "debug",
		Usage:  "run in debug mode",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_DEBUG",
	},
	cli.IntFlag{
		Name:   "http-port, p",
		Usage:  "HTTP server port",
		Value:  8089,
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_PORT",
	},
	cli.StringFlag{
		Name:   "http-host, i",
		Usage:  "HTTP server host",
		Value:  "127.0.0.1",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_HOST",
	},
	cli.StringFlag{
		Name:   "http-mode, m",
		Usage:  "debug, release or test",
		Value:  "",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_MODE",
	},
	cli.StringFlag{
		Name:   "service-name, s",
		Usage:  "register a server for a service",
		Value:  "firstService",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_SERVICE_NAME",
	},
	cli.StringFlag{
		Name:   "consul-host, ch",
		Usage:  "consul host",
		Value:  "http://192.168.33.11:8500",
		EnvVar: "SERVICE_REGISTER_AND_DISCOVERY_CONSUL_HOST",
	},
}
