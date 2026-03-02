package presentation

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rubenbuelvas/game-organizer-be/src/application"
)

type Handler struct {
	organizerService application.IOrganizerService
}

func NewHandler(organizerService application.IOrganizerService) Handler {
	return Handler{
		organizerService: organizerService,
	}
}

func (handler *Handler) Organize(c *gin.Context) {
	var dto application.OrganizerDTO
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "unable to read request body"})
		return
	}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	teams, err := handler.organizerService.Organize(dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, teams)
}
