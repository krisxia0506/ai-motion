package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
)

type CharacterHandler struct {
	characterService *service.CharacterService
}

func NewCharacterHandler(characterService *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{characterService: characterService}
}

func (h *CharacterHandler) Extract(c *gin.Context) {
	novelID := c.Param("novel_id")

	characters, err := h.characterService.ExtractCharacters(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to extract characters",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": characters,
	})
}

func (h *CharacterHandler) Get(c *gin.Context) {
	id := c.Param("id")

	character, err := h.characterService.GetCharacter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Character not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": character,
	})
}

func (h *CharacterHandler) ListByNovel(c *gin.Context) {
	novelID := c.Param("novel_id")

	characters, err := h.characterService.GetCharactersByNovelID(c.Request.Context(), novelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list characters",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": characters,
	})
}

func (h *CharacterHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCharacterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	character, err := h.characterService.UpdateCharacter(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update character",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": character,
	})
}

func (h *CharacterHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.characterService.DeleteCharacter(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete character",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Character deleted successfully",
	})
}

func (h *CharacterHandler) Merge(c *gin.Context) {
	var req struct {
		NovelID  string `json:"novel_id" binding:"required"`
		SourceID string `json:"source_id" binding:"required"`
		TargetID string `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	if err := h.characterService.MergeCharacters(c.Request.Context(), req.NovelID, req.SourceID, req.TargetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to merge characters",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Characters merged successfully",
	})
}
