package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func Client() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer func(cli *client.Client) {
		err := cli.Close()
		if err != nil {
			log.Fatal("fatal error connecting to docker: %w", err)
		}
	}(cli)

	return cli
}

func RemoveContainer(ctx *context.Context, client *client.Client, container *types.Container) error {
	options := types.ContainerRemoveOptions{
		Force: true,
	}

	return client.ContainerRemove(*ctx, container.ID, options)
}

func RemoveNetwork(ctx *context.Context, client *client.Client, network *types.NetworkResource) error {
	return client.NetworkRemove(*ctx, network.ID)
}

func RemoveVolume(ctx *context.Context, client *client.Client, volume *volume.Volume) error {
	return client.VolumeRemove(*ctx, volume.Name, true)
}
