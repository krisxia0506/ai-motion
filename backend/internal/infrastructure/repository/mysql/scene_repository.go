package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xiajiayi/ai-motion/internal/domain/scene"
)

type MySQLSceneRepository struct {
	db *sql.DB
}

func NewMySQLSceneRepository(db *sql.DB) scene.SceneRepository {
	return &MySQLSceneRepository{db: db}
}

type sceneRow struct {
	ID          string
	ChapterID   string
	SceneNumber int
	Description string
	Dialogue    string
	Location    sql.NullString
	TimeOfDay   sql.NullString
	Characters  []byte
	Prompt      sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *MySQLSceneRepository) Save(ctx context.Context, s *scene.Scene) error {
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

	query := `
		INSERT INTO scenes (
			id, chapter_id, scene_number, description, dialogue,
			location, time_of_day, characters, prompt, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			description = VALUES(description),
			dialogue = VALUES(dialogue),
			location = VALUES(location),
			time_of_day = VALUES(time_of_day),
			characters = VALUES(characters),
			prompt = VALUES(prompt),
			updated_at = VALUES(updated_at)
	`

	_, err = r.db.ExecContext(ctx, query,
		s.ID, s.ChapterID, s.SceneNumber,
		descriptionJSON, dialoguesJSON,
		s.Location, s.TimeOfDay,
		charactersJSON, s.ImagePrompt,
		s.CreatedAt, s.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save scene: %w", err)
	}

	return nil
}

func (r *MySQLSceneRepository) FindByID(ctx context.Context, id scene.SceneID) (*scene.Scene, error) {
	query := `
		SELECT id, chapter_id, scene_number, description, dialogue,
		       location, time_of_day, characters, prompt, created_at, updated_at
		FROM scenes
		WHERE id = ?
	`

	var row sceneRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&row.ID, &row.ChapterID, &row.SceneNumber,
		&row.Description, &row.Dialogue,
		&row.Location, &row.TimeOfDay,
		&row.Characters, &row.Prompt,
		&row.CreatedAt, &row.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, scene.ErrSceneNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find scene: %w", err)
	}

	return r.scanScene(&row)
}

func (r *MySQLSceneRepository) FindByChapterID(ctx context.Context, chapterID string) ([]*scene.Scene, error) {
	query := `
		SELECT id, chapter_id, scene_number, description, dialogue,
		       location, time_of_day, characters, prompt, created_at, updated_at
		FROM scenes
		WHERE chapter_id = ?
		ORDER BY scene_number
	`

	rows, err := r.db.QueryContext(ctx, query, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenes: %w", err)
	}
	defer rows.Close()

	return r.scanScenes(rows)
}

func (r *MySQLSceneRepository) FindByNovelID(ctx context.Context, novelID string) ([]*scene.Scene, error) {
	query := `
		SELECT s.id, s.chapter_id, s.scene_number, s.description, s.dialogue,
		       s.location, s.time_of_day, s.characters, s.prompt, s.created_at, s.updated_at
		FROM scenes s
		INNER JOIN chapters c ON s.chapter_id = c.id
		WHERE c.novel_id = ?
		ORDER BY c.chapter_number, s.scene_number
	`

	rows, err := r.db.QueryContext(ctx, query, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to query scenes by novel: %w", err)
	}
	defer rows.Close()

	return r.scanScenes(rows)
}

func (r *MySQLSceneRepository) Delete(ctx context.Context, id scene.SceneID) error {
	query := `DELETE FROM scenes WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete scene: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return scene.ErrSceneNotFound
	}

	return nil
}

func (r *MySQLSceneRepository) DeleteByChapterID(ctx context.Context, chapterID string) error {
	query := `DELETE FROM scenes WHERE chapter_id = ?`

	_, err := r.db.ExecContext(ctx, query, chapterID)
	if err != nil {
		return fmt.Errorf("failed to delete scenes by chapter ID: %w", err)
	}

	return nil
}

func (r *MySQLSceneRepository) BatchSave(ctx context.Context, scenes []*scene.Scene) error {
	if len(scenes) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO scenes (
			id, chapter_id, scene_number, description, dialogue,
			location, time_of_day, characters, prompt, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			description = VALUES(description),
			dialogue = VALUES(dialogue),
			location = VALUES(location),
			time_of_day = VALUES(time_of_day),
			characters = VALUES(characters),
			prompt = VALUES(prompt),
			updated_at = VALUES(updated_at)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

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

		_, err = stmt.ExecContext(ctx,
			s.ID, s.ChapterID, s.SceneNumber,
			descriptionJSON, dialoguesJSON,
			s.Location, s.TimeOfDay,
			charactersJSON, s.ImagePrompt,
			s.CreatedAt, s.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to save scene %s: %w", s.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *MySQLSceneRepository) scanScenes(rows *sql.Rows) ([]*scene.Scene, error) {
	var scenes []*scene.Scene

	for rows.Next() {
		var row sceneRow
		err := rows.Scan(
			&row.ID, &row.ChapterID, &row.SceneNumber,
			&row.Description, &row.Dialogue,
			&row.Location, &row.TimeOfDay,
			&row.Characters, &row.Prompt,
			&row.CreatedAt, &row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan scene: %w", err)
		}

		s, err := r.scanScene(&row)
		if err != nil {
			return nil, err
		}

		scenes = append(scenes, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scenes: %w", err)
	}

	return scenes, nil
}

func (r *MySQLSceneRepository) scanScene(row *sceneRow) (*scene.Scene, error) {
	var description scene.Description
	if err := json.Unmarshal([]byte(row.Description), &description); err != nil {
		return nil, fmt.Errorf("failed to unmarshal description: %w", err)
	}

	var dialogues []scene.Dialogue
	if err := json.Unmarshal([]byte(row.Dialogue), &dialogues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dialogues: %w", err)
	}

	var characters []string
	if err := json.Unmarshal(row.Characters, &characters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal characters: %w", err)
	}

	s := &scene.Scene{
		ID:           scene.SceneID(row.ID),
		ChapterID:    row.ChapterID,
		SceneNumber:  row.SceneNumber,
		Description:  description,
		Dialogues:    dialogues,
		CharacterIDs: characters,
		Status:       scene.SceneStatusPending,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	if row.Location.Valid {
		s.Location = row.Location.String
	}
	if row.TimeOfDay.Valid {
		s.TimeOfDay = row.TimeOfDay.String
	}
	if row.Prompt.Valid {
		s.ImagePrompt = row.Prompt.String
	}

	return s, nil
}
