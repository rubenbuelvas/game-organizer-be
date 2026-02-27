package presentation

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

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

func (h *Handler) GetProjectsHandler(c *gin.Context) {
	result, err := h.taskService.GetTasks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (handler *TaskHandler) GetSubtasksHandler(c *gin.Context) {
	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}
	result, err := handler.taskService.GetSubtasks(taskId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (handler *TaskHandler) CreateTaskHandler(c *gin.Context) {
	var task application.CreateTaskDTO
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "unable to read request body"})
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	err = handler.taskService.CreateTask(task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func (handler *TaskHandler) PatchTaskHandler(c *gin.Context) {
	var task application.PatchTaskDTO
	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "unable to read request body"})
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	err = handler.taskService.PatchTask(taskId, task)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func (handler *TaskHandler) DeleteTaskHandler(c *gin.Context) {
	taskId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid task ID"})
		return
	}
	err = handler.taskService.DeleteTask(taskId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
