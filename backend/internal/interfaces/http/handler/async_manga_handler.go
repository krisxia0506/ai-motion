package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
)

type AsyncMangaHandler struct {
	service *service.AsyncMangaService
}

func NewAsyncMangaHandler(service *service.AsyncMangaService) *AsyncMangaHandler {
	return &AsyncMangaHandler{
		service: service,
	}
}

func (h *AsyncMangaHandler) Generate(c *gin.Context) {
	var req dto.GenerateMangaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	task, err := h.service.StartMangaGeneration(c.Request.Context(), req.Title, req.Author, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start manga generation",
			"details": err.Error(),
		})
		return
	}

	response := &dto.TaskResponse{
		TaskID:    task.ID,
		Type:      string(task.Type),
		Status:    string(task.Status),
		Progress:  task.Progress,
		Message:   task.Message,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data": response,
		"message": "Manga generation started. Use task_id to check status.",
	})
}

func (h *AsyncMangaHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")

	task, err := h.service.GetTaskStatus(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Task not found",
			"details": err.Error(),
		})
		return
	}

	response := &dto.TaskResponse{
		TaskID:      task.ID,
		Type:        string(task.Type),
		Status:      string(task.Status),
		Progress:    task.Progress,
		Message:     task.Message,
		Result:      task.Result,
		Error:       task.Error,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
