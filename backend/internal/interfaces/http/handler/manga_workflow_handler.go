package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/middleware"
)

type MangaWorkflowHandler struct {
	workflowService *service.MangaWorkflowService
}

func NewMangaWorkflowHandler(workflowService *service.MangaWorkflowService) *MangaWorkflowHandler {
	return &MangaWorkflowHandler{
		workflowService: workflowService,
	}
}

// GenerateManga 创建漫画生成任务
func (h *MangaWorkflowHandler) GenerateManga(c *gin.Context) {
	var req dto.GenerateMangaRequest

	// 1. 解析请求
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request",
			"error", err,
			"path", c.Request.URL.Path,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10001,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	slog.Info("Received manga generation request",
		"title", req.Title,
		"content_length", len(req.Content),
	)

	// 2. 获取当前用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		slog.Warn("Unauthorized request - no user ID found",
			"path", c.Request.URL.Path,
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    20001,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	slog.Info("Creating manga task",
		"user_id", userID,
		"title", req.Title,
	)

	// 3. 创建任务
	task, err := h.workflowService.CreateTask(c.Request.Context(), userID, &req)
	if err != nil {
		slog.Error("Failed to create task",
			"error", err,
			"user_id", userID,
			"title", req.Title,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50001,
			"message": "创建任务失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	slog.Info("Task created successfully",
		"task_id", task.ID,
		"user_id", userID,
	)

	// 4. 异步执行任务
	go h.workflowService.ExecuteTask(context.Background(), task.ID)

	// 5. 立即返回任务ID
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务已创建",
		"data": gin.H{
			"task_id":    task.ID,
			"status":     task.Status,
			"created_at": task.CreatedAt,
		},
	})
}

// GetTaskStatus 获取任务状态
func (h *MangaWorkflowHandler) GetTaskStatus(c *gin.Context) {
	// 1. 获取任务ID
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10001,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	// 2. 获取当前用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    20001,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	// 3. 获取任务状态
	taskStatus, err := h.workflowService.GetTaskStatus(c.Request.Context(), userID, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    10002,
			"message": "任务不存在",
			"data":    nil,
		})
		return
	}

	// 4. 返回任务状态
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    taskStatus,
	})
}

// GetTaskList 获取任务列表
func (h *MangaWorkflowHandler) GetTaskList(c *gin.Context) {
	// 1. 获取当前用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    20001,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	// 2. 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 3. 获取任务列表
	tasks, pagination, err := h.workflowService.GetTaskList(c.Request.Context(), userID, page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    50001,
			"message": "获取任务列表失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 4. 返回任务列表
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":      tasks,
			"pagination": pagination,
		},
	})
}

// CancelTask 取消任务
func (h *MangaWorkflowHandler) CancelTask(c *gin.Context) {
	// 1. 获取任务ID
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10001,
			"message": "任务ID不能为空",
			"data":    nil,
		})
		return
	}

	// 2. 获取当前用户ID
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    20001,
			"message": "未授权",
			"data":    nil,
		})
		return
	}

	// 3. 取消任务
	err := h.workflowService.CancelTask(c.Request.Context(), userID, taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10001,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// 4. 返回成功
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "任务已取消",
		"data": gin.H{
			"task_id": taskID,
			"status":  "cancelled",
		},
	})
}
