package strategy

import (
	"context"
	"fmt"
	"singularity/internal/data"
	"singularity/internal/docker"
	"singularity/internal/util"
	"strconv"

	"github.com/moby/moby/api/types/container"
	mobyClient "github.com/moby/moby/client"
)

const MB = 1024 * 1024

type CreateContainerStrategy struct {
	dockerService *docker.Service
}

func CreateNewContainerStrategy(service *docker.Service) CreateContainerStrategy {
	return CreateContainerStrategy{
		service,
	}
}

func (strategy CreateContainerStrategy) EnsureOrCreateVolumes(blueprint data.Blueprint) {
	ctx := context.Background()
	client := strategy.dockerService.Client

	filters := mobyClient.Filters{}
	for _, volume := range blueprint.Volumes {
		filters.Add("name", volume.Id)
	}

	volumeList, err := client.VolumeList(ctx, mobyClient.VolumeListOptions{
		Filters: filters,
	})

	if err != nil {
		fmt.Printf("cannot list volumes from docker: %s\n", err)
		return
	}

	existingMap := make(map[string]bool)
	for _, v := range volumeList.Items {
		existingMap[v.Name] = true
	}

	for _, requiredVolume := range blueprint.Volumes {
		if !existingMap[requiredVolume.Id] {
			_, err := client.VolumeCreate(ctx, mobyClient.VolumeCreateOptions{
				Name:   requiredVolume.Id,
				Driver: "local",
			})

			if err != nil {
				fmt.Printf("An unexpected error occurred, cannot create volume %s: %s\n", requiredVolume.Id, err)
			} else {
				fmt.Println("The new volume has been created: ", requiredVolume.Id)
			}
		}
	}
}

func (strategy CreateContainerStrategy) CreateContainer(server *data.Server) bool {
	ctx := context.Background()
	client := strategy.dockerService.Client

	blueprint := server.Blueprint
	image := blueprint.Environment["image"]

	util.PullImageIfNotExists(client, ctx, image)
	strategy.EnsureOrCreateVolumes(blueprint)

	memoryAmount, err := strconv.ParseInt(blueprint.Environment["memory-amount"], 10, 64)
	if err != nil {
		fmt.Printf("An unexpected error occurred, cannot convert memory-amount property from blueprint env. %s", err)
		return false
	}

	var binds []string
	for _, volume := range blueprint.Volumes {
		bindEntry := fmt.Sprintf("%s:%s", volume.Id, volume.TargetFolder)
		binds = append(binds, bindEntry)
	}

	defaultServerEnv := []string{
		"SERVER_ID=" + server.Discriminator,
		"SERVER_PORT=" + strconv.Itoa(int(server.Port)),
	}

	environment := util.MergeMapValuesWithExtras(blueprint.Environment, defaultServerEnv)

	result, err := client.ContainerCreate(ctx, mobyClient.ContainerCreateOptions{
		Name:  server.Id(),
		Image: image,
		Config: &container.Config{
			Env: environment,
		},
		HostConfig: &container.HostConfig{
			NetworkMode: "host",
			Binds:       binds,
			Resources: container.Resources{
				Memory: memoryAmount * MB,
			},
		},
	})

	if err != nil {
		fmt.Printf("An unexpected error occurred, cannot create container. %s", err)
		return false
	}

	_, err = client.ContainerStart(ctx, result.ID, mobyClient.ContainerStartOptions{})
	if err != nil {
		return false
	}

	return true
}
