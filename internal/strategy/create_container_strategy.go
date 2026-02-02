package strategy

import (
	"context"
	"fmt"
	"singularity/gen/blueprint"
	"singularity/internal/data"
	"singularity/internal/docker"
	"singularity/internal/util"
	"strconv"
	"time"

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

func (strategy CreateContainerStrategy) EnsureOrCreateVolumes(bp blueprint.Blueprint) {
	ctx := context.Background()
	client := strategy.dockerService.Client

	filters := mobyClient.Filters{}
	for _, volume := range bp.Volumes {
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

	for _, requiredVolume := range bp.Volumes {
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
	image := blueprint.Image

	pullImageErr := util.PullImageIfNotExists(client, ctx, image)
	if pullImageErr != nil {
		fmt.Printf("An unexpected error occurred, cannot pull image %s: %s", image, pullImageErr)
		return false
	}

	strategy.EnsureOrCreateVolumes(blueprint)

	var binds []string
	for _, volume := range blueprint.Volumes {
		bindEntry := fmt.Sprintf("%s:%s", volume.Id, volume.TargetFolder)
		if volume.ReadOnly {
			bindEntry += ":ro"
		}
		binds = append(binds, bindEntry)
	}

	defaultServerEnv := []string{
		"SERVER_ID=" + server.Discriminator,
		"SERVER_PORT=" + strconv.Itoa(int(server.Port)),
	}
	environment := util.MergeMapValuesWithExtras(blueprint.Environment, defaultServerEnv)

	containerConfig := &container.Config{
		Image: image,
		Env:   environment,
	}

	if blueprint.Entrypoint != nil && len(*blueprint.Entrypoint) > 0 {
		containerConfig.Entrypoint = *blueprint.Entrypoint
	}

	if blueprint.Cmd != nil && len(*blueprint.Cmd) > 0 {
		containerConfig.Cmd = *blueprint.Cmd
	}

	if blueprint.WorkingDir != nil {
		containerConfig.WorkingDir = *blueprint.WorkingDir
	}

	if blueprint.UserId != nil {
		user := strconv.FormatUint(uint64(*blueprint.UserId), 10)
		if blueprint.GroupId != nil {
			user += ":" + strconv.FormatUint(uint64(*blueprint.GroupId), 10)
		}
		containerConfig.User = user
	}

	if blueprint.HealthCheck != nil {
		containerConfig.Healthcheck = &container.HealthConfig{
			Test:          blueprint.HealthCheck.Test,
			Interval:      time.Duration(blueprint.HealthCheck.Interval) * time.Second,
			Timeout:       time.Duration(blueprint.HealthCheck.Timeout) * time.Second,
			Retries:       int(blueprint.HealthCheck.Retries),
			StartPeriod:   time.Duration(blueprint.HealthCheck.StartPeriod) * time.Second,
			StartInterval: time.Duration(blueprint.HealthCheck.Interval) * time.Second,
		}
	}

	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(blueprint.NetworkMode),
		Binds:       binds,
		Resources: container.Resources{
			Memory:   int64(blueprint.Resources.Memory) * MB,
			NanoCPUs: int64(blueprint.Resources.Cpu * 1e9),
		},
	}

	if blueprint.Logging != nil {
		hostConfig.LogConfig = container.LogConfig{
			Type:   blueprint.Logging.Driver,
			Config: blueprint.Logging.Options,
		}
	}

	result, err := client.ContainerCreate(ctx, mobyClient.ContainerCreateOptions{
		Name:       server.Id(),
		Image:      image,
		Config:     containerConfig,
		HostConfig: hostConfig,
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
