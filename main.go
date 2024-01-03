package main

import (
	"github.com/pterm/pterm"
	urfave "github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func cliGetTemplateVars(c *urfave.Context) urfave.ExitCoder {
	templateName := c.Args().Get(0)
	if len(templateName) == 0 {
		return nil
	}
	vars, err := getTemplateVars(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	for key := range vars {
		pterm.NewRGB(66, 135, 245).Print(key)
		pterm.NewRGB(255, 255, 255).Print(":")
		pterm.NewRGB(3, 252, 90).Println(vars[key])
	}
	return nil
}

func cliGetTemplateNames(c *urfave.Context) urfave.ExitCoder {
	names, err := getTemplateNames()
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	for i := range names {
		name, _ := strings.CutSuffix(names[i], ".yml")
		pterm.NewRGB(66, 135, 245).Println(name)
	}
	return nil
}

func cliDeployServer(c *urfave.Context) urfave.ExitCoder {
	templateName := c.Args().Get(0)
	if len(templateName) == 0 {
		return nil
	}

	templateVars, err := getTemplateVars(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	for i, s := range c.Args().Slice() {
		if i == 0 {
			continue
		}

		arg := strings.Split(s, "=")
		if len(arg) != 2 {
			continue
		}

		templateVars[arg[0]] = arg[1]
	}

	err = deployContainer(templateName, c.String("name"), templateVars)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	return nil
}

func cliStartServer(c *urfave.Context) urfave.ExitCoder {

	return nil
}

func main() {
	app := &urfave.App{
		Name:                 "conductor",
		Version:              "1.0.0",
		Description:          "Easily create and manage game servers in a dockerized environment",
		Usage:                "conductor [template name]",
		EnableBashCompletion: true,
		Action: func(c *urfave.Context) error {
			return cliGetTemplateVars(c)
		},
		Commands: []*urfave.Command{
			{
				Name:        "templates",
				Description: "List the built in templates",
				Usage:       "conductor templates",
				Action: func(c *urfave.Context) error {
					return cliGetTemplateNames(c)
				},
			},
			{
				Name:        "deploy",
				Description: "Deploy a new server",
				Usage:       "conductor deploy [template name] [variables]",
				Flags: []urfave.Flag{
					&urfave.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Set the name of the server",
					},
				},
				Action: func(c *urfave.Context) error {
					return cliDeployServer(c)
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "remove", "rm"},
			},
			{
				Name:        "set",
				Description: "Set a servers global variables",
				Usage:       "conductor set [flags]",
			},
			{
				Name:    "start",
				Aliases: []string{"begin"},
				Usage:   "conductor start [server name]",
				Action: func(c *urfave.Context) error {
					return cliStartServer(c)
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"halt", "quit", "kill"},
				Usage:   "conductor stop [server name]",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
