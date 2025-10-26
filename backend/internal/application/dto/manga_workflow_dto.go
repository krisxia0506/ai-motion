package dto

import "time"

// Note: GenerateMangaRequest is defined in novel_dto.go

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID      string                 `json:"task_id"`
	Status      string                 `json:"status"`
	Progress    TaskProgressResponse   `json:"progress"`
	Result      *TaskResultResponse    `json:"result,omitempty"`
	Error       *TaskErrorResponse     `json:"error,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	FailedAt    *time.Time             `json:"failed_at,omitempty"`
}

// TaskProgressResponse 任务进度响应
type TaskProgressResponse struct {
	CurrentStep      string                       `json:"current_step"`
	CurrentStepIndex int                          `json:"current_step_index"`
	TotalSteps       int                          `json:"total_steps"`
	Percentage       int                          `json:"percentage"`
	Details          *TaskProgressDetailsResponse `json:"details,omitempty"`
}

// TaskProgressDetailsResponse 任务进度详情
type TaskProgressDetailsResponse struct {
	CharactersExtracted int `json:"characters_extracted"`
	CharactersGenerated int `json:"characters_generated"`
	ScenesDivided       int `json:"scenes_divided"`
	ScenesGenerated     int `json:"scenes_generated"`
}

// TaskResultResponse 任务结果响应
type TaskResultResponse struct {
	NovelID        string                       `json:"novel_id"`
	Title          string                       `json:"title"`
	CharacterCount int                          `json:"character_count"`
	SceneCount     int                          `json:"scene_count"`
	Characters     []TaskCharacterResponse      `json:"characters"`
	Scenes         []TaskSceneResponse          `json:"scenes"`
}

// TaskCharacterResponse 任务中的角色信息（简化版）
type TaskCharacterResponse struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ReferenceImageURL string `json:"reference_image_url"`
}

// TaskSceneResponse 任务中的场景信息（简化版）
type TaskSceneResponse struct {
	ID          string `json:"id"`
	SequenceNum int    `json:"sequence_num"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// TaskErrorResponse 任务错误信息
type TaskErrorResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	RetryAble  bool   `json:"retry_able"`
}

// TaskListItemResponse 任务列表项
type TaskListItemResponse struct {
	TaskID         string                `json:"task_id"`
	Title          string                `json:"title"`
	Status         string                `json:"status"`
	Progress       TaskProgressResponse  `json:"progress"`
	CharacterCount int                   `json:"character_count,omitempty"`
	SceneCount     int                   `json:"scene_count,omitempty"`
	Error          *TaskErrorResponse    `json:"error,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	CompletedAt    *time.Time            `json:"completed_at,omitempty"`
	FailedAt       *time.Time            `json:"failed_at,omitempty"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
