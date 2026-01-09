package server

import (
	"net/http"
	"singularity/internal/data"
	"singularity/internal/enum"
	"singularity/internal/manager"
	"singularity/internal/strategy"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup,
	blueprintManager *manager.BlueprintManager,
	serverManager *manager.ServerManager,
	createContainerStrategy strategy.CreateContainerStrategy,
	deleteContainerStrategy strategy.DeleteContainerStrategy) {

	group.GET("", func(context *gin.Context) {
		ListAllServersHandler(context, serverManager)
	})

	group.GET("/:id", func(context *gin.Context) {
		GetServerHandler(context, serverManager)
	})

	group.POST("", func(context *gin.Context) {
		CreateServerHandler(context, createContainerStrategy, blueprintManager, serverManager)
	})

	group.DELETE("/:id", func(context *gin.Context) {
		DeleteServerHandler(context, serverManager, deleteContainerStrategy)
	})

	group.PATCH("/:id/report", func(context *gin.Context) {
		UpdateServerReportHandler(context, serverManager)
	})

	group.PATCH("/:id/status", func(context *gin.Context) {
		UpdateServerStatusHandler(context, serverManager)
	})

	group.POST("/:id/restart", func(context *gin.Context) {
		RestartServerHandler(context, serverManager, deleteContainerStrategy, createContainerStrategy)
	})

}

func CreateServerHandler(
	ctx *gin.Context,
	createStrategy strategy.CreateContainerStrategy,
	blueprintManager *manager.BlueprintManager,
	serverManager *manager.ServerManager,
) {
	blueprintId := ctx.Query("blueprintId")
	if blueprintId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid blueprintId"})
		return
	}

	blueprint, found := blueprintManager.GetBlueprint(blueprintId)
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "blueprint not found"})
		return
	}

	server := strategy.CreateServer(*blueprint, serverManager)

	server.Status = enum.StatusCreating

	if !createStrategy.CreateContainer(server) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "container creation failed",
		})

		return
	}

	server.Status = enum.StatusActive
	serverManager.AddServer(server)

	ctx.JSON(http.StatusCreated, server)
}

func ListAllServersHandler(context *gin.Context, serverManager *manager.ServerManager) {
	servers := serverManager.GetAllServers()
	context.JSON(http.StatusOK, servers)
}

func GetServerHandler(context *gin.Context, serverManager *manager.ServerManager) {
	id := context.Param("id")

	server, _ := serverManager.GetServer(id)
	if server == nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Server not found",
		})
		return
	}

	context.JSON(http.StatusOK, server)
}

func DeleteServerHandler(
	ctx *gin.Context,
	serverManager *manager.ServerManager,
	deleteStrategy strategy.DeleteContainerStrategy,
) {
	id := ctx.Param("id")

	server, err := serverManager.GetServer(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	if !deleteStrategy.DeleteContainer(server) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete container",
		})
		return
	}

	serverManager.DeleteServer(id)
	ctx.JSON(http.StatusOK, server)
}

func UpdateServerReportHandler(context *gin.Context, serverManager *manager.ServerManager) {
	id := context.Param("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid server id"})
		return
	}

	server, _ := serverManager.GetServer(id)
	if server == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	var report data.ServerReport
	if err := context.ShouldBindJSON(&report); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid report payload"})
		return
	}

	server.Report = &report
	context.JSON(http.StatusOK, server)
}

func UpdateServerStatusHandler(ctx *gin.Context, sm *manager.ServerManager) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server id"})
		return
	}

	var payload struct {
		Status enum.Status `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if !payload.Status.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	if !sm.UpdateStatus(id, payload.Status) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": payload.Status})
}

func RestartServerHandler(
	ctx *gin.Context,
	serverManager *manager.ServerManager,
	deleteStrategy strategy.DeleteContainerStrategy,
	createStrategy strategy.CreateContainerStrategy,
) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid server id"})
		return
	}

	server, err := serverManager.GetServer(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	serverManager.UpdateStatus(id, enum.StatusRestarting)

	if !deleteStrategy.DeleteContainer(server) {
		serverManager.UpdateStatus(id, enum.StatusError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete container",
		})
		return
	}

	if !createStrategy.CreateContainer(server) {
		serverManager.UpdateStatus(id, enum.StatusError)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create container",
		})
		return
	}

	serverManager.UpdateStatus(id, enum.StatusRestarting)
	ctx.JSON(http.StatusOK, gin.H{"status": "restarted"})
}
