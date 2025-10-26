package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/response"
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
		response.AIServiceError(c, "Failed to divide chapter into scenes: "+err.Error())
		return
	}

	response.Success(c, scenes)
}

func (h *SceneHandler) Get(c *gin.Context) {
	id := c.Param("id")

	scene, err := h.sceneService.GetScene(c.Request.Context(), id)
	if err != nil {
		response.ResourceNotFound(c, "Scene not found: "+err.Error())
		return
	}

	response.Success(c, scene)
}

func (h *SceneHandler) ListByChapter(c *gin.Context) {
	chapterID := c.Param("chapter_id")

	scenes, err := h.sceneService.GetScenesByChapterID(c.Request.Context(), chapterID)
	if err != nil {
		response.InternalError(c, "Failed to list scenes: "+err.Error())
		return
	}

	response.Success(c, scenes)
}

func (h *SceneHandler) ListByNovel(c *gin.Context) {
	novelID := c.Param("novel_id")

	scenes, err := h.sceneService.GetScenesByNovelID(c.Request.Context(), novelID)
	if err != nil {
		response.InternalError(c, "Failed to list scenes: "+err.Error())
		return
	}

	response.Success(c, scenes)
}

func (h *SceneHandler) GeneratePrompt(c *gin.Context) {
	var req dto.GenerateScenePromptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	result, err := h.sceneService.GeneratePrompt(c.Request.Context(), &req)
	if err != nil {
		response.AIServiceError(c, "Failed to generate prompt: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *SceneHandler) GenerateBatchPrompts(c *gin.Context) {
	var req dto.GeneratePromptsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.sceneService.GenerateBatchPrompts(c.Request.Context(), &req); err != nil {
		response.AIServiceError(c, "Failed to generate batch prompts: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Batch prompts generated successfully", nil)
}

func (h *SceneHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.sceneService.DeleteScene(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete scene: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Scene deleted successfully", nil)
}
