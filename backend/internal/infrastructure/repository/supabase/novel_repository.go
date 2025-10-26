package supabase

import (
	"context"
	"fmt"
	"log"

	postgrest "github.com/supabase-community/postgrest-go"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type NovelRepository struct {
	client *postgrest.Client
}

func NewNovelRepository(client *postgrest.Client) novel.NovelRepository {
	return &NovelRepository{client: client}
}

// getClientWithAuth returns a client with JWT token from context if available
func (r *NovelRepository) getClientWithAuth(ctx context.Context) *postgrest.Client {
	// Try to get JWT token from context
	if jwtToken, ok := ctx.Value("jwt_token").(string); ok && jwtToken != "" {
		log.Printf("[NovelRepo] Using user JWT token for authentication (length: %d)", len(jwtToken))
		// SetAuthToken creates a new client with Authorization header
		// The apikey should be preserved from the original client
		authClient := r.client.SetAuthToken(jwtToken)
		log.Printf("[NovelRepo] Auth client created, ClientError: %v", authClient.ClientError)
		return authClient
	}
	// Fallback to service role client
	log.Printf("[NovelRepo] No JWT token in context, using service role client")
	return r.client
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

	client := r.getClientWithAuth(ctx)
	_, _, err := client.From("aimotion_novel").Upsert(data, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to save novel: %w", err)
	}

	return nil
}

func (r *NovelRepository) FindByID(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
	var novels []novel.Novel

	client := r.getClientWithAuth(ctx)
	_, err := client.From("aimotion_novel").
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

	client := r.getClientWithAuth(ctx)
	_, err := client.From("aimotion_novel").
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
	client := r.getClientWithAuth(ctx)
	_, _, err := client.From("aimotion_novel").
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

	client := r.getClientWithAuth(ctx)
	count, err := client.From("aimotion_novel").
		Select("id", "exact", false).
		ExecuteTo(&novels)

	if err != nil {
		return 0, fmt.Errorf("failed to count novels: %w", err)
	}

	return int(count), nil
}
