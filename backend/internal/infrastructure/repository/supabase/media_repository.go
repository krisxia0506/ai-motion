package supabase

import (
	"context"
	"fmt"
	"time"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/media"
)

type MediaRepository struct {
	client *postgrest.Client
}

func NewMediaRepository(client *postgrest.Client) media.MediaRepository {
	return &MediaRepository{client: client}
}

func (r *MediaRepository) Save(ctx context.Context, m *media.Media) error {
	data := map[string]interface{}{
		"id":             string(m.ID),
		"novel_id":       m.NovelID,
		"scene_id":       m.SceneID,
		"type":           string(m.Type),
		"status":         string(m.Status),
		"url":            m.URL,
		"width":          m.Metadata.Width,
		"height":         m.Metadata.Height,
		"duration":       m.Metadata.Duration,
		"format":         m.Metadata.Format,
		"file_size":      m.Metadata.FileSize,
		"generation_id":  m.GenerationID,
		"error_message":  m.ErrorMessage,
		"created_at":     m.CreatedAt,
		"updated_at":     m.UpdatedAt,
		"completed_at":   m.CompletedAt,
	}

	_, _, err := r.client.From("aimotion_media").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save media: %w", err)
	}

	return nil
}

func (r *MediaRepository) FindByID(ctx context.Context, id media.MediaID) (*media.Media, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_media").
		Select("*", "", false).
		Eq("id", string(id)).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	if len(results) == 0 {
		return nil, media.ErrMediaNotFound
	}

	return r.mapToMedia(results[0])
}

func (r *MediaRepository) FindBySceneID(ctx context.Context, sceneID string) ([]*media.Media, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_media").
		Select("*", "", false).
		Eq("scene_id", sceneID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find media by scene ID: %w", err)
	}

	var mediaList []*media.Media
	for _, result := range results {
		m, err := r.mapToMedia(result)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, m)
	}

	return mediaList, nil
}

func (r *MediaRepository) FindByNovelID(ctx context.Context, novelID string) ([]*media.Media, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_media").
		Select("*", "", false).
		Eq("novel_id", novelID).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find media by novel ID: %w", err)
	}

	var mediaList []*media.Media
	for _, result := range results {
		m, err := r.mapToMedia(result)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, m)
	}

	return mediaList, nil
}

func (r *MediaRepository) UpdateStatus(ctx context.Context, id media.MediaID, status media.MediaStatus, url string, errorMsg string) error {
	data := map[string]interface{}{
		"status":        string(status),
		"url":           url,
		"error_message": errorMsg,
		"updated_at":    time.Now(),
	}

	_, _, err := r.client.From("aimotion_media").
		Update(data, "", "").
		Eq("id", string(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to update media status: %w", err)
	}

	return nil
}

func (r *MediaRepository) Delete(ctx context.Context, id media.MediaID) error {
	_, _, err := r.client.From("aimotion_media").
		Delete("", "").
		Eq("id", string(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	return nil
}

func (r *MediaRepository) FindPendingMedia(ctx context.Context, limit int) ([]*media.Media, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_media").
		Select("*", "", false).
		Eq("status", string(media.MediaStatusPending)).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		Range(0, limit-1, "").
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find pending media: %w", err)
	}

	var mediaList []*media.Media
	for _, result := range results {
		m, err := r.mapToMedia(result)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, m)
	}

	return mediaList, nil
}

func (r *MediaRepository) mapToMedia(data map[string]interface{}) (*media.Media, error) {
	m := &media.Media{}

	if id, ok := data["id"].(string); ok {
		m.ID = media.MediaID(id)
	}
	if novelID, ok := data["novel_id"].(string); ok {
		m.NovelID = novelID
	}
	if sceneID, ok := data["scene_id"].(string); ok {
		m.SceneID = sceneID
	}
	if mediaType, ok := data["type"].(string); ok {
		m.Type = media.MediaType(mediaType)
	}
	if status, ok := data["status"].(string); ok {
		m.Status = media.MediaStatus(status)
	}
	if url, ok := data["url"].(string); ok {
		m.URL = url
	}
	if width, ok := data["width"].(float64); ok {
		m.Metadata.Width = int(width)
	}
	if height, ok := data["height"].(float64); ok {
		m.Metadata.Height = int(height)
	}
	if duration, ok := data["duration"].(float64); ok {
		m.Metadata.Duration = duration
	}
	if format, ok := data["format"].(string); ok {
		m.Metadata.Format = format
	}
	if fileSize, ok := data["file_size"].(float64); ok {
		m.Metadata.FileSize = int64(fileSize)
	}
	if generationID, ok := data["generation_id"].(string); ok {
		m.GenerationID = generationID
	}
	if errorMessage, ok := data["error_message"].(string); ok {
		m.ErrorMessage = errorMessage
	}
	if completedAtStr, ok := data["completed_at"].(string); ok && completedAtStr != "" {
		completedAt, err := time.Parse(time.RFC3339, completedAtStr)
		if err == nil {
			m.CompletedAt = &completedAt
		}
	}

	return m, nil
}
