package model

import "time"

// Novel 小说模型
type Novel struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Content   string    `json:"content" db:"content"`
	Status    string    `json:"status" db:"status"` // pending, parsing, completed, failed
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Chapter 章节模型
type Chapter struct {
	ID        string    `json:"id" db:"id"`
	NovelID   string    `json:"novel_id" db:"novel_id"`
	Number    int       `json:"number" db:"number"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
