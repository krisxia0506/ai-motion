package character

import "context"

type CharacterRepository interface {
	Save(ctx context.Context, character *Character) error
	FindByID(ctx context.Context, id CharacterID) (*Character, error)
	FindByNovelID(ctx context.Context, novelID string) ([]*Character, error)
	FindByName(ctx context.Context, novelID, name string) (*Character, error)
	Delete(ctx context.Context, id CharacterID) error
	DeleteByNovelID(ctx context.Context, novelID string) error
}
