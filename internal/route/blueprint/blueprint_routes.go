package blueprint

import (
	"net/http"
	"singularity/internal/manager"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup, blueprintManager *manager.BlueprintManager) {
	group.GET("/list", func(c *gin.Context) {
		ListAllBlueprintsHandler(c, blueprintManager)
	})

	group.GET("/:id", func(c *gin.Context) {
		GetBlueprintHandler(c, blueprintManager)
	})
}

func ListAllBlueprintsHandler(c *gin.Context, bm *manager.BlueprintManager) {
	blueprints := bm.GetAllBlueprints()
	c.JSON(http.StatusOK, blueprints)
}

func GetBlueprintHandler(c *gin.Context, bm *manager.BlueprintManager) {
	id := c.Param("id")

	blueprint, ok := bm.GetBlueprint(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Blueprint not found",
		})
		return
	}

	c.JSON(http.StatusOK, blueprint)
}
