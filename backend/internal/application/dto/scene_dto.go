package dto

import "time"

type SceneResponse struct {
	ID           string              `json:"id"`
	ChapterID    string              `json:"chapter_id"`
	NovelID      string              `json:"novel_id"`
	SceneNumber  int                 `json:"scene_number"`
	Location     string              `json:"location,omitempty"`
	TimeOfDay    string              `json:"time_of_day,omitempty"`
	Description  DescriptionResponse `json:"description"`
	Dialogues    []DialogueResponse  `json:"dialogues"`
	CharacterIDs []string            `json:"character_ids"`
	ImagePrompt  string              `json:"image_prompt,omitempty"`
	VideoPrompt  string              `json:"video_prompt,omitempty"`
	Status       string              `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type DescriptionResponse struct {
	Setting    string `json:"setting,omitempty"`
	Action     string `json:"action,omitempty"`
	Atmosphere string `json:"atmosphere,omitempty"`
	FullText   string `json:"full_text"`
}

type DialogueResponse struct {
	Speaker string `json:"speaker"`
	Content string `json:"content"`
	Emotion string `json:"emotion,omitempty"`
}

type DivideChapterRequest struct {
	ChapterID string `json:"chapter_id" binding:"required"`
}

type GeneratePromptsRequest struct {
	SceneIDs    []string `json:"scene_ids" binding:"required"`
	Style       string   `json:"style,omitempty"`
	Quality     string   `json:"quality,omitempty"`
	AspectRatio string   `json:"aspect_ratio,omitempty"`
}

type GenerateScenePromptRequest struct {
	SceneID      string   `json:"scene_id" binding:"required"`
	CharacterIDs []string `json:"character_ids,omitempty"`
	Style        string   `json:"style,omitempty"`
	Quality      string   `json:"quality,omitempty"`
	AspectRatio  string   `json:"aspect_ratio,omitempty"`
}

type GenerateScenePromptResponse struct {
	SceneID     string    `json:"scene_id"`
	ImagePrompt string    `json:"image_prompt"`
	VideoPrompt string    `json:"video_prompt,omitempty"`
	GeneratedAt time.Time `json:"generated_at"`
}

type SceneListResponse struct {
	Scenes []*SceneResponse `json:"scenes"`
	Total  int              `json:"total"`
}
