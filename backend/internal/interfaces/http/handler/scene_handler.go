package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
)

type SceneHandler struct {
	sceneService *service.SceneService
}

func NewSceneHandler(sceneService *service.SceneService) *SceneHandler {
	return &SceneHandler{sceneService: sceneService}
}

func (h *SceneHandler) DivideChapter(c *gin.Context) {
	chapterID := c.Param("chapter_id")

	scenes, err := h.sceneService.DivideChapter(c.Request.Context(), chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to divide chapter into scenes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scenes,
	})
}

func (h *SceneHandler) Get(c *gin.Context) {
	id := c.Param("id")

	scene, err := h.sceneService.GetScene(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Scene not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scene,
	})
}

func (h *SceneHandler) ListByChapter(c *gin.Context) {
	chapterID := c.Param("chapter_id")

	scenes, err := h.sceneService.GetScenesByChapterID(c.Request.Context(), chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list scenes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scenes,
	})
}

func (h *SceneHandler) ListByNovel(c *gin.Context) {
	novelID := c.Param("novel_id")

	scenes, err := h.sceneService.GetScenesByNovelID(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list scenes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": scenes,
	})
}

func (h *SceneHandler) GeneratePrompt(c *gin.Context) {
	var req dto.GenerateScenePromptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := h.sceneService.GeneratePrompt(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate prompt",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *SceneHandler) GenerateBatchPrompts(c *gin.Context) {
	var req dto.GeneratePromptsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.sceneService.GenerateBatchPrompts(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate batch prompts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch prompts generated successfully",
	})
}

func (h *SceneHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.sceneService.DeleteScene(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete scene",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Scene deleted successfully",
	})
}
