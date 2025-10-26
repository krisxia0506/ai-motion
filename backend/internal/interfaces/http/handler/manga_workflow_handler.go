package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/middleware"
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
		slog.Error("Failed to bind request",
			"error", err,
			"path", c.Request.URL.Path,
		)
		response.InvalidParams(c, "参数错误: "+err.Error())
		return
	}

	slog.Info("Received manga generation request",
		"title", req.Title,
		"content_length", len(req.Content),
	)

	userID, exists := middleware.GetUserID(c)
	if !exists {
		slog.Warn("Unauthorized request - no user ID found",
			"path", c.Request.URL.Path,
		)
		response.Unauthorized(c, "未授权")
		return
	}

	slog.Info("Creating manga task",
		"user_id", userID,
		"title", req.Title,
	)

	task, err := h.workflowService.CreateTask(c.Request.Context(), userID, &req)
	if err != nil {
		slog.Error("Failed to create task",
			"error", err,
			"user_id", userID,
			"title", req.Title,
		)
		response.DatabaseError(c, "创建任务失败: "+err.Error())
		return
	}

	slog.Info("Task created successfully",
		"task_id", task.ID,
		"user_id", userID,
	)

	go h.workflowService.ExecuteTask(context.Background(), task.ID)

	response.SuccessWithMessage(c, "任务已创建", gin.H{
		"task_id":    task.ID,
		"status":     task.Status,
		"created_at": task.CreatedAt,
	})
}

func (h *MangaWorkflowHandler) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.InvalidParams(c, "任务ID不能为空")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	taskStatus, err := h.workflowService.GetTaskStatus(c.Request.Context(), userID, taskID)
	if err != nil {
		response.ResourceNotFound(c, "任务不存在")
		return
	}

	response.Success(c, taskStatus)
}

func (h *MangaWorkflowHandler) GetTaskList(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tasks, pagination, err := h.workflowService.GetTaskList(c.Request.Context(), userID, page, pageSize, status)
	if err != nil {
		response.DatabaseError(c, "获取任务列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":      tasks,
		"pagination": pagination,
	})
}

func (h *MangaWorkflowHandler) CancelTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		response.InvalidParams(c, "任务ID不能为空")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	err := h.workflowService.CancelTask(c.Request.Context(), userID, taskID)
	if err != nil {
		response.InvalidParams(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "任务已取消", gin.H{
		"task_id": taskID,
		"status":  "cancelled",
	})
}
