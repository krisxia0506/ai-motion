package service

import (
	"context"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type NovelService struct {
	novelRepo     novel.NovelRepository
	chapterRepo   novel.ChapterRepository
	parserService *novel.ParserService
}

func NewNovelService(
	novelRepo novel.NovelRepository,
	chapterRepo novel.ChapterRepository,
	parserService *novel.ParserService,
) *NovelService {
	return &NovelService{
		novelRepo:     novelRepo,
		chapterRepo:   chapterRepo,
		parserService: parserService,
	}
}

func (s *NovelService) UploadAndParse(ctx context.Context, req *dto.UploadNovelRequest) (*dto.NovelResponse, error) {
	n, err := novel.NewNovel(req.Title, req.Author, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create novel: %w", err)
	}

	if err := s.parserService.Parse(n); err != nil {
		return nil, fmt.Errorf("failed to parse novel: %w", err)
	}

	if err := s.novelRepo.Save(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to save novel: %w", err)
	}

	if len(n.Chapters) > 0 {
		if err := s.chapterRepo.SaveBatch(ctx, n.Chapters); err != nil {
			return nil, fmt.Errorf("failed to save chapters: %w", err)
		}
	}

	return s.toNovelResponse(n), nil
}

func (s *NovelService) GetNovel(ctx context.Context, id string) (*dto.NovelResponse, error) {
	n, err := s.novelRepo.FindByID(ctx, novel.NovelID(id))
	if err != nil {
		return nil, err
	}

	return s.toNovelResponse(n), nil
}

func (s *NovelService) ListNovels(ctx context.Context, offset, limit int) (*dto.NovelListResponse, error) {
	novels, err := s.novelRepo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list novels: %w", err)
	}

	total, err := s.novelRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count novels: %w", err)
	}

	responses := make([]*dto.NovelResponse, len(novels))
	for i, n := range novels {
		responses[i] = s.toNovelResponse(n)
	}

	return &dto.NovelListResponse{
		Novels: responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (s *NovelService) GetChapters(ctx context.Context, novelID string) ([]*dto.ChapterResponse, error) {
	chapters, err := s.chapterRepo.FindByNovelID(ctx, novel.NovelID(novelID))
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	responses := make([]*dto.ChapterResponse, len(chapters))
	for i, chapter := range chapters {
		responses[i] = &dto.ChapterResponse{
			ID:            chapter.ID,
			ChapterNumber: chapter.ChapterNumber,
			Title:         chapter.Title,
			WordCount:     chapter.WordCount,
			CreatedAt:     chapter.CreatedAt,
		}
	}

	return responses, nil
}

func (s *NovelService) DeleteNovel(ctx context.Context, id string) error {
	if err := s.chapterRepo.DeleteByNovelID(ctx, novel.NovelID(id)); err != nil {
		return fmt.Errorf("failed to delete chapters: %w", err)
	}

	if err := s.novelRepo.Delete(ctx, novel.NovelID(id)); err != nil {
		return fmt.Errorf("failed to delete novel: %w", err)
	}

	return nil
}

func (s *NovelService) toNovelResponse(n *novel.Novel) *dto.NovelResponse {
	return &dto.NovelResponse{
		ID:           string(n.ID),
		Title:        n.Title,
		Author:       n.Author,
		Status:       string(n.Status),
		WordCount:    n.WordCount,
		ChapterCount: n.ChapterCount,
		CreatedAt:    n.CreatedAt,
		UpdatedAt:    n.UpdatedAt,
	}
}
