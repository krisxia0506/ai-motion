package novel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type ParserService struct {
	chapterPatterns []*regexp.Regexp
}

func NewParserService() *ParserService {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^第[一二三四五六七八九十百千零〇0-9]+章[\s:：]*(.*?)$`),
		regexp.MustCompile(`(?m)^Chapter\s+(\d+)[\s:：]*(.*?)$`),
		regexp.MustCompile(`(?m)^[0-9]+[\s.、]+(.*?)$`),
	}

	return &ParserService{
		chapterPatterns: patterns,
	}
}

func (s *ParserService) Parse(novel *Novel) error {
	if err := novel.Validate(); err != nil {
		return err
	}

	novel.UpdateStatus(NovelStatusParsing)

	chapters, err := s.extractChapters(novel)
	if err != nil {
		novel.UpdateStatus(NovelStatusFailed)
		return fmt.Errorf("failed to extract chapters: %w", err)
	}

	if len(chapters) == 0 {
		singleChapter := Chapter{
			ID:            uuid.New().String(),
			NovelID:       novel.ID,
			ChapterNumber: 1,
			Title:         "全文",
			Content:       novel.Content,
			WordCount:     novel.WordCount,
		}
		chapters = []Chapter{singleChapter}
	}

	novel.SetChapters(chapters)
	novel.UpdateStatus(NovelStatusParsed)

	return nil
}

func (s *ParserService) extractChapters(novel *Novel) ([]Chapter, error) {
	content := strings.TrimSpace(novel.Content)

	for _, pattern := range s.chapterPatterns {
		matches := pattern.FindAllStringSubmatchIndex(content, -1)
		if len(matches) > 1 {
			return s.splitByPattern(novel, content, matches), nil
		}
	}

	return nil, nil
}

func (s *ParserService) splitByPattern(novel *Novel, content string, matches [][]int) []Chapter {
	chapters := make([]Chapter, 0, len(matches))

	for i, match := range matches {
		chapterNumber := i + 1
		titleStart := match[0]
		titleEnd := match[1]

		title := strings.TrimSpace(content[titleStart:titleEnd])

		contentStart := titleEnd
		contentEnd := len(content)
		if i < len(matches)-1 {
			contentEnd = matches[i+1][0]
		}

		chapterContent := strings.TrimSpace(content[contentStart:contentEnd])

		chapter := Chapter{
			ID:            uuid.New().String(),
			NovelID:       novel.ID,
			ChapterNumber: chapterNumber,
			Title:         title,
			Content:       chapterContent,
			WordCount:     countWords(chapterContent),
		}

		chapters = append(chapters, chapter)
	}

	return chapters
}
