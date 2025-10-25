package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.generationService.GenerateSceneImage(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate image",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *GenerationHandler) GenerateVideo(c *gin.Context) {
	var req dto.GenerateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.generationService.GenerateSceneVideo(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate video",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *GenerationHandler) BatchGenerate(c *gin.Context) {
	var req dto.BatchGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	result, err := h.generationService.BatchGenerateScenes(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to batch generate",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

func (h *GenerationHandler) GetStatus(c *gin.Context) {
	sceneID := c.Param("scene_id")
	if sceneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "scene_id is required",
		})
		return
	}

	result, err := h.generationService.GetGenerationStatus(c.Request.Context(), sceneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
