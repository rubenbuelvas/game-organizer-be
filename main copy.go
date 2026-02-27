package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/rubenbuelvas/ringover/src/application"
	"github.com/rubenbuelvas/ringover/src/infrastructure"
	"github.com/rubenbuelvas/ringover/src/presentation"
)

type Config struct {
	App struct {
		Route string
	}
	DB struct {
		User   string
		DBName string
		DSN    string
	}
}

func main() {
	config := loadConfig()
	engine := gin.Default()
	applyDependencies(engine, config)
	engine.Run(config.App.Route)
}

func applyDependencies(engine *gin.Engine, config Config) {
	MySQLClient, err := infrastructure.NewMySQLClient(config.DB.DSN)
	if err != nil {
		panic(err)
	}
	taskService := application.NewTaskService(MySQLClient)
	taskHandler := presentation.NewHandler(taskService)
	presentation.InitRoutes(engine, taskHandler)
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
