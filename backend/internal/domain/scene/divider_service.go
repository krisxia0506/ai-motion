package scene

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type SceneDividerService struct {
	repo SceneRepository
}

func NewSceneDividerService(repo SceneRepository) *SceneDividerService {
	return &SceneDividerService{
		repo: repo,
	}
}

type Chapter struct {
	ID      string
	NovelID string
	Content string
}

func (s *SceneDividerService) DivideChapterIntoScenes(ctx context.Context, chapter Chapter) ([]*Scene, error) {
	boundaries := s.detectSceneBoundaries(chapter.Content)

	var scenes []*Scene
	for i, boundary := range boundaries {
		scene, err := NewScene(chapter.ID, chapter.NovelID, i+1)
		if err != nil {
			return nil, fmt.Errorf("failed to create scene %d: %w", i+1, err)
		}

		description := Description{
			FullText: boundary.Content,
		}
		if err := scene.SetDescription(description); err != nil {
			return nil, fmt.Errorf("failed to set description for scene %d: %w", i+1, err)
		}

		dialogues := s.extractDialogues(boundary.Content)
		scene.SetDialogues(dialogues)

		if boundary.Location != "" {
			scene.SetLocation(boundary.Location)
		}
		if boundary.TimeOfDay != "" {
			scene.SetTimeOfDay(boundary.TimeOfDay)
		}

		scenes = append(scenes, scene)
	}

	if err := s.repo.BatchSave(ctx, scenes); err != nil {
		return nil, fmt.Errorf("failed to save scenes: %w", err)
	}

	return scenes, nil
}

type sceneBoundary struct {
	Content   string
	Location  string
	TimeOfDay string
}

func (s *SceneDividerService) detectSceneBoundaries(content string) []sceneBoundary {
	var boundaries []sceneBoundary

	locationMarkers := []string{
		"来到", "进入", "走进", "到达", "回到",
		"离开", "前往", "去往", "出现在",
	}

	timeMarkers := []string{
		"第二天", "次日", "清晨", "早晨", "中午", "下午", "傍晚", "晚上", "深夜", "午夜",
		"天亮", "天黑", "黎明", "黄昏",
		"过了", "之后", "后来", "接着", "随后",
	}

	paragraphs := strings.Split(content, "\n\n")
	if len(paragraphs) == 1 {
		paragraphs = strings.Split(content, "\n")
	}

	var currentScene strings.Builder
	var currentLocation, currentTimeOfDay string

	for i, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		isNewScene := false
		newLocation := ""
		newTimeOfDay := ""

		for _, marker := range locationMarkers {
			if strings.Contains(para, marker) {
				isNewScene = true
				newLocation = s.extractLocation(para, marker)
				break
			}
		}

		for _, marker := range timeMarkers {
			if strings.Contains(para, marker) {
				isNewScene = true
				newTimeOfDay = marker
				break
			}
		}

		if isNewScene && currentScene.Len() > 100 {
			boundaries = append(boundaries, sceneBoundary{
				Content:   strings.TrimSpace(currentScene.String()),
				Location:  currentLocation,
				TimeOfDay: currentTimeOfDay,
			})
			currentScene.Reset()
			currentLocation = newLocation
			currentTimeOfDay = newTimeOfDay
		} else {
			if newLocation != "" {
				currentLocation = newLocation
			}
			if newTimeOfDay != "" {
				currentTimeOfDay = newTimeOfDay
			}
		}

		if currentScene.Len() > 0 {
			currentScene.WriteString("\n")
		}
		currentScene.WriteString(para)

		if i == len(paragraphs)-1 && currentScene.Len() > 0 {
			boundaries = append(boundaries, sceneBoundary{
				Content:   strings.TrimSpace(currentScene.String()),
				Location:  currentLocation,
				TimeOfDay: currentTimeOfDay,
			})
		}
	}

	if len(boundaries) == 0 && currentScene.Len() > 0 {
		boundaries = append(boundaries, sceneBoundary{
			Content:   strings.TrimSpace(currentScene.String()),
			Location:  currentLocation,
			TimeOfDay: currentTimeOfDay,
		})
	}

	return boundaries
}

func (s *SceneDividerService) extractLocation(text, marker string) string {
	index := strings.Index(text, marker)
	if index == -1 {
		return ""
	}

	after := text[index+len(marker):]
	words := strings.Fields(after)

	if len(words) > 0 {
		location := words[0]
		location = strings.TrimRight(location, "。,!?;:、")
		return location
	}

	return ""
}

func (s *SceneDividerService) extractDialogues(content string) []Dialogue {
	var dialogues []Dialogue

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`([一-龥]{2,4})(?:说道?|道|答|问|喊|叫|笑|哭)[::""]([^""]+)[""]`),
		regexp.MustCompile(`"([^""]+)"[,。]?([一-龥]{2,4})(?:说道?|道|答|问)`),
		regexp.MustCompile(`([一-龥]{2,4})[::]"([^""]+)"`),
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					speaker := match[1]
					dialogueContent := match[2]

					if len(dialogueContent) > 1 {
						emotion := s.detectEmotion(line)
						dialogues = append(dialogues, Dialogue{
							Speaker: speaker,
							Content: dialogueContent,
							Emotion: emotion,
						})
					}
				}
			}
		}
	}

	return dialogues
}

func (s *SceneDividerService) detectEmotion(text string) string {
	emotions := map[string]string{
		"笑":  "happy",
		"哭":  "sad",
		"喊":  "angry",
		"怒":  "angry",
		"惊":  "surprised",
		"叹":  "sigh",
		"激动": "excited",
		"冷静": "calm",
		"温柔": "gentle",
	}

	for keyword, emotion := range emotions {
		if strings.Contains(text, keyword) {
			return emotion
		}
	}

	return "neutral"
}

func (s *SceneDividerService) EnhanceSceneWithMetadata(ctx context.Context, sceneID SceneID, characters []string) error {
	scene, err := s.repo.FindByID(ctx, sceneID)
	if err != nil {
		return fmt.Errorf("failed to find scene: %w", err)
	}

	scene.SetCharacters(characters)

	if scene.Location == "" {
		scene.SetLocation(s.inferLocation(scene.Description.FullText))
	}

	if scene.TimeOfDay == "" {
		scene.SetTimeOfDay(s.inferTimeOfDay(scene.Description.FullText))
	}

	if err := s.repo.Save(ctx, scene); err != nil {
		return fmt.Errorf("failed to save enhanced scene: %w", err)
	}

	return nil
}

func (s *SceneDividerService) inferLocation(text string) string {
	locations := []string{
		"房间", "卧室", "客厅", "书房", "厨房",
		"街道", "街上", "路上", "广场", "公园",
		"学校", "教室", "办公室", "商店",
		"山上", "森林", "河边", "海边",
	}

	for _, loc := range locations {
		if strings.Contains(text, loc) {
			return loc
		}
	}

	return "未知地点"
}

func (s *SceneDividerService) inferTimeOfDay(text string) string {
	timeKeywords := map[string]string{
		"清晨": "morning",
		"早晨": "morning",
		"上午": "morning",
		"中午": "noon",
		"下午": "afternoon",
		"傍晚": "evening",
		"晚上": "night",
		"深夜": "night",
		"午夜": "midnight",
		"阳光": "daytime",
		"月光": "night",
		"星空": "night",
	}

	for keyword, time := range timeKeywords {
		if strings.Contains(text, keyword) {
			return time
		}
	}

	return "daytime"
}
