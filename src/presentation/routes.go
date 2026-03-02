package presentation

import (
	"github.com/gin-gonic/gin"
)

// InitRoutes wires the HTTP handlers to Gin router. Only project endpoints are
// registered in this boilerplate.

func InitRoutes(router *gin.Engine, handler Handler) {
	router.POST("/organize", handler.Organize)
}
