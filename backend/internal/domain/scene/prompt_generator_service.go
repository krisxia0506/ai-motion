package scene

import (
	"context"
	"fmt"
	"strings"
)

type Character struct {
	ID         string
	Name       string
	Appearance CharacterAppearance
}

type CharacterAppearance struct {
	PhysicalTraits   string
	ClothingStyle    string
	DistinctFeatures string
	Age              string
	Height           string
}

func (a CharacterAppearance) ToPrompt() string {
	var parts []string

	if a.Age != "" {
		parts = append(parts, a.Age)
	}
	if a.Height != "" {
		parts = append(parts, a.Height)
	}
	if a.PhysicalTraits != "" {
		parts = append(parts, a.PhysicalTraits)
	}
	if a.ClothingStyle != "" {
		parts = append(parts, a.ClothingStyle)
	}
	if a.DistinctFeatures != "" {
		parts = append(parts, a.DistinctFeatures)
	}

	return strings.Join(parts, ", ")
}

type PromptGeneratorService struct {
	sceneRepo SceneRepository
}

func NewPromptGeneratorService(sceneRepo SceneRepository) *PromptGeneratorService {
	return &PromptGeneratorService{
		sceneRepo: sceneRepo,
	}
}

type PromptStyle string

const (
	PromptStyleAnime     PromptStyle = "anime"
	PromptStyleRealistic PromptStyle = "realistic"
	PromptStyleCartoon   PromptStyle = "cartoon"
	PromptStylePainting  PromptStyle = "painting"
)

type PromptOptions struct {
	Style       PromptStyle
	Quality     string
	AspectRatio string
	Negative    []string
}

func DefaultPromptOptions() PromptOptions {
	return PromptOptions{
		Style:       PromptStyleAnime,
		Quality:     "high quality, detailed",
		AspectRatio: "16:9",
		Negative:    []string{"blurry", "low quality", "distorted"},
	}
}

func (s *PromptGeneratorService) GenerateImagePrompt(
	ctx context.Context,
	scene *Scene,
	characters []Character,
	options PromptOptions,
) (string, error) {
	var promptParts []string

	promptParts = append(promptParts, string(options.Style)+" style")

	if scene.Description.Setting != "" {
		promptParts = append(promptParts, scene.Description.Setting)
	} else if scene.Description.FullText != "" {
		setting := s.extractVisualElements(scene.Description.FullText)
		if setting != "" {
			promptParts = append(promptParts, setting)
		}
	}

	if scene.Location != "" {
		promptParts = append(promptParts, "location: "+scene.Location)
	}

	if scene.TimeOfDay != "" {
		lighting := s.mapTimeToLighting(scene.TimeOfDay)
		promptParts = append(promptParts, lighting)
	}

	if len(characters) > 0 {
		charDescriptions := s.buildCharacterDescriptions(characters)
		promptParts = append(promptParts, charDescriptions...)
	}

	if scene.Description.Action != "" {
		promptParts = append(promptParts, "action: "+scene.Description.Action)
	}

	if scene.Description.Atmosphere != "" {
		promptParts = append(promptParts, "atmosphere: "+scene.Description.Atmosphere)
	}

	promptParts = append(promptParts, options.Quality)

	prompt := strings.Join(promptParts, ", ")

	if len(options.Negative) > 0 {
		negative := strings.Join(options.Negative, ", ")
		prompt = fmt.Sprintf("%s. Negative: %s", prompt, negative)
	}

	scene.SetImagePrompt(prompt)

	if err := s.sceneRepo.Save(ctx, scene); err != nil {
		return "", fmt.Errorf("failed to save scene with prompt: %w", err)
	}

	return prompt, nil
}

func (s *PromptGeneratorService) GenerateVideoPrompt(
	ctx context.Context,
	scene *Scene,
	characters []Character,
	options PromptOptions,
) (string, error) {
	imagePrompt, err := s.GenerateImagePrompt(ctx, scene, characters, options)
	if err != nil {
		return "", err
	}

	var motionParts []string

	motionParts = append(motionParts, imagePrompt)

	if scene.Description.Action != "" {
		motion := s.convertActionToMotion(scene.Description.Action)
		motionParts = append(motionParts, "motion: "+motion)
	}

	if len(scene.Dialogues) > 0 {
		motionParts = append(motionParts, "with dialogue and lip sync")
	}

	motionParts = append(motionParts, "smooth camera movement, cinematic")

	videoPrompt := strings.Join(motionParts, ", ")

	scene.SetVideoPrompt(videoPrompt)

	if err := s.sceneRepo.Save(ctx, scene); err != nil {
		return "", fmt.Errorf("failed to save scene with video prompt: %w", err)
	}

	return videoPrompt, nil
}

func (s *PromptGeneratorService) GenerateBatchPrompts(
	ctx context.Context,
	scenes []*Scene,
	charactersMap map[string][]Character,
	options PromptOptions,
) error {
	for _, scene := range scenes {
		characters, ok := charactersMap[string(scene.ID)]
		if !ok {
			characters = []Character{}
		}

		_, err := s.GenerateImagePrompt(ctx, scene, characters, options)
		if err != nil {
			return fmt.Errorf("failed to generate prompt for scene %s: %w", scene.ID, err)
		}
	}

	return nil
}

func (s *PromptGeneratorService) extractVisualElements(text string) string {
	visualKeywords := []string{
		"房间", "街道", "森林", "山", "河", "海",
		"明亮", "昏暗", "阳光", "月光",
		"红色", "蓝色", "绿色", "金色", "白色", "黑色",
		"大", "小", "高", "低", "宽", "窄",
	}

	var elements []string
	words := strings.Fields(text)

	for _, word := range words {
		for _, keyword := range visualKeywords {
			if strings.Contains(word, keyword) {
				elements = append(elements, word)
				break
			}
		}

		if len(elements) >= 5 {
			break
		}
	}

	return strings.Join(elements, " ")
}

func (s *PromptGeneratorService) buildCharacterDescriptions(characters []Character) []string {
	var descriptions []string

	for i, char := range characters {
		var charDesc strings.Builder

		if i == 0 {
			charDesc.WriteString("main character: ")
		} else {
			charDesc.WriteString(fmt.Sprintf("character %d: ", i+1))
		}

		charDesc.WriteString(char.Name)

		appearance := char.Appearance.ToPrompt()
		if appearance != "" {
			charDesc.WriteString(" (")
			charDesc.WriteString(appearance)
			charDesc.WriteString(")")
		}

		descriptions = append(descriptions, charDesc.String())
	}

	return descriptions
}

func (s *PromptGeneratorService) mapTimeToLighting(timeOfDay string) string {
	lightingMap := map[string]string{
		"morning":   "soft morning light, warm golden hour",
		"noon":      "bright midday sun, clear lighting",
		"afternoon": "warm afternoon light",
		"evening":   "golden hour, sunset lighting",
		"night":     "moonlight, dark atmosphere, night scene",
		"midnight":  "very dark, minimal lighting, mysterious",
		"daytime":   "natural daylight, bright",
	}

	if lighting, ok := lightingMap[timeOfDay]; ok {
		return lighting
	}

	return "natural lighting"
}

func (s *PromptGeneratorService) convertActionToMotion(action string) string {
	if strings.Contains(action, "走") || strings.Contains(action, "跑") {
		return "walking/running motion"
	}
	if strings.Contains(action, "坐") || strings.Contains(action, "站") {
		return "static pose with subtle movements"
	}
	if strings.Contains(action, "打") || strings.Contains(action, "战") {
		return "dynamic action, combat movements"
	}
	if strings.Contains(action, "说") || strings.Contains(action, "笑") {
		return "talking, facial expressions"
	}

	return "natural character movement"
}

func (s *PromptGeneratorService) OptimizePromptForConsistency(
	basePrompt string,
	referenceImageURL string,
) string {
	var parts []string

	parts = append(parts, "maintain character consistency from reference image")
	parts = append(parts, basePrompt)
	parts = append(parts, fmt.Sprintf("reference: %s", referenceImageURL))

	return strings.Join(parts, ". ")
}
