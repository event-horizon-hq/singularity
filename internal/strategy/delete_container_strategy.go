package strategy

import (
	"context"
	"fmt"
	"singularity/internal/data"
	"singularity/internal/docker"

	moby "github.com/moby/moby/client"
)

type DeleteContainerStrategy struct {
	docker *docker.Service
}

func CreateNewDeleteContainerStrategy(service *docker.Service) DeleteContainerStrategy {
	return DeleteContainerStrategy{
		service,
	}
}

func (strategy DeleteContainerStrategy) DeleteContainer(server *data.Server) bool {
	ctx := context.Background()
	client := strategy.docker.Client

	_, err := client.ContainerStop(ctx, server.Id(), moby.ContainerStopOptions{})
	if err != nil {
		fmt.Printf("an unexpected error occurred while stopping container. %s", err)
		return false
	}

	_, removeErr := client.ContainerRemove(ctx, server.Id(), moby.ContainerRemoveOptions{})
	if removeErr != nil {
		fmt.Printf("an unexpected error occurred while deleting container. %s", err)
		return false
	}

	return true
}
