package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type ChapterRepository struct {
	db *sql.DB
}

func NewChapterRepository(db *sql.DB) novel.ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) Save(ctx context.Context, chapter *novel.Chapter) error {
	query := `
		INSERT INTO chapters (id, novel_id, chapter_number, title, content, word_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			content = VALUES(content),
			word_count = VALUES(word_count),
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query,
		chapter.ID, chapter.NovelID, chapter.ChapterNumber,
		chapter.Title, chapter.Content, chapter.WordCount,
	)

	if err != nil {
		return fmt.Errorf("failed to save chapter: %w", err)
	}

	return nil
}

func (r *ChapterRepository) SaveBatch(ctx context.Context, chapters []novel.Chapter) error {
	if len(chapters) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO chapters (id, novel_id, chapter_number, title, content, word_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			title = VALUES(title),
			content = VALUES(content),
			word_count = VALUES(word_count),
			updated_at = NOW()
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, chapter := range chapters {
		_, err := stmt.ExecContext(ctx,
			chapter.ID, chapter.NovelID, chapter.ChapterNumber,
			chapter.Title, chapter.Content, chapter.WordCount,
		)
		if err != nil {
			return fmt.Errorf("failed to save chapter %d: %w", chapter.ChapterNumber, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *ChapterRepository) FindByNovelID(ctx context.Context, novelID novel.NovelID) ([]novel.Chapter, error) {
	query := `
		SELECT id, novel_id, chapter_number, title, content, word_count, created_at, updated_at
		FROM chapters
		WHERE novel_id = ?
		ORDER BY chapter_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to query chapters: %w", err)
	}
	defer rows.Close()

	var chapters []novel.Chapter
	for rows.Next() {
		var chapter novel.Chapter
		err := rows.Scan(
			&chapter.ID, &chapter.NovelID, &chapter.ChapterNumber,
			&chapter.Title, &chapter.Content, &chapter.WordCount,
			&chapter.CreatedAt, &chapter.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chapter: %w", err)
		}
		chapters = append(chapters, chapter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chapters: %w", err)
	}

	return chapters, nil
}

func (r *ChapterRepository) FindByID(ctx context.Context, id string) (*novel.Chapter, error) {
	query := `
		SELECT id, novel_id, chapter_number, title, content, word_count, created_at, updated_at
		FROM chapters WHERE id = ?
	`

	var chapter novel.Chapter
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&chapter.ID, &chapter.NovelID, &chapter.ChapterNumber,
		&chapter.Title, &chapter.Content, &chapter.WordCount,
		&chapter.CreatedAt, &chapter.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("chapter not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find chapter: %w", err)
	}

	return &chapter, nil
}

func (r *ChapterRepository) DeleteByNovelID(ctx context.Context, novelID novel.NovelID) error {
	query := `DELETE FROM chapters WHERE novel_id = ?`

	_, err := r.db.ExecContext(ctx, query, novelID)
	if err != nil {
		return fmt.Errorf("failed to delete chapters: %w", err)
	}

	return nil
}
