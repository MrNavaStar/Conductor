package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	urfave "github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
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
		" && cd " + template.Info.WorkingDir +
		" && " + parseScript(template.Actions.Install) +
		" && " + parseScript(template.Actions.Adduser)

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

func main() {
	/*ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	var templateStr = "templates/mindustry.yml"

	template, err := parseTemplate(templateStr)
	if err != nil {
		log.Fatal(err)
	}

	templateVars, err := getTemplateVars(templateStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(template)

	err = deployContainer(ctx, cli, template, "stupid", templateVars)
	if err != nil {
		log.Fatal(err)
	}*/

	app := &urfave.App{
		Name:        "Conductor",
		Description: "Easily create and manage game servers in a dockerized environment",
		Usage:       "conductor [template name]",
		Action: func(c *urfave.Context) error {
			return cliGetTemplateVars(c)
		},
		Commands: []*urfave.Command{},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
