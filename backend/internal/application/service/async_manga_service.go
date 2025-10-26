package service

import (
	"context"
	"fmt"
	"log"

	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
	"github.com/xiajiayi/ai-motion/internal/domain/task"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/gemini"
)

type AsyncMangaService struct {
	taskRepo         task.TaskRepository
	novelRepo        novel.NovelRepository
	chapterRepo      novel.ChapterRepository
	characterRepo    character.CharacterRepository
	sceneRepo        scene.SceneRepository
	parserService    *novel.ParserService
	extractorService *character.CharacterExtractorService
	dividerService   *scene.SceneDividerService
	geminiClient     *gemini.Client
}

func NewAsyncMangaService(
	taskRepo task.TaskRepository,
	novelRepo novel.NovelRepository,
	chapterRepo novel.ChapterRepository,
	characterRepo character.CharacterRepository,
	sceneRepo scene.SceneRepository,
	parserService *novel.ParserService,
	extractorService *character.CharacterExtractorService,
	dividerService *scene.SceneDividerService,
	geminiClient *gemini.Client,
) *AsyncMangaService {
	return &AsyncMangaService{
		taskRepo:         taskRepo,
		novelRepo:        novelRepo,
		chapterRepo:      chapterRepo,
		characterRepo:    characterRepo,
		sceneRepo:        sceneRepo,
		parserService:    parserService,
		extractorService: extractorService,
		dividerService:   dividerService,
		geminiClient:     geminiClient,
	}
}

func (s *AsyncMangaService) StartMangaGeneration(ctx context.Context, title, author, content string) (*task.Task, error) {
	input := map[string]interface{}{
		"title":   title,
		"author":  author,
		"content": content,
	}

	t := task.NewTask(task.TaskTypeMangaGeneration, input)

	if err := s.taskRepo.Save(ctx, t); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	go s.processMangaGeneration(t)

	return t, nil
}

func (s *AsyncMangaService) GetTaskStatus(ctx context.Context, taskID string) (*task.Task, error) {
	return s.taskRepo.FindByID(ctx, taskID)
}

func (s *AsyncMangaService) processMangaGeneration(t *task.Task) {
	ctx := context.Background()

	t.Start()
	if err := s.taskRepo.Update(ctx, t); err != nil {
		log.Printf("Failed to update task status: %v", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in manga generation: %v", r)
			t.Fail(fmt.Errorf("panic: %v", r))
			_ = s.taskRepo.Update(ctx, t)
		}
	}()

	title := t.Input["title"].(string)
	author := t.Input["author"].(string)
	content := t.Input["content"].(string)

	t.UpdateProgress(10, "解析小说中...")
	_ = s.taskRepo.Update(ctx, t)

	n, err := novel.NewNovel(title, author, content)
	if err != nil {
		t.Fail(fmt.Errorf("failed to create novel: %w", err))
		_ = s.taskRepo.Update(ctx, t)
		return
	}

	if err := s.parserService.Parse(n); err != nil {
		t.Fail(fmt.Errorf("failed to parse novel: %w", err))
		_ = s.taskRepo.Update(ctx, t)
		return
	}

	if err := s.novelRepo.Save(ctx, n); err != nil {
		t.Fail(fmt.Errorf("failed to save novel: %w", err))
		_ = s.taskRepo.Update(ctx, t)
		return
	}

	for i := range n.Chapters {
		if err := s.chapterRepo.Save(ctx, &n.Chapters[i]); err != nil {
			t.Fail(fmt.Errorf("failed to save chapter: %w", err))
			_ = s.taskRepo.Update(ctx, t)
			return
		}
	}

	t.UpdateProgress(20, "提取角色中...")
	_ = s.taskRepo.Update(ctx, t)

	characters, err := s.extractorService.ExtractFromNovel(ctx, string(n.ID), n.Content)
	if err != nil {
		t.Fail(fmt.Errorf("failed to extract characters: %w", err))
		_ = s.taskRepo.Update(ctx, t)
		return
	}

	t.UpdateProgress(30, fmt.Sprintf("为 %d 个角色生成参考图...", len(characters)))
	_ = s.taskRepo.Update(ctx, t)

	for i, char := range characters {
		if s.geminiClient != nil {
			appearanceStr := fmt.Sprintf("%s, %s, %s", char.Appearance.PhysicalTraits, char.Appearance.ClothingStyle, char.Appearance.DistinctFeatures)
			prompt := fmt.Sprintf("anime style portrait of %s: %s", char.Name, appearanceStr)

			imageURL, err := s.geminiClient.TextToImage(ctx, gemini.TextToImageRequest{
				Prompt: prompt,
				Style:  "anime",
				Width:  512,
				Height: 512,
			})

			if err != nil {
				log.Printf("Failed to generate reference image for character %s: %v", char.Name, err)
			} else {
				char.ReferenceImageURL = imageURL
				if err := s.characterRepo.Save(ctx, char); err != nil {
					log.Printf("Failed to update character: %v", err)
				}
			}
		}

		progress := 30 + (i+1)*20/len(characters)
		t.UpdateProgress(progress, fmt.Sprintf("已生成 %d/%d 个角色参考图", i+1, len(characters)))
		_ = s.taskRepo.Update(ctx, t)
	}

	t.UpdateProgress(50, "划分场景中...")
	_ = s.taskRepo.Update(ctx, t)

	var allScenes []*scene.Scene
	for _, ch := range n.Chapters {
		sceneChapter := scene.Chapter{
			ID:      ch.ID,
			NovelID: string(ch.NovelID),
			Content: ch.Content,
		}
		scenes, err := s.dividerService.DivideChapterIntoScenes(ctx, sceneChapter)
		if err != nil {
			log.Printf("Failed to divide chapter %s: %v", ch.ID, err)
			continue
		}
		allScenes = append(allScenes, scenes...)
	}

	t.UpdateProgress(60, fmt.Sprintf("为 %d 个场景生成图片...", len(allScenes)))
	_ = s.taskRepo.Update(ctx, t)

	for i, sc := range allScenes {
		if s.geminiClient != nil {
			var refImageURL string
			if len(sc.CharacterIDs) > 0 {
				char, err := s.characterRepo.FindByID(ctx, character.CharacterID(sc.CharacterIDs[0]))
				if err == nil && char.ReferenceImageURL != "" {
					refImageURL = char.ReferenceImageURL
				}
			}

			prompt := fmt.Sprintf("anime style scene: %s, location: %s", sc.Description, sc.Location)

			var imageURL string
			var err error

			if refImageURL != "" {
				imageURL, err = s.geminiClient.ImageToImage(ctx, gemini.ImageToImageRequest{
					ReferenceImage: refImageURL,
					Prompt:         prompt,
					Width:          1024,
					Height:         768,
				})
			} else {
				imageURL, err = s.geminiClient.TextToImage(ctx, gemini.TextToImageRequest{
					Prompt: prompt,
					Style:  "anime",
					Width:  1024,
					Height: 768,
				})
			}

			if err != nil {
				log.Printf("Failed to generate scene image: %v", err)
			} else {
				log.Printf("Generated scene image: %s", imageURL)
			}
		}

		progress := 60 + (i+1)*35/len(allScenes)
		t.UpdateProgress(progress, fmt.Sprintf("已生成 %d/%d 个场景图片", i+1, len(allScenes)))
		_ = s.taskRepo.Update(ctx, t)
	}

	result := map[string]interface{}{
		"novel_id":        n.ID,
		"title":           n.Title,
		"character_count": len(characters),
		"scene_count":     len(allScenes),
		"status":          "completed",
		"message":         fmt.Sprintf("Successfully generated manga with %d characters and %d scenes", len(characters), len(allScenes)),
	}

	t.Complete(result)
	_ = s.taskRepo.Update(ctx, t)
}
