package supabase

import (
	"context"
	"fmt"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type ChapterRepository struct {
	client *postgrest.Client
}

func NewChapterRepository(client *postgrest.Client) novel.ChapterRepository {
	return &ChapterRepository{client: client}
}

func (r *ChapterRepository) Save(ctx context.Context, chapter *novel.Chapter) error {
	data := map[string]interface{}{
		"id":             chapter.ID,
		"novel_id":       string(chapter.NovelID),
		"chapter_number": chapter.ChapterNumber,
		"title":          chapter.Title,
		"content":        chapter.Content,
		"word_count":     chapter.WordCount,
		"created_at":     chapter.CreatedAt,
		"updated_at":     chapter.UpdatedAt,
	}

	_, _, err := r.client.From("aimotion_chapter").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save chapter: %w", err)
	}

	return nil
}

func (r *ChapterRepository) SaveBatch(ctx context.Context, chapters []novel.Chapter) error {
	if len(chapters) == 0 {
		return nil
	}

	var data []map[string]interface{}
	for _, chapter := range chapters {
		data = append(data, map[string]interface{}{
			"id":             chapter.ID,
			"novel_id":       string(chapter.NovelID),
			"chapter_number": chapter.ChapterNumber,
			"title":          chapter.Title,
			"content":        chapter.Content,
			"word_count":     chapter.WordCount,
			"created_at":     chapter.CreatedAt,
			"updated_at":     chapter.UpdatedAt,
		})
	}

	_, _, err := r.client.From("aimotion_chapter").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save batch chapters: %w", err)
	}

	return nil
}

func (r *ChapterRepository) FindByNovelID(ctx context.Context, novelID novel.NovelID) ([]novel.Chapter, error) {
	var chapters []novel.Chapter

	_, err := r.client.From("aimotion_chapter").
		Select("*", "", false).
		Eq("novel_id", string(novelID)).
		Order("chapter_number", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&chapters)

	if err != nil {
		return nil, fmt.Errorf("failed to find chapters by novel ID: %w", err)
	}

	return chapters, nil
}

func (r *ChapterRepository) FindByID(ctx context.Context, id string) (*novel.Chapter, error) {
	var chapters []novel.Chapter

	_, err := r.client.From("aimotion_chapter").
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&chapters)

	if err != nil {
		return nil, fmt.Errorf("failed to find chapter: %w", err)
	}

	if len(chapters) == 0 {
		return nil, fmt.Errorf("chapter not found")
	}

	return &chapters[0], nil
}

func (r *ChapterRepository) DeleteByNovelID(ctx context.Context, novelID novel.NovelID) error {
	_, _, err := r.client.From("aimotion_chapter").
		Delete("", "").
		Eq("novel_id", string(novelID)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete chapters by novel ID: %w", err)
	}

	return nil
}
