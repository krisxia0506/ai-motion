package media

import (
	"errors"
	"time"
)

type MediaID string
type MediaType string
type MediaStatus string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

const (
	MediaStatusPending    MediaStatus = "pending"
	MediaStatusGenerating MediaStatus = "generating"
	MediaStatusCompleted  MediaStatus = "completed"
	MediaStatusFailed     MediaStatus = "failed"
)

var (
	ErrMediaNotFound    = errors.New("media not found")
	ErrInvalidMediaType = errors.New("invalid media type")
	ErrInvalidStatus    = errors.New("invalid media status")
)

type Media struct {
	ID           MediaID
	NovelID      string // 关联的小说ID
	SceneID      string // 关联的场景ID（可选，用于场景生成模式）
	Type         MediaType
	Status       MediaStatus
	URL          string
	Metadata     MediaMetadata
	GenerationID string
	ErrorMessage string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  *time.Time
}

func NewMedia(sceneID string, mediaType MediaType) *Media {
	return &Media{
		ID:        MediaID(generateID()),
		SceneID:   sceneID,
		Type:      mediaType,
		Status:    MediaStatusPending,
		Metadata:  MediaMetadata{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewMediaForNovel 创建与小说关联的媒体实体（用于漫画生成等场景）
func NewMediaForNovel(novelID string, mediaType MediaType) *Media {
	return &Media{
		ID:        MediaID(generateID()),
		NovelID:   novelID,
		Type:      mediaType,
		Status:    MediaStatusPending,
		Metadata:  MediaMetadata{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Media) IsReady() bool {
	return m.Status == MediaStatusCompleted
}

func (m *Media) MarkGenerating(generationID string) {
	m.Status = MediaStatusGenerating
	m.GenerationID = generationID
	m.UpdatedAt = time.Now()
}

func (m *Media) MarkCompleted(url string, metadata MediaMetadata) {
	m.Status = MediaStatusCompleted
	m.URL = url
	m.Metadata = metadata
	now := time.Now()
	m.CompletedAt = &now
	m.UpdatedAt = now
}

func (m *Media) MarkFailed(errorMsg string) {
	m.Status = MediaStatusFailed
	m.ErrorMessage = errorMsg
	m.UpdatedAt = time.Now()
}

func (m *Media) Validate() error {
	if m.Type != MediaTypeImage && m.Type != MediaTypeVideo {
		return ErrInvalidMediaType
	}
	// Either NovelID or SceneID must be provided
	if m.NovelID == "" && m.SceneID == "" {
		return errors.New("either novel_id or scene_id is required")
	}
	return nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
