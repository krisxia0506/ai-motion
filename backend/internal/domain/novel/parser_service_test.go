package novel_test

import (
	"testing"

	"github.com/xiajiayi/ai-motion/internal/domain/novel"
)

func TestParserService_Parse(t *testing.T) {
	parser := novel.NewParserService()

	tests := []struct {
		name         string
		content      string
		wantChapters int
		wantErr      bool
	}{
		{
			name: "novel with chapters",
			content: `第一章 开始
这是第一章的内容。
主角张三出现了。

第二章 发展
这是第二章的内容。
张三遇到了李四。`,
			wantChapters: 2,
			wantErr:      false,
		},
		{
			name: "novel without chapter markers",
			content: `这是一个没有章节标记的小说。
内容很长。
` + repeatString("更多内容。\n", 50),
			wantChapters: 1,
			wantErr:      false,
		},
		{
			name:         "empty content",
			content:      "",
			wantChapters: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &novel.Novel{
				Title:   "Test Novel",
				Author:  "Test Author",
				Content: tt.content,
			}

			err := parser.Parse(n)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}

			if len(n.Chapters) != tt.wantChapters {
				t.Errorf("len(Chapters) = %v, want %v", len(n.Chapters), tt.wantChapters)
			}

			if n.ChapterCount != tt.wantChapters {
				t.Errorf("ChapterCount = %v, want %v", n.ChapterCount, tt.wantChapters)
			}

			if n.Status != novel.NovelStatusParsed {
				t.Errorf("Status = %v, want %v", n.Status, novel.NovelStatusParsed)
			}
		})
	}
}

func TestParserService_ParseChapters(t *testing.T) {
	parser := novel.NewParserService()

	content := `第一章 开始
这是第一章的内容。

第二章 发展
这是第二章的内容。

第三章 高潮
这是第三章的内容。`

	n := &novel.Novel{
		Title:   "Test",
		Content: content,
	}

	err := parser.Parse(n)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(n.Chapters) != 3 {
		t.Fatalf("Expected 3 chapters, got %v", len(n.Chapters))
	}

	for i, chapter := range n.Chapters {
		if chapter.ChapterNumber != i+1 {
			t.Errorf("Chapter %v: ChapterNumber = %v, want %v", i, chapter.ChapterNumber, i+1)
		}

		if chapter.Content == "" {
			t.Errorf("Chapter %v: Content is empty", i)
		}

		if chapter.WordCount == 0 {
			t.Errorf("Chapter %v: WordCount is 0", i)
		}
	}
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
