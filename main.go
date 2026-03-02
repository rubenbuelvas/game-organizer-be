package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/rubenbuelvas/game-organizer-be/src/application"
	"github.com/rubenbuelvas/game-organizer-be/src/presentation"
)

type Config struct {
	App struct {
		Route string
	}
}

func main() {
	config := loadConfig()
	engine := gin.Default()
	applyDependencies(engine, config)
	engine.Run(config.App.Route)
}

func applyDependencies(engine *gin.Engine, config Config) {
	service := application.NewOrganizerService()
	handler := presentation.NewHandler(service)
	presentation.InitRoutes(engine, handler)
}

func loadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}
