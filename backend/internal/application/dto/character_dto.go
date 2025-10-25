package dto

import "time"

type CharacterResponse struct {
	ID                string              `json:"id"`
	NovelID           string              `json:"novel_id"`
	Name              string              `json:"name"`
	Role              string              `json:"role"`
	Appearance        AppearanceResponse  `json:"appearance"`
	Personality       PersonalityResponse `json:"personality"`
	Description       string              `json:"description"`
	ReferenceImageURL string              `json:"reference_image_url,omitempty"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type AppearanceResponse struct {
	PhysicalTraits   string `json:"physical_traits,omitempty"`
	ClothingStyle    string `json:"clothing_style,omitempty"`
	DistinctFeatures string `json:"distinct_features,omitempty"`
	Age              string `json:"age,omitempty"`
	Height           string `json:"height,omitempty"`
}

type PersonalityResponse struct {
	Traits     string `json:"traits,omitempty"`
	Motivation string `json:"motivation,omitempty"`
	Background string `json:"background,omitempty"`
}

type UpdateCharacterRequest struct {
	Name              string                    `json:"name"`
	Role              string                    `json:"role"`
	Appearance        *UpdateAppearanceRequest  `json:"appearance,omitempty"`
	Personality       *UpdatePersonalityRequest `json:"personality,omitempty"`
	Description       string                    `json:"description,omitempty"`
	ReferenceImageURL string                    `json:"reference_image_url,omitempty"`
}

type UpdateAppearanceRequest struct {
	PhysicalTraits   string `json:"physical_traits"`
	ClothingStyle    string `json:"clothing_style"`
	DistinctFeatures string `json:"distinct_features"`
	Age              string `json:"age"`
	Height           string `json:"height"`
}

type UpdatePersonalityRequest struct {
	Traits     string `json:"traits"`
	Motivation string `json:"motivation"`
	Background string `json:"background"`
}

type ExtractCharactersRequest struct {
	NovelID string `json:"novel_id" binding:"required"`
}

type GenerateReferenceImageRequest struct {
	CharacterID string `json:"character_id" binding:"required"`
}

type GenerateReferenceImageResponse struct {
	CharacterID string    `json:"character_id"`
	ImageURL    string    `json:"image_url"`
	PromptUsed  string    `json:"prompt_used"`
	GeneratedAt time.Time `json:"generated_at"`
}

type CharacterListResponse struct {
	Characters []*CharacterResponse `json:"characters"`
	Total      int                  `json:"total"`
}
