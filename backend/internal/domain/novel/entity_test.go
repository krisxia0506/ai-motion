package novel

import (
	"strings"
	"testing"
	"time"
)

func TestNewNovel(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		author      string
		content     string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "valid novel",
			title:   "Test Novel",
			author:  "Test Author",
			content: strings.Repeat("This is a test content. ", 10),
			wantErr: false,
		},
		{
			name:        "empty title",
			title:       "",
			author:      "Test Author",
			content:     strings.Repeat("content ", 20),
			wantErr:     true,
			expectedErr: ErrEmptyTitle,
		},
		{
			name:        "whitespace title",
			title:       "   ",
			author:      "Test Author",
			content:     strings.Repeat("content ", 20),
			wantErr:     true,
			expectedErr: ErrEmptyTitle,
		},
		{
			name:        "content too short",
			title:       "Test Novel",
			author:      "Test Author",
			content:     "Short",
			wantErr:     true,
			expectedErr: ErrContentTooShort,
		},
		{
			name:    "valid novel with trimmed whitespace",
			title:   "  Test Novel  ",
			author:  "  Test Author  ",
			content: strings.Repeat("content ", 20),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			novel, err := NewNovel(tt.title, tt.author, tt.content)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewNovel() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("NewNovel() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewNovel() unexpected error = %v", err)
				return
			}

			if novel.ID == "" {
				t.Error("NewNovel() ID should not be empty")
			}

			if novel.Title != strings.TrimSpace(tt.title) {
				t.Errorf("NewNovel() Title = %v, want %v", novel.Title, strings.TrimSpace(tt.title))
			}

			if novel.Author != strings.TrimSpace(tt.author) {
				t.Errorf("NewNovel() Author = %v, want %v", novel.Author, strings.TrimSpace(tt.author))
			}

			if novel.Status != NovelStatusPending {
				t.Errorf("NewNovel() Status = %v, want %v", novel.Status, NovelStatusPending)
			}

			if novel.WordCount == 0 {
				t.Error("NewNovel() WordCount should not be 0")
			}

			if novel.CreatedAt.IsZero() {
				t.Error("NewNovel() CreatedAt should not be zero")
			}

			if novel.UpdatedAt.IsZero() {
				t.Error("NewNovel() UpdatedAt should not be zero")
			}
		})
	}
}

func TestNovel_Validate(t *testing.T) {
	tests := []struct {
		name        string
		novel       *Novel
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid novel",
			novel: &Novel{
				Title:   "Test Novel",
				Content: strings.Repeat("content ", 20),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			novel: &Novel{
				Title:   "",
				Content: strings.Repeat("content ", 20),
			},
			wantErr:     true,
			expectedErr: ErrEmptyTitle,
		},
		{
			name: "content too short",
			novel: &Novel{
				Title:   "Test Novel",
				Content: "Short",
			},
			wantErr:     true,
			expectedErr: ErrContentTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.novel.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("Validate() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error = %v", err)
			}
		})
	}
}

func TestNovel_UpdateStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      NovelStatus
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "valid status - pending",
			status:  NovelStatusPending,
			wantErr: false,
		},
		{
			name:    "valid status - parsing",
			status:  NovelStatusParsing,
			wantErr: false,
		},
		{
			name:    "valid status - parsed",
			status:  NovelStatusParsed,
			wantErr: false,
		},
		{
			name:    "valid status - processing",
			status:  NovelStatusProcessing,
			wantErr: false,
		},
		{
			name:    "valid status - completed",
			status:  NovelStatusCompleted,
			wantErr: false,
		},
		{
			name:    "valid status - failed",
			status:  NovelStatusFailed,
			wantErr: false,
		},
		{
			name:        "invalid status",
			status:      NovelStatus("invalid"),
			wantErr:     true,
			expectedErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			novel := &Novel{
				Title:     "Test Novel",
				Content:   strings.Repeat("content ", 20),
				UpdatedAt: time.Now().Add(-time.Hour),
			}

			oldUpdatedAt := novel.UpdatedAt
			time.Sleep(time.Millisecond)

			err := novel.UpdateStatus(tt.status)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateStatus() expected error but got nil")
					return
				}
				if tt.expectedErr != nil && err != tt.expectedErr {
					t.Errorf("UpdateStatus() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateStatus() unexpected error = %v", err)
				return
			}

			if novel.Status != tt.status {
				t.Errorf("UpdateStatus() Status = %v, want %v", novel.Status, tt.status)
			}

			if !novel.UpdatedAt.After(oldUpdatedAt) {
				t.Error("UpdateStatus() should update UpdatedAt timestamp")
			}
		})
	}
}

func TestNovel_SetChapters(t *testing.T) {
	novel := &Novel{
		Title:     "Test Novel",
		Content:   strings.Repeat("content ", 20),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	chapters := []Chapter{
		{ID: "1", Title: "Chapter 1", Content: "Content 1"},
		{ID: "2", Title: "Chapter 2", Content: "Content 2"},
		{ID: "3", Title: "Chapter 3", Content: "Content 3"},
	}

	oldUpdatedAt := novel.UpdatedAt
	time.Sleep(time.Millisecond)

	novel.SetChapters(chapters)

	if len(novel.Chapters) != len(chapters) {
		t.Errorf("SetChapters() Chapters length = %v, want %v", len(novel.Chapters), len(chapters))
	}

	if novel.ChapterCount != len(chapters) {
		t.Errorf("SetChapters() ChapterCount = %v, want %v", novel.ChapterCount, len(chapters))
	}

	if !novel.UpdatedAt.After(oldUpdatedAt) {
		t.Error("SetChapters() should update UpdatedAt timestamp")
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "empty string",
			text: "",
			want: 0,
		},
		{
			name: "whitespace only",
			text: "   ",
			want: 0,
		},
		{
			name: "english words",
			text: "Hello world this is a test",
			want: 6,
		},
		{
			name: "chinese characters",
			text: "这是一个测试文本内容",
			want: 10,
		},
		{
			name: "mixed english and chinese",
			text: "这是 test 测试",
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countWords(tt.text)
			if got != tt.want {
				t.Errorf("countWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
