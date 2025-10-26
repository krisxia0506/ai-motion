package supabase

import (
	"context"
	"fmt"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type NovelRepository struct {
	client *postgrest.Client
}

func NewNovelRepository(client *postgrest.Client) novel.NovelRepository {
	return &NovelRepository{client: client}
}

func (r *NovelRepository) Save(ctx context.Context, n *novel.Novel) error {
	data := map[string]interface{}{
		"id":            string(n.ID),
		"title":         n.Title,
		"author":        n.Author,
		"content":       n.Content,
		"status":        string(n.Status),
		"word_count":    n.WordCount,
		"chapter_count": n.ChapterCount,
		"created_at":    n.CreatedAt,
		"updated_at":    n.UpdatedAt,
	}

	_, _, err := r.client.From("novels").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save novel: %w", err)
	}

	return nil
}

func (r *NovelRepository) FindByID(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
	var novels []novel.Novel

	_, err := r.client.From("novels").
		Select("*", "", false).
		Eq("id", string(id)).
		ExecuteTo(&novels)

	if err != nil {
		return nil, fmt.Errorf("failed to find novel: %w", err)
	}

	if len(novels) == 0 {
		return nil, novel.ErrNovelNotFound
	}

	return &novels[0], nil
}

func (r *NovelRepository) FindAll(ctx context.Context, offset, limit int) ([]*novel.Novel, error) {
	var novels []*novel.Novel

	_, err := r.client.From("novels").
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&novels)

	if err != nil {
		return nil, fmt.Errorf("failed to query novels: %w", err)
	}

	return novels, nil
}

func (r *NovelRepository) Delete(ctx context.Context, id novel.NovelID) error {
	_, _, err := r.client.From("novels").
		Delete("", "").
		Eq("id", string(id)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete novel: %w", err)
	}

	return nil
}

func (r *NovelRepository) Count(ctx context.Context) (int, error) {
	var novels []novel.Novel

	count, err := r.client.From("novels").
		Select("id", "exact", false).
		ExecuteTo(&novels)

	if err != nil {
		return 0, fmt.Errorf("failed to count novels: %w", err)
	}

	return int(count), nil
}
