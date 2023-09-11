package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"docker-cleaner/internal/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

var composeProjectLabel = "com.docker.compose.project"

func findComposeProjectByLabelsRegex(labels map[string]string) (string, error) {
	for labelName, labelValue := range labels {
		if patterns, ok := config.Config.MatchContainersLabelsRegex[labelName]; ok {
			for _, r := range patterns {
				if r.MatchString(labelValue) {
					return labels[composeProjectLabel], nil
				}
			}
		}
	}
	return "", errors.New("no containers found")
}

func getDockerNetworksByComposeProject(ctx context.Context, cli *client.Client, project string) []types.NetworkResource {
	networksFilter := filters.NewArgs(
		filters.Arg("label", fmt.Sprintf("%s=%s", composeProjectLabel, project)),
	)
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{
		Filters: networksFilter,
	})
	if err != nil {
		log.Fatal(err)
	}
	return networks
}

func getDockerVolumesByComposeProject(ctx context.Context, cli *client.Client, project string) []*volume.Volume {
	volumesFilter := filters.NewArgs(
		filters.Arg("label", fmt.Sprintf("%s=%s", composeProjectLabel, project)),
	)
	volumes, err := cli.VolumeList(ctx, volume.ListOptions{
		Filters: volumesFilter,
	})
	if err != nil {
		log.Fatal(err)
	}
	return volumes.Volumes
}

func Cleaner() {
	cli := Client()
	ctx := context.Background()

	// Get all containers contains "composeProjectLabel" label
	containersFilter := filters.NewArgs(
		filters.Arg("label", composeProjectLabel),
	)
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: containersFilter,
	})
	if err != nil {
		log.Fatal(err)
	}

	var composeProjects = make(map[string]bool)
	// Match containers by regex from config
	for _, container := range containers {
		// Find containers matched for regex's defined in config
		composeProject, err := findComposeProjectByLabelsRegex(container.Labels)
		if err != nil {
			continue
		}

		// Skip matched containers if creation time is less than Config.ContainerTtlCreated
		if time.Now().Sub(time.Unix(container.Created, 0)) < config.Config.ContainerTtlCreated {
			continue
		}

		composeProjects[composeProject] = true
		containerName, _ := strings.CutPrefix(container.Names[0], "/")

		log.Printf(`Remove docker container "%s" in docker-compose project "%s"`, containerName, composeProject)
		if !config.CheckMode {
			if err := RemoveContainer(&ctx, cli, &container); err != nil {
				log.Printf(`error when deleting container: %s`, err)
			}
		}
	}

	for composeProject, _ := range composeProjects {
		log.Printf(`Found docker-compose project "%s" by label "%s"`, composeProject, composeProjectLabel)

		// Get all networks by compose project label
		networks := getDockerNetworksByComposeProject(ctx, cli, composeProject)
		for _, network := range networks {
			log.Printf(`Remove docker network "%s" in docker-compose project "%s"`, network.Name, composeProject)
			if !config.CheckMode {
				if err := RemoveNetwork(&ctx, cli, &network); err != nil {
					log.Printf(`error when deleting network: %s`, err)
				}
			}
		}

		// Get all volumes by compose project label
		volumes := getDockerVolumesByComposeProject(ctx, cli, composeProject)
		for _, vol := range volumes {
			log.Printf(`Remove docker volume "%s" in docker-compose project "%s"`, vol.Name, composeProject)
			if !config.CheckMode {
				if err := RemoveVolume(&ctx, cli, vol); err != nil {
					log.Printf(`error when deleting volume: %s`, err)
				}
			}
		}
	}
}
