package service

import (
	"context"
	"fmt"

	"github.com/xiajiayi/ai-motion/internal/application/dto"
	"github.com/xiajiayi/ai-motion/internal/domain/character"
	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

type CharacterService struct {
	characterRepo    character.CharacterRepository
	novelRepo        novel.NovelRepository
	extractorService *character.CharacterExtractorService
}

func NewCharacterService(
	characterRepo character.CharacterRepository,
	novelRepo novel.NovelRepository,
	extractorService *character.CharacterExtractorService,
) *CharacterService {
	return &CharacterService{
		characterRepo:    characterRepo,
		novelRepo:        novelRepo,
		extractorService: extractorService,
	}
}

func (s *CharacterService) ExtractCharacters(ctx context.Context, novelID string) ([]*dto.CharacterResponse, error) {
	n, err := s.novelRepo.FindByID(ctx, novel.NovelID(novelID))
	if err != nil {
		return nil, fmt.Errorf("failed to find novel: %w", err)
	}

	characters, err := s.extractorService.ExtractFromNovel(ctx, novelID, n.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract characters: %w", err)
	}

	responses := make([]*dto.CharacterResponse, len(characters))
	for i, char := range characters {
		responses[i] = s.toCharacterResponse(char)
	}

	return responses, nil
}

func (s *CharacterService) GetCharacter(ctx context.Context, id string) (*dto.CharacterResponse, error) {
	char, err := s.characterRepo.FindByID(ctx, character.CharacterID(id))
	if err != nil {
		return nil, err
	}

	return s.toCharacterResponse(char), nil
}

func (s *CharacterService) GetCharactersByNovelID(ctx context.Context, novelID string) ([]*dto.CharacterResponse, error) {
	characters, err := s.characterRepo.FindByNovelID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get characters: %w", err)
	}

	responses := make([]*dto.CharacterResponse, len(characters))
	for i, char := range characters {
		responses[i] = s.toCharacterResponse(char)
	}

	return responses, nil
}

func (s *CharacterService) UpdateCharacter(ctx context.Context, id string, req *dto.UpdateCharacterRequest) (*dto.CharacterResponse, error) {
	char, err := s.characterRepo.FindByID(ctx, character.CharacterID(id))
	if err != nil {
		return nil, fmt.Errorf("failed to find character: %w", err)
	}

	if req.Name != "" {
		char.Name = req.Name
	}

	if req.Role != "" {
		char.Role = character.CharacterRole(req.Role)
	}

	if req.Description != "" {
		char.SetDescription(req.Description)
	}

	if req.Appearance != nil {
		appearance := character.Appearance{
			PhysicalTraits:   req.Appearance.PhysicalTraits,
			ClothingStyle:    req.Appearance.ClothingStyle,
			DistinctFeatures: req.Appearance.DistinctFeatures,
			Age:              req.Appearance.Age,
			Height:           req.Appearance.Height,
		}
		if err := char.SetAppearance(appearance); err != nil {
			return nil, fmt.Errorf("failed to set appearance: %w", err)
		}
	}

	if req.Personality != nil {
		personality := character.Personality{
			Traits:     req.Personality.Traits,
			Motivation: req.Personality.Motivation,
			Background: req.Personality.Background,
		}
		char.SetPersonality(personality)
	}

	if req.ReferenceImageURL != "" {
		char.SetReferenceImage(req.ReferenceImageURL)
	}

	if err := s.characterRepo.Save(ctx, char); err != nil {
		return nil, fmt.Errorf("failed to save character: %w", err)
	}

	return s.toCharacterResponse(char), nil
}

func (s *CharacterService) DeleteCharacter(ctx context.Context, id string) error {
	if err := s.characterRepo.Delete(ctx, character.CharacterID(id)); err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	return nil
}

func (s *CharacterService) MergeCharacters(ctx context.Context, novelID string, sourceID, targetID string) error {
	err := s.extractorService.MergeCharacters(
		ctx,
		novelID,
		character.CharacterID(sourceID),
		character.CharacterID(targetID),
	)
	if err != nil {
		return fmt.Errorf("failed to merge characters: %w", err)
	}

	return nil
}

func (s *CharacterService) toCharacterResponse(char *character.Character) *dto.CharacterResponse {
	return &dto.CharacterResponse{
		ID:      string(char.ID),
		NovelID: char.NovelID,
		Name:    char.Name,
		Role:    string(char.Role),
		Appearance: dto.AppearanceResponse{
			PhysicalTraits:   char.Appearance.PhysicalTraits,
			ClothingStyle:    char.Appearance.ClothingStyle,
			DistinctFeatures: char.Appearance.DistinctFeatures,
			Age:              char.Appearance.Age,
			Height:           char.Appearance.Height,
		},
		Personality: dto.PersonalityResponse{
			Traits:     char.Personality.Traits,
			Motivation: char.Personality.Motivation,
			Background: char.Personality.Background,
		},
		Description:       char.Description,
		ReferenceImageURL: char.ReferenceImageURL,
		CreatedAt:         char.CreatedAt,
		UpdatedAt:         char.UpdatedAt,
	}
}
