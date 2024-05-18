package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
)

type Server struct {
	Name     string
	Template string
	Vars     []string
	Running  bool
}

func (server Server) deploy(cli *client.Client) error {
	template, err := getTemplate(server.Template)
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(context.Background(), template.Info.Container, image.PullOptions{})
	if err != nil {
		return err
	}

	defer reader.Close()
	// cli.ImagePull is asynchronous.
	// The reader needs to be read completely for the pull operation to complete.
	// If stdout is not required, consider using io.Discard instead of os.Stdout.
	io.Copy(os.Stdout, reader)

	created, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: template.Info.Container,
			Cmd:   []string{"/bin/sh", "-c", template.getDeployCmd()},
		},
		&container.HostConfig{}, nil, nil, "conductor-"+server.Name)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), created.ID, container.StartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(context.Background(), created.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(context.Background(), created.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		return err
	}

	_, err = cli.ContainerCommit(context.Background(), created.ID, container.CommitOptions{
		Reference: "conductor-" + server.Name,
		Comment:   "Created from template: " + template.Name,
		Author:    "conductor",
	})
	if err != nil {
		return err
	}

	statusCh, errCh = cli.ContainerWait(context.Background(), created.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	err = cli.ContainerRemove(context.Background(), "conductor-"+server.Name, container.RemoveOptions{RemoveVolumes: true})
	if err != nil {
		return err
	}

	created, err = cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: "conductor-" + server.Name,
			Cmd:   []string{"/bin/sh", template.getStartCmd()},
			User:  template.Info.User,
		},
		&container.HostConfig{}, nil, nil, "conductor-"+server.Name)
	if err != nil {
		return err
	}

	statusCh, errCh = cli.ContainerWait(context.Background(), created.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	_, err = cli.ContainerCommit(context.Background(), created.ID, container.CommitOptions{
		Reference: "conductor-" + server.Name,
		Comment:   "Created from template: " + template.Name,
		Author:    "conductor",
	})
	if err != nil {
		return err
	}

	if server.Running {
		err := server.start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (server Server) start() error {
	return nil
}

func (server Server) stop() error {
	return nil
}

func (server Server) delete(cli *client.Client) error {
	err := server.stop()
	if err != nil {
		return err
	}

	err = cli.ContainerKill(context.Background(), "conductor-"+server.Name, "SIGKILL")
	if err != nil {
		return err
	}

	err = os.RemoveAll(getAppDir() + "/servers/" + server.Name)
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(context.Background(), "conductor-"+server.Name, container.RemoveOptions{RemoveVolumes: true})
	if err != nil {
		return err
	}
	return nil
}

func (server Server) backup() error {
	return nil
}
