package novel

import "context"

type NovelRepository interface {
	Save(ctx context.Context, novel *Novel) error
	FindByID(ctx context.Context, id NovelID) (*Novel, error)
	FindAll(ctx context.Context, offset, limit int) ([]*Novel, error)
	Delete(ctx context.Context, id NovelID) error
	Count(ctx context.Context) (int, error)
}

type ChapterRepository interface {
	Save(ctx context.Context, chapter *Chapter) error
	SaveBatch(ctx context.Context, chapters []Chapter) error
	FindByNovelID(ctx context.Context, novelID NovelID) ([]Chapter, error)
	FindByID(ctx context.Context, id string) (*Chapter, error)
	DeleteByNovelID(ctx context.Context, novelID NovelID) error
}
