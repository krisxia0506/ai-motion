package supabase

import (
	"context"
	"encoding/json"
	"fmt"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
)

type SceneRepository struct {
	client *postgrest.Client
}

func NewSceneRepository(client *postgrest.Client) scene.SceneRepository {
	return &SceneRepository{client: client}
}

func (r *SceneRepository) Save(ctx context.Context, s *scene.Scene) error {
	descriptionJSON, err := json.Marshal(s.Description)
	if err != nil {
		return fmt.Errorf("failed to marshal description: %w", err)
	}

	dialoguesJSON, err := json.Marshal(s.Dialogues)
	if err != nil {
		return fmt.Errorf("failed to marshal dialogues: %w", err)
	}

	charactersJSON, err := json.Marshal(s.CharacterIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal characters: %w", err)
	}

	data := map[string]interface{}{
		"id":           string(s.ID),
		"chapter_id":   s.ChapterID,
		"novel_id":     s.NovelID,
		"scene_number": s.SceneNumber,
		"location":     s.Location,
		"time_of_day":  s.TimeOfDay,
		"description":  string(descriptionJSON),
		"dialogues":    string(dialoguesJSON),
		"character_ids": string(charactersJSON),
		"image_prompt": s.ImagePrompt,
		"video_prompt": s.VideoPrompt,
		"status":       string(s.Status),
		"created_at":   s.CreatedAt,
		"updated_at":   s.UpdatedAt,
	}

	_, _, err = r.client.From("aimotion_scene").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save scene: %w", err)
	}

	return nil
}

func (r *SceneRepository) FindByID(ctx context.Context, id scene.SceneID) (*scene.Scene, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_scene").
		Select("*", "", false).
		Eq("id", string(id)).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find scene: %w", err)
	}

	if len(results) == 0 {
		return nil, scene.ErrSceneNotFound
	}

	return r.mapToScene(results[0])
}

func (r *SceneRepository) FindByChapterID(ctx context.Context, chapterID string) ([]*scene.Scene, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_scene").
		Select("*", "", false).
		Eq("chapter_id", chapterID).
		Order("scene_number", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find scenes by chapter ID: %w", err)
	}

	var scenes []*scene.Scene
	for _, result := range results {
		s, err := r.mapToScene(result)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, s)
	}

	return scenes, nil
}

func (r *SceneRepository) FindByNovelID(ctx context.Context, novelID string) ([]*scene.Scene, error) {
	var results []map[string]interface{}

	_, err := r.client.From("aimotion_scene").
		Select("*", "", false).
		Eq("novel_id", novelID).
		Order("scene_number", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find scenes by novel ID: %w", err)
	}

	var scenes []*scene.Scene
	for _, result := range results {
		s, err := r.mapToScene(result)
		if err != nil {
			return nil, err
		}
		scenes = append(scenes, s)
	}

	return scenes, nil
}

func (r *SceneRepository) Delete(ctx context.Context, id scene.SceneID) error {
	_, _, err := r.client.From("aimotion_scene").
		Delete("", "").
		Eq("id", string(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete scene: %w", err)
	}

	return nil
}

func (r *SceneRepository) DeleteByChapterID(ctx context.Context, chapterID string) error {
	_, _, err := r.client.From("aimotion_scene").
		Delete("", "").
		Eq("chapter_id", chapterID).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete scenes by chapter ID: %w", err)
	}

	return nil
}

func (r *SceneRepository) BatchSave(ctx context.Context, scenes []*scene.Scene) error {
	if len(scenes) == 0 {
		return nil
	}

	var data []map[string]interface{}
	for _, s := range scenes {
		descriptionJSON, err := json.Marshal(s.Description)
		if err != nil {
			return fmt.Errorf("failed to marshal description: %w", err)
		}

		dialoguesJSON, err := json.Marshal(s.Dialogues)
		if err != nil {
			return fmt.Errorf("failed to marshal dialogues: %w", err)
		}

		charactersJSON, err := json.Marshal(s.CharacterIDs)
		if err != nil {
			return fmt.Errorf("failed to marshal characters: %w", err)
		}

		data = append(data, map[string]interface{}{
			"id":            string(s.ID),
			"chapter_id":    s.ChapterID,
			"novel_id":      s.NovelID,
			"scene_number":  s.SceneNumber,
			"location":      s.Location,
			"time_of_day":   s.TimeOfDay,
			"description":   string(descriptionJSON),
			"dialogues":     string(dialoguesJSON),
			"character_ids": string(charactersJSON),
			"image_prompt":  s.ImagePrompt,
			"video_prompt":  s.VideoPrompt,
			"status":        string(s.Status),
			"created_at":    s.CreatedAt,
			"updated_at":    s.UpdatedAt,
		})
	}

	_, _, err := r.client.From("aimotion_scene").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to batch save scenes: %w", err)
	}

	return nil
}

func (r *SceneRepository) mapToScene(data map[string]interface{}) (*scene.Scene, error) {
	s := &scene.Scene{}

	if id, ok := data["id"].(string); ok {
		s.ID = scene.SceneID(id)
	}
	if chapterID, ok := data["chapter_id"].(string); ok {
		s.ChapterID = chapterID
	}
	if novelID, ok := data["novel_id"].(string); ok {
		s.NovelID = novelID
	}
	if sceneNumber, ok := data["scene_number"].(float64); ok {
		s.SceneNumber = int(sceneNumber)
	}
	if location, ok := data["location"].(string); ok {
		s.Location = location
	}
	if timeOfDay, ok := data["time_of_day"].(string); ok {
		s.TimeOfDay = timeOfDay
	}
	if imagePrompt, ok := data["image_prompt"].(string); ok {
		s.ImagePrompt = imagePrompt
	}
	if videoPrompt, ok := data["video_prompt"].(string); ok {
		s.VideoPrompt = videoPrompt
	}
	if status, ok := data["status"].(string); ok {
		s.Status = scene.SceneStatus(status)
	}

	if descriptionStr, ok := data["description"].(string); ok && descriptionStr != "" {
		if err := json.Unmarshal([]byte(descriptionStr), &s.Description); err != nil {
			return nil, fmt.Errorf("failed to unmarshal description: %w", err)
		}
	}

	if dialoguesStr, ok := data["dialogues"].(string); ok && dialoguesStr != "" {
		if err := json.Unmarshal([]byte(dialoguesStr), &s.Dialogues); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dialogues: %w", err)
		}
	}

	if charactersStr, ok := data["character_ids"].(string); ok && charactersStr != "" {
		if err := json.Unmarshal([]byte(charactersStr), &s.CharacterIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal character_ids: %w", err)
		}
	}

	return s, nil
}
