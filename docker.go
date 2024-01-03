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
	"os"
)

func deployContainer(templateName string, name string, templateVars map[string]string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}
	defer cli.Close()

	template, err := parseTemplate(templateName)
	if err != nil {
		return urfave.Exit(err.Error(), 1)
	}

	out, err := cli.ImagePull(ctx, template.Info.Container, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	/*	volume, err := cli.VolumeCreate(ctx, volume.CreateOptions{})
		if err != nil {
			return nil
		}*/

	createdContainer, err := cli.ContainerCreate(ctx, &container.Config{
		Image: template.Info.Container,
		Tty:   true,
	}, nil, nil, nil, name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	var installCmd = parseTemplateVars(templateVars) +
		"mkdir " + template.Info.WorkingDir +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.Install +
		"\n" + template.Actions.Update +
		"\n" + template.Actions.Adduser +
		"\n" + saveTemplateVarsCmd(templateVars)

	err = runCommandInContainer(ctx, cli, createdContainer.ID, "root", installCmd)
	if err != nil {
		return err
	}

	return nil
}

func runCommandInContainer(ctx context.Context, cli *client.Client, containerId string, user string, cmd string) error {
	exec, err := cli.ContainerExecCreate(ctx, containerId, types.ExecConfig{
		User:         user,
		Cmd:          []string{"sh", "-c", cmd},
		Tty:          false,
		AttachStdout: true,
		AttachStderr: true,
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
