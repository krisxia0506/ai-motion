package dto

import "time"

type GenerateImageRequest struct {
	SceneID        string `json:"scene_id" binding:"required"`
	Prompt         string `json:"prompt" binding:"required"`
	NegativePrompt string `json:"negative_prompt"`
	ReferenceImage string `json:"reference_image"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Quality        string `json:"quality"`
	Style          string `json:"style"`
}

type GenerateVideoRequest struct {
	SceneID  string `json:"scene_id" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	Prompt   string `json:"prompt"`
	Duration int    `json:"duration"`
}

type BatchGenerateRequest struct {
	SceneIDs       []string `json:"scene_ids" binding:"required"`
	GenerateImages bool     `json:"generate_images"`
	GenerateVideos bool     `json:"generate_videos"`
}

type MediaResponse struct {
	ID           string        `json:"id"`
	SceneID      string        `json:"scene_id"`
	Type         string        `json:"type"`
	Status       string        `json:"status"`
	URL          string        `json:"url"`
	Metadata     MediaMetadata `json:"metadata"`
	GenerationID string        `json:"generation_id"`
	ErrorMessage string        `json:"error_message,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	CompletedAt  *time.Time    `json:"completed_at,omitempty"`
}

type MediaMetadata struct {
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Duration   float64 `json:"duration,omitempty"`
	Format     string  `json:"format"`
	FileSize   int64   `json:"file_size"`
	Resolution string  `json:"resolution"`
}

type GenerationStatusResponse struct {
	TotalTasks     int             `json:"total_tasks"`
	CompletedTasks int             `json:"completed_tasks"`
	FailedTasks    int             `json:"failed_tasks"`
	PendingTasks   int             `json:"pending_tasks"`
	Media          []MediaResponse `json:"media"`
}

type BatchGenerateResponse struct {
	JobID       string    `json:"job_id"`
	TotalScenes int       `json:"total_scenes"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
