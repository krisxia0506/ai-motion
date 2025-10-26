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

	// 步骤1: 解析小说
	if err := s.stepParseNovel(ctx, taskEntity, novelEntity); err != nil {
		taskEntity.MarkFailed(50002, fmt.Sprintf("解析小说失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤2: 提取角色
	characters, err := s.stepExtractCharacters(ctx, taskEntity, novelEntity)
	if err != nil {
		taskEntity.MarkFailed(50003, fmt.Sprintf("提取角色失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤3: 生成角色参考图
	if err := s.stepGenerateCharacterReferences(ctx, taskEntity, characters); err != nil {
		taskEntity.MarkFailed(40001, fmt.Sprintf("生成角色参考图失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤4: 划分场景
	scenes, err := s.stepDivideScenes(ctx, taskEntity, novelEntity)
	if err != nil {
		taskEntity.MarkFailed(50004, fmt.Sprintf("划分场景失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤5: 生成场景图片
	if err := s.stepGenerateSceneImages(ctx, taskEntity, scenes, characters); err != nil {
		taskEntity.MarkFailed(40001, fmt.Sprintf("生成场景图片失败: %v", err))
		s.taskRepo.Save(ctx, taskEntity)
		return
	}

	// 步骤6: 完成
	taskEntity.MarkCompleted()
	s.taskRepo.Save(ctx, taskEntity)
}

// 步骤1: 解析小说
func (s *MangaWorkflowService) stepParseNovel(ctx context.Context, t *task.Task, n *novel.Novel) error {
	if t.IsCancelled() {
		return fmt.Errorf("task cancelled")
	}

	t.UpdateProgress("解析小说", 1, 17, t.ProgressDetails)
	s.taskRepo.Save(ctx, t)

	// 使用 NovelParserService 解析小说
	if err := s.parserService.Parse(n); err != nil {
		return err
	}

	// 保存章节
	if len(n.Chapters) > 0 {
		if err := s.chapterRepo.SaveBatch(ctx, n.Chapters); err != nil {
			return fmt.Errorf("failed to save chapters: %w", err)
		}
	}

	return nil
}

// 步骤2: 提取角色
func (s *MangaWorkflowService) stepExtractCharacters(ctx context.Context, t *task.Task, n *novel.Novel) ([]*character.Character, error) {
	if t.IsCancelled() {
		return nil, fmt.Errorf("task cancelled")
	}

	t.UpdateProgress("提取角色", 2, 33, t.ProgressDetails)
	s.taskRepo.Save(ctx, t)

	// 使用 CharacterExtractorService 提取角色
	characters, err := s.extractorService.ExtractFromNovel(ctx, string(n.ID), n.Content)
	if err != nil {
		return nil, err
	}

	// 更新进度详情
	details := t.ProgressDetails
	details.CharactersExtracted = len(characters)
	t.UpdateProgress("提取角色", 2, 33, details)
	s.taskRepo.Save(ctx, t)

	return characters, nil
}

// 步骤3: 生成角色参考图
func (s *MangaWorkflowService) stepGenerateCharacterReferences(ctx context.Context, t *task.Task, characters []*character.Character) error {
	if t.IsCancelled() {
		return fmt.Errorf("task cancelled")
	}

	details := t.ProgressDetails
	details.CharactersExtracted = len(characters)

	for i, char := range characters {
		if t.IsCancelled() {
			return fmt.Errorf("task cancelled")
		}

		// 更新当前进度
		details.CharactersGenerated = i
		t.UpdateProgress("生成角色参考图", 3, 50, details)
		s.taskRepo.Save(ctx, t)

		// 生成参考图
		if err := s.generateCharacterReferenceImage(ctx, char); err != nil {
			return fmt.Errorf("failed to generate reference for character %s: %w", char.Name, err)
		}
	}

	// 完成参考图生成
	details.CharactersGenerated = len(characters)
	t.UpdateProgress("生成角色参考图", 3, 50, details)
	s.taskRepo.Save(ctx, t)

	return nil
}

// 步骤4: 划分场景
func (s *MangaWorkflowService) stepDivideScenes(ctx context.Context, t *task.Task, n *novel.Novel) ([]*scene.Scene, error) {
	if t.IsCancelled() {
		return nil, fmt.Errorf("task cancelled")
	}

	details := t.ProgressDetails
	t.UpdateProgress("划分场景", 4, 67, details)
	s.taskRepo.Save(ctx, t)

	// 使用 SceneDividerService 划分场景
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
				return nil, fmt.Errorf("failed to divide chapter %s: %w", chapter.ID, err)
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
			return nil, fmt.Errorf("failed to divide novel: %w", err)
		}

		allScenes = scenes
	}

	// 更新进度详情
	details.ScenesDivided = len(allScenes)
	t.UpdateProgress("划分场景", 4, 67, details)
	s.taskRepo.Save(ctx, t)

	return allScenes, nil
}

// 步骤5: 生成场景图片
func (s *MangaWorkflowService) stepGenerateSceneImages(ctx context.Context, t *task.Task, scenes []*scene.Scene, characters []*character.Character) error {
	if t.IsCancelled() {
		return fmt.Errorf("task cancelled")
	}

	details := t.ProgressDetails

	for i, scn := range scenes {
		if t.IsCancelled() {
			return fmt.Errorf("task cancelled")
		}

		// 更新当前进度
		details.ScenesGenerated = i
		t.UpdateProgress("生成场景图片", 5, 83, details)
		s.taskRepo.Save(ctx, t)

		// 匹配角色到场景
		sceneCharacters := s.matchCharactersToScene(scn, characters)
		scn.SetCharacters(sceneCharacters)

		if err := s.sceneRepo.Save(ctx, scn); err != nil {
			return fmt.Errorf("failed to save scene: %w", err)
		}

		// 生成场景图片
		if err := s.generateSceneImage(ctx, scn, characters); err != nil {
			return fmt.Errorf("failed to generate image for scene %d: %w", scn.SceneNumber, err)
		}
	}

	// 完成场景图片生成
	details.ScenesGenerated = len(scenes)
	t.UpdateProgress("生成场景图片", 5, 83, details)
	s.taskRepo.Save(ctx, t)

	return nil
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
			TotalSteps:       6,
			Percentage:       taskEntity.ProgressPercentage,
			Details: &dto.TaskProgressDetailsResponse{
				CharactersExtracted: taskEntity.ProgressDetails.CharactersExtracted,
				CharactersGenerated: taskEntity.ProgressDetails.CharactersGenerated,
				ScenesDivided:       taskEntity.ProgressDetails.ScenesDivided,
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

	// 加载角色
	characters, err := s.characterRepo.FindByNovelID(ctx, t.NovelID)
	if err != nil {
		return nil, err
	}

	// 加载场景
	scenes, err := s.sceneRepo.FindByNovelID(ctx, t.NovelID)
	if err != nil {
		return nil, err
	}

	// 构建角色列表
	characterResponses := make([]dto.TaskCharacterResponse, 0, len(characters))
	for _, char := range characters {
		characterResponses = append(characterResponses, dto.TaskCharacterResponse{
			ID:                string(char.ID),
			Name:              char.Name,
			ReferenceImageURL: char.ReferenceImageURL,
		})
	}

	// 构建场景列表
	sceneResponses := make([]dto.TaskSceneResponse, 0, len(scenes))
	for _, sc := range scenes {
		sceneResponses = append(sceneResponses, dto.TaskSceneResponse{
			ID:          string(sc.ID),
			SequenceNum: sc.SceneNumber,
			Description: sc.Description.FullText,
			ImageURL:    "", // TODO: 从 Media 中获取图片URL
		})
	}

	return &dto.TaskResultResponse{
		NovelID:        string(novelEntity.ID),
		Title:          novelEntity.Title,
		CharacterCount: len(characters),
		SceneCount:     len(scenes),
		Characters:     characterResponses,
		Scenes:         sceneResponses,
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
			item.CharacterCount = t.ProgressDetails.CharactersGenerated
			item.SceneCount = t.ProgressDetails.ScenesGenerated
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
