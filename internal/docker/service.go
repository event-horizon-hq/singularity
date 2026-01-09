package docker

import (
	"github.com/moby/moby/client"
)

type Service struct {
	Client *client.Client
}

func NewDockerService(client *client.Client) *Service {
	return &Service{
		Client: client,
	}
}
