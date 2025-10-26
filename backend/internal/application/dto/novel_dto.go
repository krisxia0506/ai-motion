package dto

import "time"

type UploadNovelRequest struct {
	Title   string `json:"title" binding:"required"`
	Author  string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type NovelResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Status       string    `json:"status"`
	WordCount    int       `json:"word_count"`
	ChapterCount int       `json:"chapter_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type NovelListResponse struct {
	Novels []*NovelResponse `json:"novels"`
	Total  int              `json:"total"`
	Offset int              `json:"offset"`
	Limit  int              `json:"limit"`
}

type ChapterResponse struct {
	ID            string    `json:"id"`
	ChapterNumber int       `json:"chapter_number"`
	Title         string    `json:"title"`
	WordCount     int       `json:"word_count"`
	CreatedAt     time.Time `json:"created_at"`
}

type GenerateMangaRequest struct {
	Title   string `json:"title" binding:"required"`
	Author  string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type MangaWorkflowResponse struct {
	NovelID        string `json:"novel_id"`
	Title          string `json:"title"`
	CharacterCount int    `json:"character_count"`
	SceneCount     int    `json:"scene_count"`
	Status         string `json:"status"`
	Message        string `json:"message"`
}
