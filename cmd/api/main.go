package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"singularity/internal/auth"
	"singularity/internal/auth/middleware"
	"singularity/internal/config"
	"singularity/internal/docker"
	"singularity/internal/factory"
	"singularity/internal/manager"
	"singularity/internal/repository"
	"singularity/internal/route/blueprint"
	"singularity/internal/route/metrics"
	"singularity/internal/route/server"
	"singularity/internal/strategy"

	"github.com/gin-gonic/gin"
	moby "github.com/moby/moby/client"
)

func main() {
	configuration := ReadConfigData()
	authenticationService := auth.NewAuthenticationService(configuration.JwtSecretKey)

	client, err := factory.InitMongo(configuration.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("singularity")

	blueprintManager := manager.CreateNewBlueprintManager()
	serverRepository := repository.NewServerRepository(database)

	serverRepository.EnsureIndexes(context.Background())

	serverManager := manager.CreateNewServerManager(serverRepository)

	dockerClient, dockerErr := moby.New()
	if dockerErr != nil {
		log.Fatal(dockerErr)
		return
	}

	dockerService := docker.NewDockerService(dockerClient)

	RegisterServers(serverRepository, serverManager)

	ReadBlueprints(configuration, blueprintManager)
	ReadAccessToken(authenticationService)
	StartRouter(
		authenticationService,
		blueprintManager,
		serverManager,
		dockerService,
		configuration)
}

func RegisterServers(serverRepository *repository.ServerRepository, serverManager *manager.ServerManager) {
	ctx := context.Background()

	servers, err := serverRepository.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, foundServer := range servers {
		ok := serverManager.LoadServer(foundServer)
		if !ok {
			fmt.Println("Duplicated data in foundServer manager. Server ID: " + foundServer.Id())
		}
	}
}

func ReadAccessToken(authentication *auth.AuthenticationService) {
	tokenFile := "access_token"

	log.Println("Verifying if access_token exists...")
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		jwtToken, err := authentication.GenerateSystemToken()
		if err != nil {
			log.Fatalf("Can't generate new token because: %v", err)
		}

		if err := os.WriteFile(tokenFile, []byte(jwtToken), 0600); err != nil {
			log.Fatalf("Can't save file new token because: %v", err)
		}

		log.Println("The new token has been generated and save on access_token file.")
	} else {
		log.Println("The access_token file already exists.")
	}
}

func ReadConfigData() *config.Config {
	cfg, configErr := config.Load("./config.toml")
	if configErr != nil {
		log.Fatal(configErr)
		return nil
	}

	log.Println("Config loaded successfully.")
	return cfg
}

func ReadBlueprints(config *config.Config, blueprintManager *manager.BlueprintManager) {

	if _, blueprintErr := blueprintManager.LoadBlueprints(config.BlueprintsFolder); blueprintErr != nil {
		log.Fatal(blueprintErr)
	}

	log.Println("Blueprints loaded successfully.")
}

func StartRouter(authenticationService *auth.AuthenticationService,
	blueprintManager *manager.BlueprintManager,
	serverManager *manager.ServerManager,
	dockerService *docker.Service,
	config *config.Config,
) {
	router := gin.Default()
	trustedProxiesErr := router.SetTrustedProxies(config.TrustedProxies)

	if trustedProxiesErr != nil {
		log.Fatal(trustedProxiesErr)
		return
	}

	createContainerStrategy := strategy.CreateNewContainerStrategy(dockerService)
	deleteContainerStrategy := strategy.CreateNewDeleteContainerStrategy(dockerService)

	metricsGroup := router.Group("/v1/metrics")
	metricsGroup.Use(middleware.Auth(authenticationService))
	metrics.Register(metricsGroup, serverManager)

	blueprintGroup := router.Group("/v1/blueprints")
	blueprintGroup.Use(middleware.Auth(authenticationService))
	blueprint.Register(blueprintGroup, blueprintManager, config.BlueprintsFolder)

	serverGroup := router.Group("/v1/servers")
	serverGroup.Use(middleware.Auth(authenticationService))
	server.Register(
		serverGroup,
		blueprintManager,
		serverManager,
		createContainerStrategy,
		deleteContainerStrategy)

	err := router.Run()
	if err != nil {
		return
	}
}
