package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pterm/pterm"
	urfave "github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"strings"
)

func deployContainer(ctx context.Context, cli *client.Client, template Template, name string, templateVars map[string]string) error {
	out, err := cli.ImagePull(ctx, "docker.io/library/"+template.Info.Container, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	createdContainer, err := cli.ContainerCreate(ctx, &container.Config{
		Image: template.Info.Container,
		Tty:   true,
	}, nil, nil, nil, name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	var installCmd = parseTemplateVars(templateVars) +
		"mkdir " + template.Info.WorkingDir +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.Install +
		"\n" + template.Actions.Adduser

	println(installCmd)

	err = runCommandInContainer(ctx, cli, createdContainer.ID, "root", installCmd)
	if err != nil {
		return err
	}

	return nil
}

func runCommandInContainer(ctx context.Context, cli *client.Client, containerId string, user string, cmd string) error {
	exec, err := cli.ContainerExecCreate(ctx, containerId, types.ExecConfig{
		User: user,
		Cmd:  []string{"sh", "-c", cmd},
	})
	if err != nil {
		return err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), exec.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	scanner := bufio.NewScanner(resp.Reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

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
	/*templateName := c.Args().Get(0)
	if len(templateName) == 0 {
		return nil
	}

	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}
	defer docker.Close()

	template, err := parseTemplate(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	templateVars, err := getTemplateVars(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	err := os.MkdirAll("/var/lib/conductor", os.ModePerm)
	if err != nil {
		return nil
	}

	create, err := docker.VolumeCreate(ctx, volume.CreateOptions{})
	if err != nil {
		return nil
	}*/

	return nil
}

func cliStartServer(c *urfave.Context) urfave.ExitCoder {

	return nil
}

func main() {
	app := &urfave.App{
		Name:        "conductor",
		Description: "Easily create and manage game servers in a dockerized environment",
		Usage:       "conductor [template name]",
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
				Usage:       "conductor deploy [template name] [flags]",
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
