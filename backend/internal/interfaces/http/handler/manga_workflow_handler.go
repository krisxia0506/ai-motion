package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/response"
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
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	if req.Title == "" {
		response.InvalidParams(c, "Title is required")
		return
	}

	if req.Content == "" {
		response.InvalidParams(c, "Content is required")
		return
	}

	result, err := h.workflowService.GenerateMangaFromNovel(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "Failed to generate manga: "+err.Error())
		return
	}

	response.Success(c, result)
}
