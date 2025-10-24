package model

import "time"

// Scene 场景模型
type Scene struct {
	ID          string    `json:"id" db:"id"`
	ChapterID   string    `json:"chapter_id" db:"chapter_id"`
	NovelID     string    `json:"novel_id" db:"novel_id"`
	SequenceNum int       `json:"sequence_num" db:"sequence_num"` // 场景序号
	Description string    `json:"description" db:"description"`   // 场景描述
	Dialogue    string    `json:"dialogue" db:"dialogue"`         // 对话内容
	Characters  []string  `json:"characters" db:"-"`              // 出现的角色ID列表
	ImageURL    string    `json:"image_url" db:"image_url"`       // 生成的图片URL
	AudioURL    string    `json:"audio_url" db:"audio_url"`       // 生成的音频URL
	Duration    int       `json:"duration" db:"duration"`         // 持续时间(秒)
	Status      string    `json:"status" db:"status"`             // pending, generating, completed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
