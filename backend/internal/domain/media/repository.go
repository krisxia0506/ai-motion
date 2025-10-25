package media

import "context"

type MediaRepository interface {
	Save(ctx context.Context, media *Media) error
	FindByID(ctx context.Context, id MediaID) (*Media, error)
	FindBySceneID(ctx context.Context, sceneID string) ([]*Media, error)
	UpdateStatus(ctx context.Context, id MediaID, status MediaStatus, url string, errorMsg string) error
	Delete(ctx context.Context, id MediaID) error
	FindPendingMedia(ctx context.Context, limit int) ([]*Media, error)
}
