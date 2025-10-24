package novel

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NovelID string

type NovelStatus string

const (
	NovelStatusPending    NovelStatus = "pending"
	NovelStatusParsing    NovelStatus = "parsing"
	NovelStatusParsed     NovelStatus = "parsed"
	NovelStatusProcessing NovelStatus = "processing"
	NovelStatusCompleted  NovelStatus = "completed"
	NovelStatusFailed     NovelStatus = "failed"
)

var (
	ErrNovelNotFound   = errors.New("novel not found")
	ErrInvalidContent  = errors.New("invalid novel content")
	ErrEmptyTitle      = errors.New("novel title cannot be empty")
	ErrContentTooShort = errors.New("novel content is too short")
	ErrInvalidStatus   = errors.New("invalid novel status")
)

type Novel struct {
	ID           NovelID
	Title        string
	Author       string
	Content      string
	Status       NovelStatus
	WordCount    int
	ChapterCount int
	Chapters     []Chapter
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Chapter struct {
	ID            string
	NovelID       NovelID
	ChapterNumber int
	Title         string
	Content       string
	WordCount     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewNovel(title, author, content string) (*Novel, error) {
	if err := validateNovelInput(title, author, content); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Novel{
		ID:        NovelID(uuid.New().String()),
		Title:     strings.TrimSpace(title),
		Author:    strings.TrimSpace(author),
		Content:   content,
		Status:    NovelStatusPending,
		WordCount: countWords(content),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (n *Novel) Validate() error {
	if n.Title == "" {
		return ErrEmptyTitle
	}
	if len(n.Content) < 100 {
		return ErrContentTooShort
	}
	return nil
}

func (n *Novel) UpdateStatus(status NovelStatus) error {
	validStatuses := map[NovelStatus]bool{
		NovelStatusPending:    true,
		NovelStatusParsing:    true,
		NovelStatusParsed:     true,
		NovelStatusProcessing: true,
		NovelStatusCompleted:  true,
		NovelStatusFailed:     true,
	}

	if !validStatuses[status] {
		return ErrInvalidStatus
	}

	n.Status = status
	n.UpdatedAt = time.Now()
	return nil
}

func (n *Novel) SetChapters(chapters []Chapter) {
	n.Chapters = chapters
	n.ChapterCount = len(chapters)
	n.UpdatedAt = time.Now()
}

func validateNovelInput(title, author, content string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	if len(content) < 100 {
		return ErrContentTooShort
	}
	return nil
}

func countWords(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	words := strings.Fields(text)
	chineseCount := 0
	for _, char := range text {
		if char >= 0x4E00 && char <= 0x9FFF {
			chineseCount++
		}
	}

	if chineseCount > len(words) {
		return chineseCount
	}
	return len(words)
}
