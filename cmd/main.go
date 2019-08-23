package main

import (
	"github.com/fionawp/service-registration-and-discovery/commands"
	"github.com/urfave/cli"
	"os"
)

var version = "development"

func main() {
	app := cli.NewApp()
	app.Name = "service-registration-and-discovery"
	app.Usage = "service-registration-and-discovery"
	app.Version = version
	app.Copyright = "(c) 2019 The service-registration-and-discovery contributors <fionawp@126.com>"
	app.EnableBashCompletion = true
	app.Flags = commands.GlobalFlags

	app.Commands = []cli.Command{
		commands.ConfigCommand,
		commands.StartCommand,
	}

	app.Run(os.Args)
}
