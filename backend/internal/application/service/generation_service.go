package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/media"
	"github.com/xiajiayi/ai-motion/internal/domain/scene"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/gemini"
	"github.com/xiajiayi/ai-motion/internal/infrastructure/ai/sora"
)

type GenerationService struct {
	mediaRepo    media.MediaRepository
	sceneRepo    scene.SceneRepository
	geminiClient *gemini.Client
	soraClient   *sora.Client
}

func NewGenerationService(
	mediaRepo media.MediaRepository,
	sceneRepo scene.SceneRepository,
	geminiClient *gemini.Client,
	soraClient *sora.Client,
) *GenerationService {
	return &GenerationService{
		mediaRepo:    mediaRepo,
		sceneRepo:    sceneRepo,
		geminiClient: geminiClient,
		soraClient:   soraClient,
	}
}

func (s *GenerationService) GenerateSceneImage(ctx context.Context, req *dto.GenerateImageRequest) (*dto.MediaResponse, error) {
	sceneEntity, err := s.sceneRepo.FindByID(ctx, scene.SceneID(req.SceneID))
	if err != nil {
		return nil, fmt.Errorf("failed to find scene: %w", err)
	}

	mediaEntity := media.NewMedia(string(sceneEntity.ID), media.MediaTypeImage)

	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return nil, fmt.Errorf("failed to save media: %w", err)
	}

	var imageURL string
	if req.ReferenceImage != "" {
		geminiReq := gemini.ImageToImageRequest{
			ReferenceImage: req.ReferenceImage,
			Prompt:         req.Prompt,
			NegativePrompt: req.NegativePrompt,
			Width:          req.Width,
			Height:         req.Height,
		}
		imageURL, err = s.geminiClient.ImageToImage(ctx, geminiReq)
	} else {
		geminiReq := gemini.TextToImageRequest{
			Prompt:         req.Prompt,
			NegativePrompt: req.NegativePrompt,
			Width:          req.Width,
			Height:         req.Height,
			Quality:        req.Quality,
			Style:          req.Style,
		}
		imageURL, err = s.geminiClient.TextToImage(ctx, geminiReq)
	}

	if err != nil {
		mediaEntity.MarkFailed(err.Error())
		s.mediaRepo.Save(ctx, mediaEntity)
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	metadata := media.NewImageMetadata(req.Width, req.Height, "image/jpeg", 0)
	mediaEntity.MarkCompleted(imageURL, metadata)

	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return s.toMediaDTO(mediaEntity), nil
}

func (s *GenerationService) GenerateSceneVideo(ctx context.Context, req *dto.GenerateVideoRequest) (*dto.MediaResponse, error) {
	sceneEntity, err := s.sceneRepo.FindByID(ctx, scene.SceneID(req.SceneID))
	if err != nil {
		return nil, fmt.Errorf("failed to find scene: %w", err)
	}

	mediaEntity := media.NewMedia(string(sceneEntity.ID), media.MediaTypeVideo)

	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return nil, fmt.Errorf("failed to save media: %w", err)
	}

	soraReq := sora.ImageToVideoRequest{
		ImageURL: req.ImageURL,
		Prompt:   req.Prompt,
		Duration: req.Duration,
	}

	videoID, err := s.soraClient.ImageToVideo(ctx, soraReq)
	if err != nil {
		mediaEntity.MarkFailed(err.Error())
		s.mediaRepo.Save(ctx, mediaEntity)
		return nil, fmt.Errorf("failed to generate video: %w", err)
	}

	mediaEntity.MarkGenerating(videoID)
	if err := s.mediaRepo.Save(ctx, mediaEntity); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return s.toMediaDTO(mediaEntity), nil
}

func (s *GenerationService) BatchGenerateScenes(ctx context.Context, req *dto.BatchGenerateRequest) (*dto.BatchGenerateResponse, error) {
	jobID := fmt.Sprintf("job-%d", time.Now().Unix())

	for _, sceneID := range req.SceneIDs {
		if req.GenerateImages {
			sceneEntity, err := s.sceneRepo.FindByID(ctx, scene.SceneID(sceneID))
			if err != nil {
				continue
			}

			mediaEntity := media.NewMedia(string(sceneEntity.ID), media.MediaTypeImage)
			s.mediaRepo.Save(ctx, mediaEntity)
		}

		if req.GenerateVideos {
			sceneEntity, err := s.sceneRepo.FindByID(ctx, scene.SceneID(sceneID))
			if err != nil {
				continue
			}

			mediaEntity := media.NewMedia(string(sceneEntity.ID), media.MediaTypeVideo)
			s.mediaRepo.Save(ctx, mediaEntity)
		}
	}

	return &dto.BatchGenerateResponse{
		JobID:       jobID,
		TotalScenes: len(req.SceneIDs),
		Status:      "pending",
		CreatedAt:   time.Now(),
	}, nil
}

func (s *GenerationService) GetGenerationStatus(ctx context.Context, sceneID string) (*dto.GenerationStatusResponse, error) {
	mediaList, err := s.mediaRepo.FindBySceneID(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("failed to find media: %w", err)
	}

	response := &dto.GenerationStatusResponse{
		TotalTasks: len(mediaList),
		Media:      make([]dto.MediaResponse, 0, len(mediaList)),
	}

	for _, m := range mediaList {
		switch m.Status {
		case media.MediaStatusCompleted:
			response.CompletedTasks++
		case media.MediaStatusFailed:
			response.FailedTasks++
		case media.MediaStatusPending:
			response.PendingTasks++
		}

		response.Media = append(response.Media, *s.toMediaDTO(m))
	}

	return response, nil
}

func (s *GenerationService) toMediaDTO(m *media.Media) *dto.MediaResponse {
	return &dto.MediaResponse{
		ID:      string(m.ID),
		SceneID: m.SceneID,
		Type:    string(m.Type),
		Status:  string(m.Status),
		URL:     m.URL,
		Metadata: dto.MediaMetadata{
			Width:      m.Metadata.Width,
			Height:     m.Metadata.Height,
			Duration:   m.Metadata.Duration,
			Format:     m.Metadata.Format,
			FileSize:   m.Metadata.FileSize,
			Resolution: m.Metadata.Resolution,
		},
		GenerationID: m.GenerationID,
		ErrorMessage: m.ErrorMessage,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		CompletedAt:  m.CompletedAt,
	}
}
