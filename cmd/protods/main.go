package main

import (
	"github.com/urfave/cli"
	"os"
)

var rootCommands []cli.Command
var rootFlags []cli.Flag

func main() {
	app := cli.NewApp()
	app.Name = "protods"
	app.Usage = "use any backing datastructure or datastore for complex protobuf structures"
	app.Authors = []cli.Author{
		cli.Author{Name: "Christian Stewart", Email: "christian@paral.in"},
	}
	app.Commands = rootCommands
	app.Flags = rootFlags

	if err := app.Run(os.Args); err != nil {
		_, _ = os.Stderr.WriteString("Error: ")
		_, _ = os.Stderr.WriteString(err.Error())
		_, _ = os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
