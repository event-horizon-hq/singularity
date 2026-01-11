package metrics

import (
	"fmt"
	"net/http"
	"singularity/internal/manager"

	"github.com/gin-gonic/gin"
)

type PrometheusTargetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func Register(group *gin.RouterGroup, serverManager *manager.ServerManager) {
	group.GET("targets", func(ctx *gin.Context) {
		GetTargetsHandler(ctx, serverManager)
	})
}

func GetTargetsHandler(context *gin.Context, serverManager *manager.ServerManager) {
	var targets []string
	for _, server := range serverManager.GetAllServers() {
		targets = append(targets, fmt.Sprintf("%s:%d", "127.0.0.1:", *server.MetricsPort))
	}

	response := []PrometheusTargetGroup{
		{
			Targets: targets,
			Labels: map[string]string{
				"service": "hytale",
			},
		},
	}

	context.JSON(http.StatusOK, response)
}
