package model

import "time"

// Character 角色模型
type Character struct {
	ID              string            `json:"id" db:"id"`
	NovelID         string            `json:"novel_id" db:"novel_id"`
	Name            string            `json:"name" db:"name"`
	Description     string            `json:"description" db:"description"`
	Appearance      string            `json:"appearance" db:"appearance"`           // 外貌描述
	Personality     string            `json:"personality" db:"personality"`         // 性格描述
	ReferenceImages []string          `json:"reference_images" db:"-"`              // 参考图片URL列表
	LoraModelPath   string            `json:"lora_model_path" db:"lora_model_path"` // LoRA 模型路径
	VoiceProfile    string            `json:"voice_profile" db:"voice_profile"`     // 语音配置
	Attributes      map[string]string `json:"attributes" db:"-"`                    // 其他属性
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}
