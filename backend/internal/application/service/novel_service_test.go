package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/application/service"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type MockNovelRepository struct {
	SaveFunc     func(ctx context.Context, n *novel.Novel) error
	FindByIDFunc func(ctx context.Context, id novel.NovelID) (*novel.Novel, error)
	FindAllFunc  func(ctx context.Context, offset, limit int) ([]*novel.Novel, error)
	CountFunc    func(ctx context.Context) (int, error)
	DeleteFunc   func(ctx context.Context, id novel.NovelID) error
}

func (m *MockNovelRepository) Save(ctx context.Context, n *novel.Novel) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, n)
	}
	return nil
}

func (m *MockNovelRepository) FindByID(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockNovelRepository) FindAll(ctx context.Context, offset, limit int) ([]*novel.Novel, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, offset, limit)
	}
	return []*novel.Novel{}, nil
}

func (m *MockNovelRepository) Count(ctx context.Context) (int, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

func (m *MockNovelRepository) Delete(ctx context.Context, id novel.NovelID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type MockChapterRepository struct {
	SaveFunc            func(ctx context.Context, chapter *novel.Chapter) error
	SaveBatchFunc       func(ctx context.Context, chapters []novel.Chapter) error
	FindByNovelIDFunc   func(ctx context.Context, novelID novel.NovelID) ([]novel.Chapter, error)
	DeleteByNovelIDFunc func(ctx context.Context, novelID novel.NovelID) error
	FindByIDFunc        func(ctx context.Context, id string) (*novel.Chapter, error)
}

func (m *MockChapterRepository) Save(ctx context.Context, chapter *novel.Chapter) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, chapter)
	}
	return nil
}

func (m *MockChapterRepository) SaveBatch(ctx context.Context, chapters []novel.Chapter) error {
	if m.SaveBatchFunc != nil {
		return m.SaveBatchFunc(ctx, chapters)
	}
	return nil
}

func (m *MockChapterRepository) FindByNovelID(ctx context.Context, novelID novel.NovelID) ([]novel.Chapter, error) {
	if m.FindByNovelIDFunc != nil {
		return m.FindByNovelIDFunc(ctx, novelID)
	}
	return []novel.Chapter{}, nil
}

func (m *MockChapterRepository) DeleteByNovelID(ctx context.Context, novelID novel.NovelID) error {
	if m.DeleteByNovelIDFunc != nil {
		return m.DeleteByNovelIDFunc(ctx, novelID)
	}
	return nil
}

func (m *MockChapterRepository) FindByID(ctx context.Context, id string) (*novel.Chapter, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func TestNovelService_UploadAndParse(t *testing.T) {
	mockNovelRepo := &MockNovelRepository{}
	mockChapterRepo := &MockChapterRepository{}
	parserService := novel.NewParserService()

	novelService := service.NewNovelService(mockNovelRepo, mockChapterRepo, parserService)

	tests := []struct {
		name    string
		request *dto.UploadNovelRequest
		wantErr bool
	}{
		{
			name: "valid novel",
			request: &dto.UploadNovelRequest{
				Title:   "Test Novel",
				Author:  "Test Author",
				Content: repeatString("这是一个测试小说的内容。", 20),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			request: &dto.UploadNovelRequest{
				Title:   "",
				Author:  "Test Author",
				Content: repeatString("内容", 20),
			},
			wantErr: true,
		},
		{
			name: "content too short",
			request: &dto.UploadNovelRequest{
				Title:   "Test",
				Author:  "Test",
				Content: "Short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := novelService.UploadAndParse(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UploadAndParse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UploadAndParse() unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Errorf("UploadAndParse() returned nil response")
				return
			}

			if resp.Title != tt.request.Title {
				t.Errorf("Title = %v, want %v", resp.Title, tt.request.Title)
			}
		})
	}
}

func TestNovelService_GetNovel(t *testing.T) {
	testNovel := &novel.Novel{
		ID:      "test-id",
		Title:   "Test Novel",
		Author:  "Test Author",
		Content: "Test content that is long enough to pass validation.",
		Status:  novel.NovelStatusParsed,
	}

	mockNovelRepo := &MockNovelRepository{
		FindByIDFunc: func(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
			if id == testNovel.ID {
				return testNovel, nil
			}
			return nil, novel.ErrNovelNotFound
		},
	}

	novelService := service.NewNovelService(mockNovelRepo, &MockChapterRepository{}, novel.NewParserService())

	resp, err := novelService.GetNovel(context.Background(), string(testNovel.ID))

	if err != nil {
		t.Errorf("GetNovel() unexpected error: %v", err)
	}

	if resp.ID != string(testNovel.ID) {
		t.Errorf("ID = %v, want %v", resp.ID, testNovel.ID)
	}

	if resp.Title != testNovel.Title {
		t.Errorf("Title = %v, want %v", resp.Title, testNovel.Title)
	}
}

func TestNovelService_GetNovel_NotFound(t *testing.T) {
	mockNovelRepo := &MockNovelRepository{
		FindByIDFunc: func(ctx context.Context, id novel.NovelID) (*novel.Novel, error) {
			return nil, novel.ErrNovelNotFound
		},
	}

	novelService := service.NewNovelService(mockNovelRepo, &MockChapterRepository{}, novel.NewParserService())

	_, err := novelService.GetNovel(context.Background(), "non-existent")

	if err == nil {
		t.Errorf("GetNovel() expected error, got nil")
	}
}

func TestNovelService_ListNovels(t *testing.T) {
	mockNovels := []*novel.Novel{
		{ID: "1", Title: "Novel 1", Content: "Content 1 with enough length to pass validation checks."},
		{ID: "2", Title: "Novel 2", Content: "Content 2 with enough length to pass validation checks."},
	}

	mockNovelRepo := &MockNovelRepository{
		FindAllFunc: func(ctx context.Context, offset, limit int) ([]*novel.Novel, error) {
			return mockNovels, nil
		},
		CountFunc: func(ctx context.Context) (int, error) {
			return len(mockNovels), nil
		},
	}

	novelService := service.NewNovelService(mockNovelRepo, &MockChapterRepository{}, novel.NewParserService())

	resp, err := novelService.ListNovels(context.Background(), 0, 10)

	if err != nil {
		t.Errorf("ListNovels() unexpected error: %v", err)
	}

	if len(resp.Novels) != 2 {
		t.Errorf("len(Novels) = %v, want 2", len(resp.Novels))
	}

	if resp.Total != 2 {
		t.Errorf("Total = %v, want 2", resp.Total)
	}
}

func TestNovelService_DeleteNovel(t *testing.T) {
	mockNovelRepo := &MockNovelRepository{
		DeleteFunc: func(ctx context.Context, id novel.NovelID) error {
			return nil
		},
	}

	mockChapterRepo := &MockChapterRepository{
		DeleteByNovelIDFunc: func(ctx context.Context, novelID novel.NovelID) error {
			return nil
		},
	}

	novelService := service.NewNovelService(mockNovelRepo, mockChapterRepo, novel.NewParserService())

	err := novelService.DeleteNovel(context.Background(), "test-id")

	if err != nil {
		t.Errorf("DeleteNovel() unexpected error: %v", err)
	}
}

func TestNovelService_DeleteNovel_Error(t *testing.T) {
	mockNovelRepo := &MockNovelRepository{
		DeleteFunc: func(ctx context.Context, id novel.NovelID) error {
			return errors.New("delete error")
		},
	}

	mockChapterRepo := &MockChapterRepository{
		DeleteByNovelIDFunc: func(ctx context.Context, novelID novel.NovelID) error {
			return nil
		},
	}

	novelService := service.NewNovelService(mockNovelRepo, mockChapterRepo, novel.NewParserService())

	err := novelService.DeleteNovel(context.Background(), "test-id")

	if err == nil {
		t.Errorf("DeleteNovel() expected error, got nil")
	}
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
