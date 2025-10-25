package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xiajiayi/ai-motion/internal/domain/media"
)

type MediaRepository struct {
	db *sql.DB
}

func NewMediaRepository(db *sql.DB) media.MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Save(ctx context.Context, m *media.Media) error {
	query := `
		INSERT INTO media (
			id, scene_id, type, status, url, width, height, duration,
			format, file_size, generation_id, error_message, created_at, updated_at, completed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			url = VALUES(url),
			width = VALUES(width),
			height = VALUES(height),
			duration = VALUES(duration),
			format = VALUES(format),
			file_size = VALUES(file_size),
			generation_id = VALUES(generation_id),
			error_message = VALUES(error_message),
			updated_at = VALUES(updated_at),
			completed_at = VALUES(completed_at)
	`

	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.SceneID, m.Type, m.Status, m.URL,
		m.Metadata.Width, m.Metadata.Height, m.Metadata.Duration,
		m.Metadata.Format, m.Metadata.FileSize,
		m.GenerationID, m.ErrorMessage,
		m.CreatedAt, m.UpdatedAt, m.CompletedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save media: %w", err)
	}

	return nil
}

func (r *MediaRepository) FindByID(ctx context.Context, id media.MediaID) (*media.Media, error) {
	query := `
		SELECT id, scene_id, type, status, url, width, height, duration,
			   format, file_size, generation_id, error_message, created_at, updated_at, completed_at
		FROM media
		WHERE id = ?
	`

	var m media.Media
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.SceneID, &m.Type, &m.Status, &m.URL,
		&m.Metadata.Width, &m.Metadata.Height, &m.Metadata.Duration,
		&m.Metadata.Format, &m.Metadata.FileSize,
		&m.GenerationID, &m.ErrorMessage,
		&m.CreatedAt, &m.UpdatedAt, &completedAt,
	)

	if err == sql.ErrNoRows {
		return nil, media.ErrMediaNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	if completedAt.Valid {
		m.CompletedAt = &completedAt.Time
	}

	return &m, nil
}

func (r *MediaRepository) FindBySceneID(ctx context.Context, sceneID string) ([]*media.Media, error) {
	query := `
		SELECT id, scene_id, type, status, url, width, height, duration,
			   format, file_size, generation_id, error_message, created_at, updated_at, completed_at
		FROM media
		WHERE scene_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sceneID)
	if err != nil {
		return nil, fmt.Errorf("failed to query media: %w", err)
	}
	defer rows.Close()

	var mediaList []*media.Media
	for rows.Next() {
		var m media.Media
		var completedAt sql.NullTime

		err := rows.Scan(
			&m.ID, &m.SceneID, &m.Type, &m.Status, &m.URL,
			&m.Metadata.Width, &m.Metadata.Height, &m.Metadata.Duration,
			&m.Metadata.Format, &m.Metadata.FileSize,
			&m.GenerationID, &m.ErrorMessage,
			&m.CreatedAt, &m.UpdatedAt, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}

		if completedAt.Valid {
			m.CompletedAt = &completedAt.Time
		}

		mediaList = append(mediaList, &m)
	}

	return mediaList, nil
}

func (r *MediaRepository) UpdateStatus(ctx context.Context, id media.MediaID, status media.MediaStatus, url string, errorMsg string) error {
	query := `
		UPDATE media
		SET status = ?, url = ?, error_message = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, url, errorMsg, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update media status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return media.ErrMediaNotFound
	}

	return nil
}

func (r *MediaRepository) Delete(ctx context.Context, id media.MediaID) error {
	query := `DELETE FROM media WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return media.ErrMediaNotFound
	}

	return nil
}

func (r *MediaRepository) FindPendingMedia(ctx context.Context, limit int) ([]*media.Media, error) {
	query := `
		SELECT id, scene_id, type, status, url, width, height, duration,
			   format, file_size, generation_id, error_message, created_at, updated_at, completed_at
		FROM media
		WHERE status = ?
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, media.MediaStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending media: %w", err)
	}
	defer rows.Close()

	var mediaList []*media.Media
	for rows.Next() {
		var m media.Media
		var completedAt sql.NullTime

		err := rows.Scan(
			&m.ID, &m.SceneID, &m.Type, &m.Status, &m.URL,
			&m.Metadata.Width, &m.Metadata.Height, &m.Metadata.Duration,
			&m.Metadata.Format, &m.Metadata.FileSize,
			&m.GenerationID, &m.ErrorMessage,
			&m.CreatedAt, &m.UpdatedAt, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media: %w", err)
		}

		if completedAt.Valid {
			m.CompletedAt = &completedAt.Time
		}

		mediaList = append(mediaList, &m)
	}

	return mediaList, nil
}
