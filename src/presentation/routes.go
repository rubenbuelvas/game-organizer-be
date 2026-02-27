package presentation

import (
	"github.com/gin-gonic/gin"
)

// InitRoutes wires the HTTP handlers to Gin router. Only project endpoints are
// registered in this boilerplate.

func InitRoutes(router *gin.Engine, handler Handler) {
	projGroup := router.Group("/projects")
	projGroup.GET("", projectHandler.ListProjects)
	projGroup.GET(":id", projectHandler.GetProject)
	projGroup.POST("", projectHandler.CreateProject)
	projGroup.PUT(":id", projectHandler.UpdateProject)
	projGroup.DELETE(":id", projectHandler.DeleteProject)
}
