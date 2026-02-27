package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rubenbuelvas/ringover/src/application"
	"github.com/rubenbuelvas/ringover/src/infrastructure"
	"github.com/rubenbuelvas/ringover/src/presentation"
)

func main() {
	// initialize repository, service and handler
	repo := infrastructure.NewMemoryProjectRepository()
	service := application.NewProjectService(repo)
	handler := presentation.NewProjectHandler(service)

	// prepare Gin router
	r := gin.Default()
	presentation.InitRoutes(r, handler)

	// start server
	r.Run() // listen on :8080 by default
}
