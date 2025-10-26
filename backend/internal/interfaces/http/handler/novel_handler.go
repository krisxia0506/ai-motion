package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/interfaces/http/response"
)

type NovelHandler struct {
	novelService *service.NovelService
}

func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{novelService: novelService}
}

func (h *NovelHandler) Upload(c *gin.Context) {
	var req dto.UploadNovelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.InvalidParams(c, "Invalid request: "+err.Error())
		return
	}

	novel, err := h.novelService.UploadAndParse(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "Failed to upload and parse novel: "+err.Error())
		return
	}

	response.Success(c, novel)
}

func (h *NovelHandler) Get(c *gin.Context) {
	id := c.Param("id")

	novel, err := h.novelService.GetNovel(c.Request.Context(), id)
	if err != nil {
		response.ResourceNotFound(c, "Novel not found: "+err.Error())
		return
	}

	response.Success(c, novel)
}

func (h *NovelHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit > 100 {
		limit = 100
	}

	novels, err := h.novelService.ListNovels(c.Request.Context(), offset, limit)
	if err != nil {
		response.InternalError(c, "Failed to list novels: "+err.Error())
		return
	}

	response.Success(c, novels)
}

func (h *NovelHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.novelService.DeleteNovel(c.Request.Context(), id); err != nil {
		response.InternalError(c, "Failed to delete novel: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "Novel deleted successfully", nil)
}

func (h *NovelHandler) GetChapters(c *gin.Context) {
	novelID := c.Param("id")

	chapters, err := h.novelService.GetChapters(c.Request.Context(), novelID)
	if err != nil {
		response.InternalError(c, "Failed to get chapters: "+err.Error())
		return
	}

	response.Success(c, chapters)
}
