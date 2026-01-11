package strategy

import (
	"math/rand"
	"singularity/internal/data"
	"singularity/internal/enum"
	"singularity/internal/manager"
)

func CreateServer(blueprint data.Blueprint, serverManager *manager.ServerManager) *data.Server {
	availablePort := getAvailablePort(serverManager)
	discriminator := GenerateDiscriminator()

	server := data.Server{
		Blueprint:     blueprint,
		Discriminator: discriminator,
		Port:          availablePort,
		Status:        enum.StatusActive,
		MetricsPort:   nil,
		Report:        nil,
	}

	state := serverManager.AddServer(&server)
	if !state {
		return nil
	}

	return &server
}


func getAvailablePort(serverManager *manager.ServerManager) int {
	var highestServerPort int = -1
	for _, server := range serverManager.GetAllServers() {
		if server.Port < 25565 {
			continue
		}

		if server.Port > highestServerPort {
			highestServerPort = server.Port
		}
	}

	if highestServerPort == -1 {
		return 25565
	}

	return highestServerPort + 1
}

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789"

func GenerateDiscriminator() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(result)
}
