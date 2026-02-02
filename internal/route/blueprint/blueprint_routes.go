package blueprint

import (
	"net/http"
	"singularity/internal/dto"
	"singularity/internal/manager"

	"github.com/gin-gonic/gin"
)

func Register(group *gin.RouterGroup, blueprintManager *manager.BlueprintManager, blueprintsFolder string) {
	group.GET("/list", func(c *gin.Context) {
		ListAllBlueprintsHandler(c, blueprintManager)
	})

	group.GET("/:id", func(c *gin.Context) {
		GetBlueprintHandler(c, blueprintManager)
	})

	group.POST("/reload", func(c *gin.Context) {
		ReloadBlueprintsHandler(c, blueprintManager, blueprintsFolder)
	})
}

func ListAllBlueprintsHandler(c *gin.Context, bm *manager.BlueprintManager) {
	blueprints := bm.GetAllBlueprints()
	response := make([]dto.BlueprintResponse, len(blueprints))
	for i, bp := range blueprints {
		response[i] = dto.NewBlueprintResponse(bp)
	}
	c.JSON(http.StatusOK, response)
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

	c.JSON(http.StatusOK, dto.NewBlueprintResponse(blueprint))
}

func ReloadBlueprintsHandler(c *gin.Context, bm *manager.BlueprintManager, blueprintsFolder string) {
	count, err := bm.ReloadBlueprints(blueprintsFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reloaded": count,
	})
}
