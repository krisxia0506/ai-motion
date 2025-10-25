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
	SceneID      string
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
	if m.SceneID == "" {
		return errors.New("scene_id is required")
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
