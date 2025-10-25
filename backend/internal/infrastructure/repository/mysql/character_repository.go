package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/domain/character"
)

type MySQLCharacterRepository struct {
	db *sql.DB
}

func NewMySQLCharacterRepository(db *sql.DB) character.CharacterRepository {
	return &MySQLCharacterRepository{db: db}
}

func (r *MySQLCharacterRepository) Save(ctx context.Context, char *character.Character) error {
	appearanceJSON, err := json.Marshal(char.Appearance)
	if err != nil {
		return fmt.Errorf("failed to marshal appearance: %w", err)
	}

	personalityJSON, err := json.Marshal(char.Personality)
	if err != nil {
		return fmt.Errorf("failed to marshal personality: %w", err)
	}

	query := `
		INSERT INTO characters (
			id, novel_id, name, role, appearance, personality, 
			description, reference_image_url, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			role = VALUES(role),
			appearance = VALUES(appearance),
			personality = VALUES(personality),
			description = VALUES(description),
			reference_image_url = VALUES(reference_image_url),
			updated_at = VALUES(updated_at)
	`

	_, err = r.db.ExecContext(ctx, query,
		char.ID, char.NovelID, char.Name, char.Role,
		appearanceJSON, personalityJSON,
		char.Description, char.ReferenceImageURL,
		char.CreatedAt, char.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save character: %w", err)
	}

	return nil
}

func (r *MySQLCharacterRepository) FindByID(ctx context.Context, id character.CharacterID) (*character.Character, error) {
	query := `
		SELECT id, novel_id, name, role, appearance, personality, 
		       description, reference_image_url, created_at, updated_at
		FROM characters
		WHERE id = ?
	`

	var char character.Character
	var appearanceJSON, personalityJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&char.ID, &char.NovelID, &char.Name, &char.Role,
		&appearanceJSON, &personalityJSON,
		&char.Description, &char.ReferenceImageURL,
		&char.CreatedAt, &char.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, character.ErrCharacterNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	if err := json.Unmarshal(appearanceJSON, &char.Appearance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal appearance: %w", err)
	}

	if err := json.Unmarshal(personalityJSON, &char.Personality); err != nil {
		return nil, fmt.Errorf("failed to unmarshal personality: %w", err)
	}

	return &char, nil
}

func (r *MySQLCharacterRepository) FindByNovelID(ctx context.Context, novelID string) ([]*character.Character, error) {
	query := `
		SELECT id, novel_id, name, role, appearance, personality, 
		       description, reference_image_url, created_at, updated_at
		FROM characters
		WHERE novel_id = ?
		ORDER BY role, name
	`

	rows, err := r.db.QueryContext(ctx, query, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %w", err)
	}
	defer rows.Close()

	var characters []*character.Character

	for rows.Next() {
		var char character.Character
		var appearanceJSON, personalityJSON []byte

		err := rows.Scan(
			&char.ID, &char.NovelID, &char.Name, &char.Role,
			&appearanceJSON, &personalityJSON,
			&char.Description, &char.ReferenceImageURL,
			&char.CreatedAt, &char.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}

		if err := json.Unmarshal(appearanceJSON, &char.Appearance); err != nil {
			return nil, fmt.Errorf("failed to unmarshal appearance: %w", err)
		}

		if err := json.Unmarshal(personalityJSON, &char.Personality); err != nil {
			return nil, fmt.Errorf("failed to unmarshal personality: %w", err)
		}

		characters = append(characters, &char)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating characters: %w", err)
	}

	return characters, nil
}

func (r *MySQLCharacterRepository) FindByName(ctx context.Context, novelID, name string) (*character.Character, error) {
	query := `
		SELECT id, novel_id, name, role, appearance, personality, 
		       description, reference_image_url, created_at, updated_at
		FROM characters
		WHERE novel_id = ? AND name = ?
	`

	var char character.Character
	var appearanceJSON, personalityJSON []byte

	err := r.db.QueryRowContext(ctx, query, novelID, name).Scan(
		&char.ID, &char.NovelID, &char.Name, &char.Role,
		&appearanceJSON, &personalityJSON,
		&char.Description, &char.ReferenceImageURL,
		&char.CreatedAt, &char.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, character.ErrCharacterNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	if err := json.Unmarshal(appearanceJSON, &char.Appearance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal appearance: %w", err)
	}

	if err := json.Unmarshal(personalityJSON, &char.Personality); err != nil {
		return nil, fmt.Errorf("failed to unmarshal personality: %w", err)
	}

	return &char, nil
}

func (r *MySQLCharacterRepository) Delete(ctx context.Context, id character.CharacterID) error {
	query := `DELETE FROM characters WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return character.ErrCharacterNotFound
	}

	return nil
}

func (r *MySQLCharacterRepository) DeleteByNovelID(ctx context.Context, novelID string) error {
	query := `DELETE FROM characters WHERE novel_id = ?`

	_, err := r.db.ExecContext(ctx, query, novelID)
	if err != nil {
		return fmt.Errorf("failed to delete characters by novel ID: %w", err)
	}

	return nil
}
