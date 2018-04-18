package main

import (
	"github.com/paralin/protods/generate"
	_ "github.com/paralin/protods/generate/itypes"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var generateOutputPath = "."

func init() {
	var subCommands []cli.Command
	generate.ForEachGenerator(func(name string, gen generate.Generator) bool {
		subCommands = append(subCommands, cli.Command{
			Name:  name,
			Usage: gen.GetUsage(),
			Action: func(c *cli.Context) error {
				protoPathArg := c.Args().Get(0)
				if protoPathArg == "" {
					return errors.New("specify proto file to use")
				}

				return generate.Generate(gen, protoPathArg, generateOutputPath)
			},
		})
		return true
	})

	rootCommands = append(rootCommands, cli.Command{
		Name:        "generate",
		Usage:       "generate protods go structures",
		Subcommands: subCommands,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "go_out, o",
				Usage:       "output go code to `PATH`",
				Destination: &generateOutputPath,
				Value:       generateOutputPath,
			},
		},
	})
}
