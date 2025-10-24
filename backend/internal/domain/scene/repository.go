package scene

import "context"

type SceneRepository interface {
	Save(ctx context.Context, scene *Scene) error
	FindByID(ctx context.Context, id SceneID) (*Scene, error)
	FindByChapterID(ctx context.Context, chapterID string) ([]*Scene, error)
	FindByNovelID(ctx context.Context, novelID string) ([]*Scene, error)
	Delete(ctx context.Context, id SceneID) error
	DeleteByChapterID(ctx context.Context, chapterID string) error
	BatchSave(ctx context.Context, scenes []*Scene) error
}
