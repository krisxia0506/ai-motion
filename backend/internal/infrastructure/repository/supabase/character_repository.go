package supabase

import (
	"context"
	"encoding/json"
	"fmt"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
)

type CharacterRepository struct {
	client *postgrest.Client
}

func NewCharacterRepository(client *postgrest.Client) character.CharacterRepository {
	return &CharacterRepository{client: client}
}

func (r *CharacterRepository) Save(ctx context.Context, char *character.Character) error {
	appearanceJSON, err := json.Marshal(char.Appearance)
	if err != nil {
		return fmt.Errorf("failed to marshal appearance: %w", err)
	}

	personalityJSON, err := json.Marshal(char.Personality)
	if err != nil {
		return fmt.Errorf("failed to marshal personality: %w", err)
	}

	data := map[string]interface{}{
		"id":                  string(char.ID),
		"novel_id":            char.NovelID,
		"name":                char.Name,
		"role":                string(char.Role),
		"appearance":          string(appearanceJSON),
		"personality":         string(personalityJSON),
		"description":         char.Description,
		"reference_image_url": char.ReferenceImageURL,
		"created_at":          char.CreatedAt,
		"updated_at":          char.UpdatedAt,
	}

	_, _, err = r.client.From("characters").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}

	return nil
}

func (r *CharacterRepository) FindByID(ctx context.Context, id character.CharacterID) (*character.Character, error) {
	var results []map[string]interface{}

	_, err := r.client.From("characters").
		Select("*", "", false).
		Eq("id", string(id)).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	if len(results) == 0 {
		return nil, character.ErrCharacterNotFound
	}

	return r.mapToCharacter(results[0])
}

func (r *CharacterRepository) FindByNovelID(ctx context.Context, novelID string) ([]*character.Character, error) {
	var results []map[string]interface{}

	_, err := r.client.From("characters").
		Select("*", "", false).
		Eq("novel_id", novelID).
		Order("role", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %w", err)
	}

	var characters []*character.Character
	for _, result := range results {
		char, err := r.mapToCharacter(result)
		if err != nil {
			return nil, err
		}
		characters = append(characters, char)
	}

	return characters, nil
}

func (r *CharacterRepository) FindByName(ctx context.Context, novelID, name string) (*character.Character, error) {
	var results []map[string]interface{}

	_, err := r.client.From("characters").
		Select("*", "", false).
		Eq("novel_id", novelID).
		Eq("name", name).
		ExecuteTo(&results)

	if err != nil {
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	if len(results) == 0 {
		return nil, character.ErrCharacterNotFound
	}

	return r.mapToCharacter(results[0])
}

func (r *CharacterRepository) Delete(ctx context.Context, id character.CharacterID) error {
	_, _, err := r.client.From("characters").
		Delete("", "").
		Eq("id", string(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	return nil
}

func (r *CharacterRepository) DeleteByNovelID(ctx context.Context, novelID string) error {
	_, _, err := r.client.From("characters").
		Delete("", "").
		Eq("novel_id", novelID).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete characters by novel ID: %w", err)
	}

	return nil
}

func (r *CharacterRepository) mapToCharacter(data map[string]interface{}) (*character.Character, error) {
	char := &character.Character{}

	if id, ok := data["id"].(string); ok {
		char.ID = character.CharacterID(id)
	}
	if novelID, ok := data["novel_id"].(string); ok {
		char.NovelID = novelID
	}
	if name, ok := data["name"].(string); ok {
		char.Name = name
	}
	if role, ok := data["role"].(string); ok {
		char.Role = character.CharacterRole(role)
	}
	if description, ok := data["description"].(string); ok {
		char.Description = description
	}
	if refImageURL, ok := data["reference_image_url"].(string); ok {
		char.ReferenceImageURL = refImageURL
	}

	if appearanceStr, ok := data["appearance"].(string); ok && appearanceStr != "" {
		if err := json.Unmarshal([]byte(appearanceStr), &char.Appearance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal appearance: %w", err)
		}
	}

	if personalityStr, ok := data["personality"].(string); ok && personalityStr != "" {
		if err := json.Unmarshal([]byte(personalityStr), &char.Personality); err != nil {
			return nil, fmt.Errorf("failed to unmarshal personality: %w", err)
		}
	}

	return char, nil
}
