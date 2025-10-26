package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/response"
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
		response.AIServiceError(c, "Failed to extract characters: "+err.Error())
		return
	}

	response.Success(c, characters)
}

func (h *CharacterHandler) Get(c *gin.Context) {
	id := c.Param("id")

	character, err := h.characterService.GetCharacter(c.Request.Context(), id)
	if err != nil {
		response.ResourceNotFound(c, "Character not found: "+err.Error())
		return
	}

	response.Success(c, character)
}

func (h *CharacterHandler) ListByNovel(c *gin.Context) {
	novelID := c.Param("novel_id")

	characters, err := h.characterService.GetCharactersByNovelID(c.Request.Context(), novelID)
	if err != nil {
		response.InternalError(c, "Failed to list characters: "+err.Error())
		return
	}

	response.Success(c, characters)
}

func (h *CharacterHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCharacterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	character, err := h.characterService.UpdateCharacter(c.Request.Context(), id, &req)
	if err != nil {
		response.InternalError(c, "Failed to update character: "+err.Error())
		return
	}

	response.Success(c, character)
}

func (h *CharacterHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.characterService.DeleteCharacter(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete character: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Character deleted successfully", nil)
}

func (h *CharacterHandler) Merge(c *gin.Context) {
	var req struct {
		NovelID  string `json:"novel_id" binding:"required"`
		SourceID string `json:"source_id" binding:"required"`
		TargetID string `json:"target_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.characterService.MergeCharacters(c.Request.Context(), req.NovelID, req.SourceID, req.TargetID); err != nil {
		response.InternalError(c, "Failed to merge characters: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Characters merged successfully", nil)
}
