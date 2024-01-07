package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	urfave "github.com/urfave/cli/v2"
	"io"
	"os"
)

func deployContainer(templateName string, serverName string, templateVars map[string]string) error {
	var directory = getAppDir() + "/servers/" + serverName
	if serverExists(serverName) {
		return errors.New("there is already a server with that name")
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
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

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}

	createdContainer, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: template.Info.Container,
			Tty:   true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: directory,
					Target: template.Info.WorkingDir,
				},
			},
		}, nil, nil, serverName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, createdContainer.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	var rootInstallCmd = parseTemplateVars(templateVars) +
		"mkdir " + template.Info.WorkingDir +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.RootInstall

	err = runCommandInContainer(createdContainer.ID, "root", rootInstallCmd)
	if err != nil {
		return err
	}

	var installCmd = parseTemplateVars(templateVars) +
		"\ncd " + template.Info.WorkingDir +
		"\n" + template.Actions.Install +
		"\n" + template.Actions.Update +
		"\n" + saveTemplateVarsCmd(templateVars)

	err = runCommandInContainer(createdContainer.ID, template.Info.User, installCmd)
	if err != nil {
		return err
	}

	return nil
}

func deleteContainer(serverName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerKill(ctx, serverName, "SIGKILL")
	if err != nil {
		return err
	}

	err = os.RemoveAll(getAppDir() + "/servers/" + serverName)
	if err != nil {
		return err
	}

	cli.ContainerRemove(ctx, serverName, types.ContainerRemoveOptions{RemoveVolumes: true})
	return nil
}

func runCommandInContainer(serverName string, user string, cmd string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	exec, err := cli.ContainerExecCreate(ctx, serverName, types.ExecConfig{
		User:         user,
		Cmd:          []string{"bash", "-c", cmd},
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
