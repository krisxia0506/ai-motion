package service

import (
	"context"
	"fmt"
	"log"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/media"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
	"github.com/xiajiayi/ai-motion/internal/domain/task"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/gemini"
)

type MangaWorkflowService struct {
	taskRepo         task.Repository
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
	taskRepo task.Repository,
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
		taskRepo:         taskRepo,
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

// CreateTask 创建漫画生成任务
func (s *MangaWorkflowService) CreateTask(ctx context.Context, userID string, req *dto.GenerateMangaRequest) (*task.Task, error) {
	// 1. 创建Novel实体
	author := req.Author
	if author == "" {
		author = "Unknown"
	}
	novelEntity, err := novel.NewNovel(req.Title, author, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create novel: %w", err)
	}

	// 2. 保存Novel到数据库
	if err := s.novelRepo.Save(ctx, novelEntity); err != nil {
		return nil, fmt.Errorf("failed to save novel: %w", err)
	}

	// 3. 创建Task实体
	taskEntity := task.NewTask(userID, string(novelEntity.ID))

	// 4. 保存Task到数据库
	if err := s.taskRepo.Save(ctx, taskEntity); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	return taskEntity, nil
}

// ExecuteTask 异步执行任务（在goroutine中调用）
func (s *MangaWorkflowService) ExecuteTask(ctx context.Context, taskID string) {
	// 1. 加载任务
	taskEntity, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		log.Printf("Failed to load task %s: %v", taskID, err)
		return
	}

	// 2. 加载小说
	novelEntity, err := s.novelRepo.FindByID(ctx, novel.NovelID(taskEntity.NovelID))
	if err != nil {
		log.Printf("Failed to load novel %s: %v", taskEntity.NovelID, err)
		taskEntity.MarkFailed(50001, "加载小说失败")
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤1: 生成漫画图片(直接调用 Gemini 生成10张)
	if err := s.stepGenerateMangaImages(ctx, taskEntity, novelEntity); err != nil {
		taskEntity.MarkFailed(40001, fmt.Sprintf("生成漫画失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤2: 完成
	taskEntity.MarkCompleted()
	s.taskRepo.Save(ctx, taskEntity)
}

// 步骤1: 生成漫画图片（直接调用 Gemini 生成10张漫画）
func (s *MangaWorkflowService) stepGenerateMangaImages(ctx context.Context, t *task.Task, n *novel.Novel) error {
	if t.IsCancelled() {
		return fmt.Errorf("task cancelled")
	}

	const totalImages = 10
	details := t.ProgressDetails

	// 准备小说内容摘要（限制长度）
	contentSummary := n.Content
	if len(contentSummary) > 2000 {
		contentSummary = contentSummary[:2000] + "..."
	}

	for i := 0; i < totalImages; i++ {
		if t.IsCancelled() {
			return fmt.Errorf("task cancelled")
		}

		// 更新进度
		percentage := (i + 1) * 100 / totalImages
		details.ScenesGenerated = i
		t.UpdateProgress(fmt.Sprintf("生成漫画图片 %d/%d", i+1, totalImages), 1, percentage, details)
		s.taskRepo.Save(ctx, t)

		// 构建 Prompt：让 Gemini 根据小说内容生成漫画风格的图片
		prompt := s.buildMangaPanelPrompt(n.Title, contentSummary, i+1, totalImages)

		// 调用 Gemini 文生图接口
		req := gemini.TextToImageRequest{
			Prompt: prompt,
			Width:  1344,
			Height: 768,
		}

		imageURL, err := s.geminiClient.TextToImage(ctx, req)
		if err != nil {
			log.Printf("Failed to generate manga image %d: %v", i+1, err)
			return fmt.Errorf("failed to generate manga image %d: %w", i+1, err)
		}

		// 保存 Media 实体（使用 NewMediaForNovel 关联小说）
		mediaEntity := media.NewMediaForNovel(string(n.ID), media.MediaTypeImage)
		metadata := media.NewImageMetadata(1344, 768, "image/jpeg", 0)
		mediaEntity.MarkCompleted(imageURL, metadata)

		if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
			return fmt.Errorf("failed to save media %d: %w", i+1, err)
		}

		log.Printf("Successfully generated manga image %d/%d for novel %s", i+1, totalImages, n.ID)
	}

	// 完成所有图片生成
	details.ScenesGenerated = totalImages
	t.UpdateProgress("生成漫画图片完成", 1, 100, details)
	s.taskRepo.Save(ctx, t)

	return nil
}

// buildMangaPanelPrompt 构建漫画面板的 Prompt
func (s *MangaWorkflowService) buildMangaPanelPrompt(title, content string, panelNum, totalPanels int) string {
	// 为不同的面板生成不同的 Prompt，确保故事连贯性
	stageDesc := ""
	switch {
	case panelNum <= 2:
		stageDesc = "opening scene, introducing the setting and atmosphere"
	case panelNum <= 4:
		stageDesc = "introducing main characters and their personalities"
	case panelNum <= 7:
		stageDesc = "developing the plot with key story events"
	case panelNum <= 9:
		stageDesc = "building toward the climax with tension and drama"
	default:
		stageDesc = "conclusion or cliffhanger ending"
	}

	prompt := fmt.Sprintf(
		"Create a beautiful anime manga panel (image %d of %d) based on the novel titled '%s'. "+
			"Story context: %s. "+
			"This panel should depict: %s. "+
			"Style: Professional anime manga art with clean lines, dynamic composition, expressive characters, "+
			"cinematic lighting, and detailed backgrounds. "+
			"Use vibrant colors and dramatic angles typical of high-quality manga adaptations. "+
			"Make it visually engaging and emotionally resonant.",
		panelNum, totalPanels, title, content, stageDesc,
	)

	return prompt
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

// GetTaskStatus 获取任务状态
func (s *MangaWorkflowService) GetTaskStatus(ctx context.Context, userID, taskID string) (*dto.TaskStatusResponse, error) {
	// 根据用户ID和任务ID查找任务（权限验证）
	taskEntity, err := s.taskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("task not found or access denied")
	}

	response := &dto.TaskStatusResponse{
		TaskID:      taskEntity.ID,
		Status:      string(taskEntity.Status),
		Progress: dto.TaskProgressResponse{
			CurrentStep:      taskEntity.ProgressStep,
			CurrentStepIndex: taskEntity.ProgressStepIndex,
			TotalSteps:       2,
			Percentage:       taskEntity.ProgressPercentage,
			Details: &dto.TaskProgressDetailsResponse{
				CharactersExtracted: 0,
				CharactersGenerated: 0,
				ScenesDivided:       0,
				ScenesGenerated:     taskEntity.ProgressDetails.ScenesGenerated,
			},
		},
		CreatedAt:   taskEntity.CreatedAt,
		UpdatedAt:   taskEntity.UpdatedAt,
		CompletedAt: taskEntity.CompletedAt,
		FailedAt:    taskEntity.FailedAt,
	}

	// 如果任务失败，添加错误信息
	if taskEntity.Status == task.TaskStatusFailed {
		response.Error = &dto.TaskErrorResponse{
			Code:      taskEntity.ErrorCode,
			Message:   taskEntity.ErrorMessage,
			RetryAble: taskEntity.IsRetryable(),
		}
	}

	// 如果任务完成，添加结果信息
	if taskEntity.Status == task.TaskStatusCompleted {
		result, err := s.buildTaskResult(ctx, taskEntity)
		if err != nil {
			log.Printf("Failed to build task result: %v", err)
		} else {
			response.Result = result
		}
	}

	return response, nil
}

// buildTaskResult 构建任务结果
func (s *MangaWorkflowService) buildTaskResult(ctx context.Context, t *task.Task) (*dto.TaskResultResponse, error) {
	// 加载小说
	novelEntity, err := s.novelRepo.FindByID(ctx, novel.NovelID(t.NovelID))
	if err != nil {
		return nil, err
	}

	// 使用 FindByNovelID 查询所有关联的媒体文件
	mediaList, err := s.mediaRepo.FindByNovelID(ctx, t.NovelID)
	if err != nil {
		log.Printf("Failed to find media by novel ID: %v", err)
		return nil, err
	}

	// 构建漫画图片响应
	var mangaImages []dto.TaskSceneResponse
	for i, mediaEntity := range mediaList {
		if mediaEntity.Status == media.MediaStatusCompleted {
			mangaImages = append(mangaImages, dto.TaskSceneResponse{
				ID:          string(mediaEntity.ID),
				SequenceNum: i + 1,
				Description: fmt.Sprintf("漫画面板 %d", i+1),
				ImageURL:    mediaEntity.URL,
			})
		}
	}

	return &dto.TaskResultResponse{
		NovelID:        string(novelEntity.ID),
		Title:          novelEntity.Title,
		CharacterCount: 0,
		SceneCount:     len(mangaImages),
		Characters:     []dto.TaskCharacterResponse{},
		Scenes:         mangaImages,
	}, nil
}

// GetTaskList 获取任务列表
func (s *MangaWorkflowService) GetTaskList(ctx context.Context, userID string, page, pageSize int, status string) ([]*dto.TaskListItemResponse, *dto.PaginationResponse, error) {
	// 查询任务列表
	tasks, total, err := s.taskRepo.FindByUserID(ctx, userID, page, pageSize, status)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch tasks: %w", err)
	}

	// 构建响应列表
	items := make([]*dto.TaskListItemResponse, 0, len(tasks))
	for _, t := range tasks {
		item := &dto.TaskListItemResponse{
			TaskID: t.ID,
			Status: string(t.Status),
			Progress: dto.TaskProgressResponse{
				CurrentStep:      t.ProgressStep,
				CurrentStepIndex: t.ProgressStepIndex,
				TotalSteps:       6,
				Percentage:       t.ProgressPercentage,
			},
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			CompletedAt: t.CompletedAt,
			FailedAt:    t.FailedAt,
		}

		// 加载小说标题
		if novelEntity, err := s.novelRepo.FindByID(ctx, novel.NovelID(t.NovelID)); err == nil {
			item.Title = novelEntity.Title
		}

		// 如果任务完成，添加统计信息
		if t.Status == task.TaskStatusCompleted {
			item.CharacterCount = 0
			item.SceneCount = t.ProgressDetails.ScenesGenerated // 漫画图片数量
		}

		// 如果任务失败，添加错误信息
		if t.Status == task.TaskStatusFailed {
			item.Error = &dto.TaskErrorResponse{
				Code:      t.ErrorCode,
				Message:   t.ErrorMessage,
				RetryAble: t.IsRetryable(),
			}
		}

		items = append(items, item)
	}

	// 计算分页信息
	totalPages := (total + pageSize - 1) / pageSize
	pagination := &dto.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return items, pagination, nil
}

// CancelTask 取消任务
func (s *MangaWorkflowService) CancelTask(ctx context.Context, userID, taskID string) error {
	// 根据用户ID和任务ID查找任务（权限验证）
	taskEntity, err := s.taskRepo.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return fmt.Errorf("task not found or access denied")
	}

	// 检查任务状态
	if taskEntity.Status != task.TaskStatusPending && taskEntity.Status != task.TaskStatusProcessing {
		return fmt.Errorf("任务已完成或已取消，无法取消")
	}

	// 取消任务
	taskEntity.Cancel()
	if err := s.taskRepo.Save(ctx, taskEntity); err != nil {
		return fmt.Errorf("failed to cancel task: %w", err)
	}

	return nil
}
