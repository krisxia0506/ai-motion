package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
)

type MangaWorkflowHandler struct {
	workflowService *service.MangaWorkflowService
}

func NewMangaWorkflowHandler(workflowService *service.MangaWorkflowService) *MangaWorkflowHandler {
	return &MangaWorkflowHandler{
		workflowService: workflowService,
	}
}

func (h *MangaWorkflowHandler) GenerateManga(c *gin.Context) {
	var req dto.GenerateMangaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Content is required",
		})
		return
	}

	response, err := h.workflowService.GenerateMangaFromNovel(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate manga",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
