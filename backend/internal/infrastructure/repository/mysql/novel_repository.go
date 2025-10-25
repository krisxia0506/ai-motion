package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type NovelRepository struct {
	db *sql.DB
}

func NewNovelRepository(db *sql.DB) novel.NovelRepository {
	return &NovelRepository{db: db}
}

func (r *NovelRepository) Save(ctx context.Context, n *novel.Novel) error {
	query := `
		INSERT INTO novels (id, title, author, content, status, word_count, chapter_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			author = VALUES(author),
			content = VALUES(content),
			status = VALUES(status),
			word_count = VALUES(word_count),
			chapter_count = VALUES(chapter_count),
			updated_at = VALUES(updated_at)
	`

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.Title, n.Author, n.Content, n.Status,
		n.WordCount, n.ChapterCount, n.CreatedAt, n.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save novel: %w", err)
	}

	return nil
}

func (r *NovelRepository) FindByID(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
	query := `
		SELECT id, title, author, content, status, word_count, chapter_count, created_at, updated_at
		FROM novels WHERE id = ?
	`

	n := &novel.Novel{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&n.ID, &n.Title, &n.Author, &n.Content, &n.Status,
		&n.WordCount, &n.ChapterCount, &n.CreatedAt, &n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, novel.ErrNovelNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find novel: %w", err)
	}

	return n, nil
}

func (r *NovelRepository) FindAll(ctx context.Context, offset, limit int) ([]*novel.Novel, error) {
	query := `
		SELECT id, title, author, content, status, word_count, chapter_count, created_at, updated_at
		FROM novels
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query novels: %w", err)
	}
	defer rows.Close()

	var novels []*novel.Novel
	for rows.Next() {
		n := &novel.Novel{}
		err := rows.Scan(
			&n.ID, &n.Title, &n.Author, &n.Content, &n.Status,
			&n.WordCount, &n.ChapterCount, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan novel: %w", err)
		}
		novels = append(novels, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating novels: %w", err)
	}

	return novels, nil
}

func (r *NovelRepository) Delete(ctx context.Context, id novel.NovelID) error {
	query := `DELETE FROM novels WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete novel: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return novel.ErrNovelNotFound
	}

	return nil
}

func (r *NovelRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM novels`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count novels: %w", err)
	}

	return count, nil
}
