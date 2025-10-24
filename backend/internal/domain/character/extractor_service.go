package character

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type CharacterExtractorService struct {
	repo CharacterRepository
}

func NewCharacterExtractorService(repo CharacterRepository) *CharacterExtractorService {
	return &CharacterExtractorService{
		repo: repo,
	}
}

type ExtractedCharacter struct {
	Name        string
	Role        CharacterRole
	Description string
	Appearances []string
}

func (s *CharacterExtractorService) ExtractFromNovel(ctx context.Context, novelID, content string) ([]*Character, error) {
	extractedChars := s.extractCharacterNames(content)

	existingChars, err := s.repo.FindByNovelID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing characters: %w", err)
	}

	existingNames := make(map[string]bool)
	for _, char := range existingChars {
		existingNames[char.Name] = true
	}

	var newCharacters []*Character

	for _, extracted := range extractedChars {
		if existingNames[extracted.Name] {
			continue
		}

		char, err := NewCharacter(novelID, extracted.Name, extracted.Role)
		if err != nil {
			continue
		}

		char.SetDescription(extracted.Description)

		if len(extracted.Appearances) > 0 {
			appearance := Appearance{
				PhysicalTraits: strings.Join(extracted.Appearances, "; "),
			}
			char.SetAppearance(appearance)
		}

		if err := s.repo.Save(ctx, char); err != nil {
			return nil, fmt.Errorf("failed to save character %s: %w", char.Name, err)
		}

		newCharacters = append(newCharacters, char)
	}

	return newCharacters, nil
}

func (s *CharacterExtractorService) extractCharacterNames(content string) []ExtractedCharacter {
	var characters []ExtractedCharacter
	characterMap := make(map[string]*ExtractedCharacter)

	patterns := []struct {
		regex *regexp.Regexp
		role  CharacterRole
	}{
		{regexp.MustCompile(`([一-龥]{2,4})(?:说道?|道|答|问|喊|叫|笑|哭|想|心想|暗想)`), CharacterRoleMinor},
		{regexp.MustCompile(`"([一-龥]{2,4}),`), CharacterRoleMinor},
		{regexp.MustCompile(`([一-龥]{2,4})(?:心中|眼中|脸上|手中)`), CharacterRoleMinor},
		{regexp.MustCompile(`主角([一-龥]{2,4})`), CharacterRoleMain},
		{regexp.MustCompile(`([一-龥]{2,4})是(?:一个|一位|个)`), CharacterRoleSupporting},
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.regex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) < 2 {
					continue
				}

				name := strings.TrimSpace(match[1])
				if !isValidCharacterName(name) {
					continue
				}

				if existing, ok := characterMap[name]; ok {
					if pattern.role == CharacterRoleMain {
						existing.Role = CharacterRoleMain
					} else if pattern.role == CharacterRoleSupporting && existing.Role == CharacterRoleMinor {
						existing.Role = CharacterRoleSupporting
					}
				} else {
					characterMap[name] = &ExtractedCharacter{
						Name: name,
						Role: pattern.role,
					}
				}
			}
		}

		s.extractAppearanceDescriptions(line, characterMap)
	}

	for _, char := range characterMap {
		characters = append(characters, *char)
	}

	return s.rankAndFilterCharacters(characters, content)
}

func (s *CharacterExtractorService) extractAppearanceDescriptions(line string, characterMap map[string]*ExtractedCharacter) {
	appearanceKeywords := []string{
		"长发", "短发", "黑发", "金发", "白发",
		"美丽", "英俊", "高大", "矮小", "瘦弱", "强壮",
		"眼睛", "面容", "身材", "穿着", "衣服",
		"年轻", "年老", "中年", "少年", "少女",
	}

	for name, char := range characterMap {
		if strings.Contains(line, name) {
			for _, keyword := range appearanceKeywords {
				if strings.Contains(line, keyword) {
					sentences := splitSentences(line)
					for _, sentence := range sentences {
						if strings.Contains(sentence, name) && strings.Contains(sentence, keyword) {
							char.Appearances = append(char.Appearances, strings.TrimSpace(sentence))
							break
						}
					}
				}
			}
		}
	}
}

func (s *CharacterExtractorService) rankAndFilterCharacters(characters []ExtractedCharacter, content string) []ExtractedCharacter {
	type charFreq struct {
		char  ExtractedCharacter
		count int
	}

	var ranked []charFreq

	for _, char := range characters {
		count := strings.Count(content, char.Name)
		ranked = append(ranked, charFreq{char: char, count: count})
	}

	for i := 0; i < len(ranked); i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].count > ranked[i].count {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	var result []ExtractedCharacter
	for _, r := range ranked {
		if r.count >= 3 {
			result = append(result, r.char)
		}
	}

	if len(result) > 0 && result[0].Role != CharacterRoleMain {
		result[0].Role = CharacterRoleMain
	}

	for i := 1; i < len(result) && i < 5; i++ {
		if result[i].Role == CharacterRoleMinor {
			result[i].Role = CharacterRoleSupporting
		}
	}

	return result
}

func isValidCharacterName(name string) bool {
	if len(name) < 2 || len(name) > 4 {
		return false
	}

	for _, r := range name {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}

	invalidNames := map[string]bool{
		"他们": true, "她们": true, "我们": true,
		"这个": true, "那个": true, "什么": true,
		"怎么": true, "为什么": true, "如何": true,
		"现在": true, "然后": true, "接着": true,
		"突然": true, "忽然": true, "立刻": true,
		"马上": true, "一直": true, "总是": true,
		"已经": true, "正在": true, "刚刚": true,
		"于是": true, "因此": true, "所以": true,
	}

	return !invalidNames[name]
}

func splitSentences(text string) []string {
	separators := []string{"。", "!", "?", "!", "?", "…", "\n"}

	sentences := []string{text}
	for _, sep := range separators {
		var newSentences []string
		for _, s := range sentences {
			parts := strings.Split(s, sep)
			for _, part := range parts {
				if strings.TrimSpace(part) != "" {
					newSentences = append(newSentences, part)
				}
			}
		}
		sentences = newSentences
	}

	return sentences
}

func (s *CharacterExtractorService) MergeCharacters(ctx context.Context, novelID string, sourceID, targetID CharacterID) error {
	source, err := s.repo.FindByID(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("failed to find source character: %w", err)
	}

	target, err := s.repo.FindByID(ctx, targetID)
	if err != nil {
		return fmt.Errorf("failed to find target character: %w", err)
	}

	if source.NovelID != novelID || target.NovelID != novelID {
		return fmt.Errorf("characters must belong to the same novel")
	}

	if source.Appearance.PhysicalTraits != "" && target.Appearance.PhysicalTraits == "" {
		target.Appearance.PhysicalTraits = source.Appearance.PhysicalTraits
	}
	if source.Appearance.ClothingStyle != "" && target.Appearance.ClothingStyle == "" {
		target.Appearance.ClothingStyle = source.Appearance.ClothingStyle
	}

	if source.Personality.Traits != "" && target.Personality.Traits == "" {
		target.Personality.Traits = source.Personality.Traits
	}

	if err := s.repo.Save(ctx, target); err != nil {
		return fmt.Errorf("failed to save merged character: %w", err)
	}

	if err := s.repo.Delete(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete source character: %w", err)
	}

	return nil
}
