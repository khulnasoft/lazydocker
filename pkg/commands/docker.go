package commands

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/jesseduffield/lazydocker/pkg/config"
	"github.com/jesseduffield/lazydocker/pkg/i18n"
	"github.com/sirupsen/logrus"
)

// DockerCommand is our main git interface
type DockerCommand struct {
	Log       *logrus.Entry
	OSCommand *OSCommand
	Tr        *i18n.Localizer
	Config    config.AppConfigurer
	Client    *client.Client
}

// NewDockerCommand it runs git commands
func NewDockerCommand(log *logrus.Entry, osCommand *OSCommand, tr *i18n.Localizer, config config.AppConfigurer) (*DockerCommand, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &DockerCommand{
		Log:       log,
		OSCommand: osCommand,
		Tr:        tr,
		Config:    config,
		Client:    cli,
	}, nil
}

// GetContainers returns a slice of docker containers
func (c *DockerCommand) GetContainers() ([]*Container, error) {
	containers, err := c.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	ownContainers := make([]*Container, len(containers))

	for i, container := range containers {
		serviceName, ok := container.Labels["com.docker.compose.service"]
		if !ok {
			serviceName = ""
			c.Log.Warn("Could not get service name from docker container")
		}
		ownContainers[i] = &Container{
			ID:          container.ID,
			Name:        strings.TrimLeft(container.Names[0], "/"),
			ServiceName: serviceName,
			Container:   container,
			Client:      c.Client,
			OSCommand:   c.OSCommand,
			Log:         c.Log,
		}
	}

	return ownContainers, nil
}

// GetImages returns a slice of docker images
func (c *DockerCommand) GetImages() ([]*Image, error) {
	images, err := c.Client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	ownImages := make([]*Image, len(images))

	for i, image := range images {
		// func (cli *Client) ImageHistory(ctx context.Context, imageID string) ([]image.HistoryResponseItem, error)

		history, err := c.Client.ImageHistory(context.Background(), image.ID)
		if err != nil {
			return nil, err
		}

		name := "none"
		tags := image.RepoTags
		if len(tags) > 0 {
			name = tags[0]
		}

		nameParts := strings.Split(name, ":")

		ownImages[i] = &Image{
			ID:        image.ID,
			Name:      nameParts[0],
			Tag:       nameParts[1],
			Image:     image,
			Client:    c.Client,
			OSCommand: c.OSCommand,
			Log:       c.Log,
			History:   history,
		}
	}

	return ownImages, nil
}

// PruneImages prunes images
func (c *DockerCommand) PruneImages() error {
	_, err := c.Client.ImagesPrune(context.Background(), filters.Args{})
	return err
}