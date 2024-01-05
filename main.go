package main

import (
	"fmt"
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

func cliGetTemplateNames() urfave.ExitCoder {
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

	serverName := c.Args().Get(1)
	if len(templateName) == 0 {
		return nil
	}

	templateVars, err := overrideTemplateVars(templateName, c.Args().Slice())
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	err = deployContainer(templateName, serverName, templateVars)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	if c.Bool("start") {
		return cliStartServer(c)
	}
	return nil
}

func cliDeleteServer(c *urfave.Context) urfave.ExitCoder {
	serverName := c.Args().Get(0)
	if len(serverName) == 0 {
		return nil
	}

	if !serverExists(serverName) {
		return urfave.Exit("No server with that name", 1)
	}

	var answer string
	if !c.Bool("yes") {
		fmt.Print("Are you sure you want to delete " + serverName + "? (y/N) ")
		_, err := fmt.Scanln(&answer)
		if err != nil {
			return nil
		}
	}

	if c.Bool("yes") || strings.Contains(answer, "y") || strings.Contains(answer, "Y") {
		err := deleteContainer(serverName)
		if err != nil {
			return urfave.Exit(err.Error(), 1)
		}
	}

	return nil
}

func cliStartServer(c *urfave.Context) urfave.ExitCoder {
	/*templateName := c.Args().Get(0)
	if len(templateName) == 0 {
		return nil
	}

	template, err := parseTemplate(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	templateVars, err := parseServerTemplateVars(c.String("name"))
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	var startCmd = templateVars + "\n" + template.Actions.Start

	err := runCommandInContainer()
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}*/

	return nil
}

/*func cliSetServerTemplateVars(c *urfave.Context) urfave.ExitCoder {
	serverName := c.Args().Get(0)
	if len(serverName) == 0 {
		return nil
	}

	templateVars, err := overrideTemplateVars(templateName, c.Args().Slice())
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}
}*/

func main() {
	app := &urfave.App{
		Name:                 "conductor",
		Version:              "1.0.0",
		Description:          "Easily create and manage game servers in a dockerized environment",
		Usage:                "conductor [template name]",
		Args:                 true,
		EnableBashCompletion: true,
		Action: func(c *urfave.Context) error {
			return cliGetTemplateVars(c)
		},
		Commands: []*urfave.Command{
			{
				Name:        "templates",
				Description: "List the built in templates",
				Usage:       "conductor templates",
				Args:        false,
				Action: func(c *urfave.Context) error {
					return cliGetTemplateNames()
				},
			},
			{
				Name:        "deploy",
				Description: "Deploy a new server",
				Usage:       "conductor deploy [flags] [template name] [server name] [variable overrides]",
				Args:        true,
				Flags: []urfave.Flag{
					&urfave.BoolFlag{
						Name:    "start",
						Aliases: []string{"s"},
						Usage:   "Start the server after it has been deployed",
					},
				},
				Action: func(c *urfave.Context) error {
					return cliDeployServer(c)
				},
			},
			{
				Name:        "delete",
				Aliases:     []string{"del", "remove", "rm"},
				Description: "Delete a server",
				Usage:       "conductor delete [server name]",
				Flags: []urfave.Flag{
					&urfave.StringFlag{
						Name:    "yes",
						Aliases: []string{"y"},
						Usage:   "skip the confirmation",
					},
				},
				Action: func(c *urfave.Context) error {
					return cliDeleteServer(c)
				},
			},
			{
				Name:        "set",
				Description: "Set a servers global variables",
				Usage:       "conductor set [server name] [variable overrides]",
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
