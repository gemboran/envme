package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/viper"
	"io"
	"os"
)

// RunContainer runs a container with the given image and environment variables
// It's equivalent to running `docker run -e PORT=8080 backend:latest` or with `-d` flag
func RunContainer(containerName, image, networkName string, detach bool) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, viper.GetString("image"), types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	if err := EnsureNetworkExists(ctx, cli, networkName); err != nil {
		return err
	}

	hostConfig := &container.HostConfig{
		ExtraHosts:  []string{"host.docker.internal:host-gateway"},
		NetworkMode: container.NetworkMode(networkName),
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Env:   viper.GetStringSlice("env"),
	}, hostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	if !detach {
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
		if err != nil {
			return err
		}

		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}

	return nil
}

// ListContainers lists all the running containers
// It's equivalent to running `docker ps`
func ListContainers() ([]types.Container, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// StopContainer stops a running container
// It's equivalent to running `docker stop <container_id>`
func StopContainer(containerID string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return err
	}

	return nil
}

// RemoveContainer removes a container
// It's equivalent to running `docker rm <container_id>`
func RemoveContainer(containerID string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}
