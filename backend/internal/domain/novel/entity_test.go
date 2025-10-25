package novel_test

import (
	"testing"

	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

func TestNewNovel(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		author  string
		content string
		wantErr bool
	}{
		{
			name:    "valid novel",
			title:   "Test Novel",
			author:  "Test Author",
			content: "This is a test novel with enough content to pass validation. It has more than 100 characters to ensure it meets the minimum requirement.",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			author:  "Test Author",
			content: "This is a test novel with enough content to pass validation. It has more than 100 characters to ensure it meets the minimum requirement.",
			wantErr: true,
		},
		{
			name:    "whitespace title",
			title:   "   ",
			author:  "Test Author",
			content: "This is a test novel with enough content to pass validation. It has more than 100 characters to ensure it meets the minimum requirement.",
			wantErr: true,
		},
		{
			name:    "content too short",
			title:   "Test Novel",
			author:  "Test Author",
			content: "Short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := novel.NewNovel(tt.title, tt.author, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewNovel() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewNovel() unexpected error: %v", err)
				return
			}

			if n.Title != tt.title {
				t.Errorf("Title = %v, want %v", n.Title, tt.title)
			}

			if n.Author != tt.author {
				t.Errorf("Author = %v, want %v", n.Author, tt.author)
			}

			if n.Status != novel.NovelStatusPending {
				t.Errorf("Status = %v, want %v", n.Status, novel.NovelStatusPending)
			}

			if n.WordCount == 0 {
				t.Errorf("WordCount should be greater than 0")
			}
		})
	}
}

func TestNovel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		novel   *novel.Novel
		wantErr bool
	}{
		{
			name: "valid novel",
			novel: &novel.Novel{
				Title:   "Test Novel",
				Content: "This is a test novel with enough content to pass validation. It has more than 100 characters to ensure it meets the minimum requirement for validation and testing purposes.",
			},
			wantErr: false,
		},
		{
			name: "empty title",
			novel: &novel.Novel{
				Title:   "",
				Content: "This is a test novel with enough content to pass validation. It has more than 100 characters.",
			},
			wantErr: true,
		},
		{
			name: "content too short",
			novel: &novel.Novel{
				Title:   "Test Novel",
				Content: "Short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.novel.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("Validate() expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestNovel_UpdateStatus(t *testing.T) {
	n := &novel.Novel{
		Title:   "Test",
		Content: "This is a test novel with enough content to pass validation. It has more than 100 characters.",
		Status:  novel.NovelStatusPending,
	}

	tests := []struct {
		name    string
		status  novel.NovelStatus
		wantErr bool
	}{
		{
			name:    "valid status - parsing",
			status:  novel.NovelStatusParsing,
			wantErr: false,
		},
		{
			name:    "valid status - parsed",
			status:  novel.NovelStatusParsed,
			wantErr: false,
		},
		{
			name:    "valid status - completed",
			status:  novel.NovelStatusCompleted,
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  novel.NovelStatus("invalid"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := n.UpdateStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStatus() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error: %v", err)
			}

			if n.Status != tt.status {
				t.Errorf("Status = %v, want %v", n.Status, tt.status)
			}
		})
	}
}

func TestNovel_SetChapters(t *testing.T) {
	n := &novel.Novel{
		Title:   "Test",
		Content: "This is a test novel with enough content.",
	}

	chapters := []novel.Chapter{
		{ID: "1", ChapterNumber: 1, Title: "Chapter 1"},
		{ID: "2", ChapterNumber: 2, Title: "Chapter 2"},
	}

	n.SetChapters(chapters)

	if n.ChapterCount != 2 {
		t.Errorf("ChapterCount = %v, want %v", n.ChapterCount, 2)
	}

	if len(n.Chapters) != 2 {
		t.Errorf("len(Chapters) = %v, want %v", len(n.Chapters), 2)
	}
}
