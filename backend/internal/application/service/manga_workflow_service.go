package service

import (
	"context"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/media"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/gemini"
)

type MangaWorkflowService struct {
	novelRepo        novel.NovelRepository
	chapterRepo      novel.ChapterRepository
	characterRepo    character.CharacterRepository
	sceneRepo        scene.SceneRepository
	mediaRepo        media.MediaRepository
	parserService    *novel.ParserService
	extractorService *character.CharacterExtractorService
	dividerService   *scene.SceneDividerService
	geminiClient     *gemini.Client
}

func NewMangaWorkflowService(
	novelRepo novel.NovelRepository,
	chapterRepo novel.ChapterRepository,
	characterRepo character.CharacterRepository,
	sceneRepo scene.SceneRepository,
	mediaRepo media.MediaRepository,
	parserService *novel.ParserService,
	extractorService *character.CharacterExtractorService,
	dividerService *scene.SceneDividerService,
	geminiClient *gemini.Client,
) *MangaWorkflowService {
	return &MangaWorkflowService{
		novelRepo:        novelRepo,
		chapterRepo:      chapterRepo,
		characterRepo:    characterRepo,
		sceneRepo:        sceneRepo,
		mediaRepo:        mediaRepo,
		parserService:    parserService,
		extractorService: extractorService,
		dividerService:   dividerService,
		geminiClient:     geminiClient,
	}
}

func (s *MangaWorkflowService) GenerateMangaFromNovel(ctx context.Context, req *dto.GenerateMangaRequest) (*dto.MangaWorkflowResponse, error) {
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

	characters, err := s.extractorService.ExtractFromNovel(ctx, string(n.ID), n.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract characters: %w", err)
	}

	for _, char := range characters {
		if err := s.generateCharacterReferenceImage(ctx, char); err != nil {
			return nil, fmt.Errorf("failed to generate reference image for character %s: %w", char.Name, err)
		}
	}

	var allScenes []*scene.Scene
	if len(n.Chapters) > 0 {
		for _, chapter := range n.Chapters {
			chapterData := scene.Chapter{
				ID:      chapter.ID,
				NovelID: string(n.ID),
				Content: chapter.Content,
			}

			scenes, err := s.dividerService.DivideChapterIntoScenes(ctx, chapterData)
			if err != nil {
				return nil, fmt.Errorf("failed to divide chapter %s into scenes: %w", chapter.ID, err)
			}

			allScenes = append(allScenes, scenes...)
		}
	} else {
		chapterData := scene.Chapter{
			ID:      string(n.ID),
			NovelID: string(n.ID),
			Content: n.Content,
		}

		scenes, err := s.dividerService.DivideChapterIntoScenes(ctx, chapterData)
		if err != nil {
			return nil, fmt.Errorf("failed to divide novel into scenes: %w", err)
		}

		allScenes = scenes
	}

	for _, scn := range allScenes {
		sceneCharacters := s.matchCharactersToScene(scn, characters)
		scn.SetCharacters(sceneCharacters)

		if err := s.sceneRepo.Save(ctx, scn); err != nil {
			return nil, fmt.Errorf("failed to update scene %s: %w", scn.ID, err)
		}

		if err := s.generateSceneImage(ctx, scn, characters); err != nil {
			return nil, fmt.Errorf("failed to generate image for scene %s: %w", scn.ID, err)
		}
	}

	if err := n.UpdateStatus(novel.NovelStatusCompleted); err != nil {
		return nil, fmt.Errorf("failed to update novel status: %w", err)
	}
	if err := s.novelRepo.Save(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to save novel status: %w", err)
	}

	return &dto.MangaWorkflowResponse{
		NovelID:        string(n.ID),
		Title:          n.Title,
		CharacterCount: len(characters),
		SceneCount:     len(allScenes),
		Status:         "completed",
		Message:        fmt.Sprintf("Successfully generated manga with %d characters and %d scenes", len(characters), len(allScenes)),
	}, nil
}

func (s *MangaWorkflowService) generateCharacterReferenceImage(ctx context.Context, char *character.Character) error {
	prompt := s.buildCharacterPrompt(char)

	req := gemini.TextToImageRequest{
		Prompt: prompt,
		Width:  1024,
		Height: 1024,
		Style:  "anime",
	}

	imageURL, err := s.geminiClient.TextToImage(ctx, req)
	if err != nil {
		return fmt.Errorf("gemini text-to-image failed: %w", err)
	}

	char.SetReferenceImage(imageURL)

	if err := s.characterRepo.Save(ctx, char); err != nil {
		return fmt.Errorf("failed to save character with reference image: %w", err)
	}

	return nil
}

func (s *MangaWorkflowService) buildCharacterPrompt(char *character.Character) string {
	prompt := fmt.Sprintf("anime character portrait: %s", char.Name)

	if char.Appearance.PhysicalTraits != "" {
		prompt += fmt.Sprintf(", %s", char.Appearance.PhysicalTraits)
	}

	if char.Appearance.ClothingStyle != "" {
		prompt += fmt.Sprintf(", wearing %s", char.Appearance.ClothingStyle)
	}

	if char.Appearance.Age != "" {
		prompt += fmt.Sprintf(", %s", char.Appearance.Age)
	}

	if char.Description != "" {
		prompt += fmt.Sprintf(", %s", char.Description)
	}

	prompt += ", high quality anime art style, detailed, clean background"

	return prompt
}

func (s *MangaWorkflowService) matchCharactersToScene(scn *scene.Scene, characters []*character.Character) []string {
	var matchedIDs []string

	for _, char := range characters {
		for _, dialogue := range scn.Dialogues {
			if dialogue.Speaker == char.Name {
				matchedIDs = append(matchedIDs, string(char.ID))
				break
			}
		}
	}

	return matchedIDs
}

func (s *MangaWorkflowService) generateSceneImage(ctx context.Context, scn *scene.Scene, characters []*character.Character) error {
	charMap := make(map[string]*character.Character)
	for _, char := range characters {
		charMap[string(char.ID)] = char
	}

	var referenceImages []string
	var characterDescriptions []string

	for _, charID := range scn.CharacterIDs {
		if char, ok := charMap[charID]; ok {
			if char.ReferenceImageURL != "" {
				referenceImages = append(referenceImages, char.ReferenceImageURL)
			}
			characterDescriptions = append(characterDescriptions, char.Name)
		}
	}

	prompt := s.buildScenePrompt(scn, characterDescriptions)

	mediaEntity := media.NewMedia(string(scn.ID), media.MediaTypeImage)
	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return fmt.Errorf("failed to create media entity: %w", err)
	}

	var imageURL string
	var err error

	if len(referenceImages) > 0 {
		req := gemini.ImageToImageRequest{
			ReferenceImage: referenceImages[0],
			Prompt:         prompt,
			Width:          1024,
			Height:         768,
			Strength:       0.6,
		}
		imageURL, err = s.geminiClient.ImageToImage(ctx, req)
	} else {
		req := gemini.TextToImageRequest{
			Prompt: prompt,
			Width:  1024,
			Height: 768,
			Style:  "anime",
		}
		imageURL, err = s.geminiClient.TextToImage(ctx, req)
	}

	if err != nil {
		mediaEntity.MarkFailed(err.Error())
		s.mediaRepo.Save(ctx, mediaEntity)
		return fmt.Errorf("failed to generate scene image: %w", err)
	}

	metadata := media.NewImageMetadata(1024, 768, "image/jpeg", 0)
	mediaEntity.MarkCompleted(imageURL, metadata)

	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return fmt.Errorf("failed to save media: %w", err)
	}

	return nil
}

func (s *MangaWorkflowService) buildScenePrompt(scn *scene.Scene, characterNames []string) string {
	prompt := "anime manga panel: "

	if scn.Location != "" {
		prompt += fmt.Sprintf("scene in %s", scn.Location)
	}

	if scn.TimeOfDay != "" {
		prompt += fmt.Sprintf(", time: %s", scn.TimeOfDay)
	}

	if len(characterNames) > 0 {
		prompt += fmt.Sprintf(", featuring characters: %s", characterNames[0])
		for i := 1; i < len(characterNames); i++ {
			prompt += fmt.Sprintf(" and %s", characterNames[i])
		}
	}

	if scn.Description.FullText != "" {
		descriptionPreview := scn.Description.FullText
		if len(descriptionPreview) > 200 {
			descriptionPreview = descriptionPreview[:200]
		}
		prompt += fmt.Sprintf(", scene: %s", descriptionPreview)
	}

	prompt += ", anime art style, manga panel composition, high quality"

	return prompt
}
