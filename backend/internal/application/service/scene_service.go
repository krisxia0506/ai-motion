package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
)

type SceneService struct {
	sceneRepo          scene.SceneRepository
	chapterRepo        novel.ChapterRepository
	characterRepo      character.CharacterRepository
	dividerService     *scene.SceneDividerService
	promptGeneratorSvc *scene.PromptGeneratorService
}

func NewSceneService(
	sceneRepo scene.SceneRepository,
	chapterRepo novel.ChapterRepository,
	characterRepo character.CharacterRepository,
	dividerService *scene.SceneDividerService,
	promptGeneratorSvc *scene.PromptGeneratorService,
) *SceneService {
	return &SceneService{
		sceneRepo:          sceneRepo,
		chapterRepo:        chapterRepo,
		characterRepo:      characterRepo,
		dividerService:     dividerService,
		promptGeneratorSvc: promptGeneratorSvc,
	}
}

func (s *SceneService) DivideChapter(ctx context.Context, chapterID string) ([]*dto.SceneResponse, error) {
	chapter, err := s.chapterRepo.FindByID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find chapter: %w", err)
	}

	domainChapter := scene.Chapter{
		ID:      chapter.ID,
		NovelID: string(chapter.NovelID),
		Content: chapter.Content,
	}

	scenes, err := s.dividerService.DivideChapterIntoScenes(ctx, domainChapter)
	if err != nil {
		return nil, fmt.Errorf("failed to divide chapter: %w", err)
	}

	responses := make([]*dto.SceneResponse, len(scenes))
	for i, sc := range scenes {
		responses[i] = s.toSceneResponse(sc)
	}

	return responses, nil
}

func (s *SceneService) GetScene(ctx context.Context, id string) (*dto.SceneResponse, error) {
	sc, err := s.sceneRepo.FindByID(ctx, scene.SceneID(id))
	if err != nil {
		return nil, err
	}

	return s.toSceneResponse(sc), nil
}

func (s *SceneService) GetScenesByChapterID(ctx context.Context, chapterID string) ([]*dto.SceneResponse, error) {
	scenes, err := s.sceneRepo.FindByChapterID(ctx, chapterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scenes: %w", err)
	}

	responses := make([]*dto.SceneResponse, len(scenes))
	for i, sc := range scenes {
		responses[i] = s.toSceneResponse(sc)
	}

	return responses, nil
}

func (s *SceneService) GetScenesByNovelID(ctx context.Context, novelID string) ([]*dto.SceneResponse, error) {
	scenes, err := s.sceneRepo.FindByNovelID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scenes: %w", err)
	}

	responses := make([]*dto.SceneResponse, len(scenes))
	for i, sc := range scenes {
		responses[i] = s.toSceneResponse(sc)
	}

	return responses, nil
}

func (s *SceneService) GeneratePrompt(ctx context.Context, req *dto.GenerateScenePromptRequest) (*dto.GenerateScenePromptResponse, error) {
	sc, err := s.sceneRepo.FindByID(ctx, scene.SceneID(req.SceneID))
	if err != nil {
		return nil, fmt.Errorf("failed to find scene: %w", err)
	}

	var characters []scene.Character
	for _, charID := range req.CharacterIDs {
		char, err := s.characterRepo.FindByID(ctx, character.CharacterID(charID))
		if err != nil {
			continue
		}

		characters = append(characters, scene.Character{
			ID:   string(char.ID),
			Name: char.Name,
			Appearance: scene.CharacterAppearance{
				PhysicalTraits:   char.Appearance.PhysicalTraits,
				ClothingStyle:    char.Appearance.ClothingStyle,
				DistinctFeatures: char.Appearance.DistinctFeatures,
				Age:              char.Appearance.Age,
				Height:           char.Appearance.Height,
			},
		})
	}

	options := s.buildPromptOptions(req)

	imagePrompt, err := s.promptGeneratorSvc.GenerateImagePrompt(ctx, sc, characters, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image prompt: %w", err)
	}

	return &dto.GenerateScenePromptResponse{
		SceneID:     req.SceneID,
		ImagePrompt: imagePrompt,
		GeneratedAt: time.Now(),
	}, nil
}

func (s *SceneService) GenerateBatchPrompts(ctx context.Context, req *dto.GeneratePromptsRequest) error {
	scenes := make([]*scene.Scene, 0, len(req.SceneIDs))

	for _, sceneID := range req.SceneIDs {
		sc, err := s.sceneRepo.FindByID(ctx, scene.SceneID(sceneID))
		if err != nil {
			return fmt.Errorf("failed to find scene %s: %w", sceneID, err)
		}
		scenes = append(scenes, sc)
	}

	charactersMap := make(map[string][]scene.Character)

	for _, sc := range scenes {
		var characters []scene.Character
		for _, charID := range sc.CharacterIDs {
			char, err := s.characterRepo.FindByID(ctx, character.CharacterID(charID))
			if err != nil {
				continue
			}

			characters = append(characters, scene.Character{
				ID:   string(char.ID),
				Name: char.Name,
				Appearance: scene.CharacterAppearance{
					PhysicalTraits:   char.Appearance.PhysicalTraits,
					ClothingStyle:    char.Appearance.ClothingStyle,
					DistinctFeatures: char.Appearance.DistinctFeatures,
					Age:              char.Appearance.Age,
					Height:           char.Appearance.Height,
				},
			})
		}
		charactersMap[string(sc.ID)] = characters
	}

	options := s.buildPromptOptionsFromBatch(req)

	if err := s.promptGeneratorSvc.GenerateBatchPrompts(ctx, scenes, charactersMap, options); err != nil {
		return fmt.Errorf("failed to generate batch prompts: %w", err)
	}

	return nil
}

func (s *SceneService) DeleteScene(ctx context.Context, id string) error {
	if err := s.sceneRepo.Delete(ctx, scene.SceneID(id)); err != nil {
		return fmt.Errorf("failed to delete scene: %w", err)
	}

	return nil
}

func (s *SceneService) buildPromptOptions(req *dto.GenerateScenePromptRequest) scene.PromptOptions {
	options := scene.DefaultPromptOptions()

	if req.Style != "" {
		options.Style = scene.PromptStyle(req.Style)
	}

	if req.Quality != "" {
		options.Quality = req.Quality
	}

	if req.AspectRatio != "" {
		options.AspectRatio = req.AspectRatio
	}

	return options
}

func (s *SceneService) buildPromptOptionsFromBatch(req *dto.GeneratePromptsRequest) scene.PromptOptions {
	options := scene.DefaultPromptOptions()

	if req.Style != "" {
		options.Style = scene.PromptStyle(req.Style)
	}

	if req.Quality != "" {
		options.Quality = req.Quality
	}

	if req.AspectRatio != "" {
		options.AspectRatio = req.AspectRatio
	}

	return options
}

func (s *SceneService) toSceneResponse(sc *scene.Scene) *dto.SceneResponse {
	dialogues := make([]dto.DialogueResponse, len(sc.Dialogues))
	for i, d := range sc.Dialogues {
		dialogues[i] = dto.DialogueResponse{
			Speaker: d.Speaker,
			Content: d.Content,
			Emotion: d.Emotion,
		}
	}

	return &dto.SceneResponse{
		ID:          string(sc.ID),
		ChapterID:   sc.ChapterID,
		NovelID:     sc.NovelID,
		SceneNumber: sc.SceneNumber,
		Location:    sc.Location,
		TimeOfDay:   sc.TimeOfDay,
		Description: dto.DescriptionResponse{
			Setting:    sc.Description.Setting,
			Action:     sc.Description.Action,
			Atmosphere: sc.Description.Atmosphere,
			FullText:   sc.Description.FullText,
		},
		Dialogues:    dialogues,
		CharacterIDs: sc.CharacterIDs,
		ImagePrompt:  sc.ImagePrompt,
		VideoPrompt:  sc.VideoPrompt,
		Status:       string(sc.Status),
		CreatedAt:    sc.CreatedAt,
		UpdatedAt:    sc.UpdatedAt,
	}
}
