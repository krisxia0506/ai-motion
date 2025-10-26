package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/response"
)

type GenerationHandler struct {
	generationService *service.GenerationService
}

func NewGenerationHandler(generationService *service.GenerationService) *GenerationHandler {
	return &GenerationHandler{
		generationService: generationService,
	}
}

func (h *GenerationHandler) GenerateImage(c *gin.Context) {
	var req dto.GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.generationService.GenerateSceneImage(c.Request.Context(), &req)
	if err != nil {
		response.GenerationError(c, "Failed to generate image: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *GenerationHandler) GenerateVideo(c *gin.Context) {
	var req dto.GenerateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.generationService.GenerateSceneVideo(c.Request.Context(), &req)
	if err != nil {
		response.GenerationError(c, "Failed to generate video: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *GenerationHandler) BatchGenerate(c *gin.Context) {
	var req dto.BatchGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.generationService.BatchGenerateScenes(c.Request.Context(), &req)
	if err != nil {
		response.GenerationError(c, "Failed to batch generate: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *GenerationHandler) GetStatus(c *gin.Context) {
	sceneID := c.Param("scene_id")
	if sceneID == "" {
		response.InvalidParams(c, "scene_id is required")
		return
	}

	result, err := h.generationService.GetGenerationStatus(c.Request.Context(), sceneID)
	if err != nil {
		response.InternalError(c, "Failed to get status: "+err.Error())
		return
	}

	response.Success(c, result)
}
