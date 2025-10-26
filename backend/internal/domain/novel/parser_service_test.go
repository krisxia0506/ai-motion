package novel

import (
	"strings"
	"testing"
)

func TestNewParserService(t *testing.T) {
	service := NewParserService()

	if service == nil {
		t.Error("NewParserService() should not return nil")
	}

	if len(service.chapterPatterns) == 0 {
		t.Error("NewParserService() should initialize chapter patterns")
	}
}

func TestParserService_Parse(t *testing.T) {
	tests := []struct {
		name           string
		novel          *Novel
		wantErr        bool
		wantChapters   int
		expectedStatus NovelStatus
	}{
		{
			name: "novel with chinese chapter markers",
			novel: &Novel{
				ID:      "test-1",
				Title:   "Test Novel",
				Content: "第一章 开始\n这是第一章的内容。\n\n第二章 继续\n这是第二章的内容。\n\n第三章 结束\n这是第三章的内容。",
				Status:  NovelStatusPending,
			},
			wantErr:        false,
			wantChapters:   3,
			expectedStatus: NovelStatusParsed,
		},
		{
			name: "novel with english chapter markers",
			novel: &Novel{
				ID:      "test-2",
				Title:   "Test Novel",
				Content: "Chapter 1 Beginning\n" + strings.Repeat("This is chapter 1 content with enough text. ", 10) + "\n\nChapter 2 Middle\n" + strings.Repeat("This is chapter 2 content with enough text. ", 10),
				Status:  NovelStatusPending,
			},
			wantErr:        false,
			wantChapters:   2,
			expectedStatus: NovelStatusParsed,
		},
		{
			name: "novel with numeric chapter markers",
			novel: &Novel{
				ID:      "test-3",
				Title:   "Test Novel",
				Content: "1. First Chapter\n" + strings.Repeat("Content of first chapter with enough text. ", 10) + "\n\n2. Second Chapter\n" + strings.Repeat("Content of second chapter with enough text. ", 10),
				Status:  NovelStatusPending,
			},
			wantErr:        false,
			wantChapters:   2,
			expectedStatus: NovelStatusParsed,
		},
		{
			name: "novel without chapter markers - single chapter",
			novel: &Novel{
				ID:        "test-4",
				Title:     "Test Novel",
				Content:   strings.Repeat("This is a long novel without chapter markers. ", 20),
				Status:    NovelStatusPending,
				WordCount: 140,
			},
			wantErr:        false,
			wantChapters:   1,
			expectedStatus: NovelStatusParsed,
		},
		{
			name: "invalid novel - empty title",
			novel: &Novel{
				ID:      "test-5",
				Title:   "",
				Content: "Some content",
				Status:  NovelStatusPending,
			},
			wantErr:        true,
			wantChapters:   0,
			expectedStatus: NovelStatusPending,
		},
		{
			name: "invalid novel - content too short",
			novel: &Novel{
				ID:      "test-6",
				Title:   "Test",
				Content: "Short",
				Status:  NovelStatusPending,
			},
			wantErr:        true,
			wantChapters:   0,
			expectedStatus: NovelStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewParserService()
			err := service.Parse(tt.novel)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got nil")
					return
				}
				if tt.novel.Status != tt.expectedStatus {
					t.Errorf("Parse() Status = %v, want %v", tt.novel.Status, tt.expectedStatus)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
				return
			}

			if len(tt.novel.Chapters) != tt.wantChapters {
				t.Errorf("Parse() chapters count = %v, want %v", len(tt.novel.Chapters), tt.wantChapters)
			}

			if tt.novel.ChapterCount != tt.wantChapters {
				t.Errorf("Parse() ChapterCount = %v, want %v", tt.novel.ChapterCount, tt.wantChapters)
			}

			if tt.novel.Status != tt.expectedStatus {
				t.Errorf("Parse() Status = %v, want %v", tt.novel.Status, tt.expectedStatus)
			}

			for i, chapter := range tt.novel.Chapters {
				if chapter.ID == "" {
					t.Errorf("Parse() Chapter[%d] ID should not be empty", i)
				}
				if chapter.NovelID != tt.novel.ID {
					t.Errorf("Parse() Chapter[%d] NovelID = %v, want %v", i, chapter.NovelID, tt.novel.ID)
				}
				if chapter.ChapterNumber != i+1 {
					t.Errorf("Parse() Chapter[%d] ChapterNumber = %v, want %v", i, chapter.ChapterNumber, i+1)
				}
				if chapter.Title == "" && tt.wantChapters > 1 {
					t.Errorf("Parse() Chapter[%d] Title should not be empty", i)
				}
				if chapter.Content == "" {
					t.Errorf("Parse() Chapter[%d] Content should not be empty", i)
				}
			}
		})
	}
}

func TestParserService_ExtractChapters(t *testing.T) {
	service := NewParserService()

	tests := []struct {
		name    string
		novel   *Novel
		wantLen int
	}{
		{
			name: "multiple chinese chapters",
			novel: &Novel{
				ID:      "test-1",
				Content: "第一章 开始\n内容1\n\n第二章 中间\n内容2\n\n第三章 结束\n内容3",
			},
			wantLen: 3,
		},
		{
			name: "no chapter markers",
			novel: &Novel{
				ID:      "test-2",
				Content: "Just plain text without any chapter markers.",
			},
			wantLen: 0,
		},
		{
			name: "single chapter marker - not enough",
			novel: &Novel{
				ID:      "test-3",
				Content: "第一章 唯一\n这是唯一的内容",
			},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chapters, err := service.extractChapters(tt.novel)

			if err != nil {
				t.Errorf("extractChapters() unexpected error = %v", err)
				return
			}

			if len(chapters) != tt.wantLen {
				t.Errorf("extractChapters() chapters count = %v, want %v", len(chapters), tt.wantLen)
			}
		})
	}
}
