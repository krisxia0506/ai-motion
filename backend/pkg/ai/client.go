package ai

// ImageGenerateRequest 图像生成请求
type ImageGenerateRequest struct {
	Prompt           string   `json:"prompt"`
	NegativePrompt   string   `json:"negative_prompt,omitempty"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	Steps            int      `json:"steps"`
	CharacterLoRA    string   `json:"character_lora,omitempty"`    // 角色 LoRA 模型路径
	ReferenceImages  []string `json:"reference_images,omitempty"`  // 参考图片
}

// ImageGenerateResponse 图像生成响应
type ImageGenerateResponse struct {
	ImageURL string `json:"image_url"`
	Seed     int64  `json:"seed"`
}

// TextAnalyzeRequest 文本分析请求
type TextAnalyzeRequest struct {
	Text    string `json:"text"`
	Type    string `json:"type"`    // character, scene, dialogue
	Context string `json:"context,omitempty"`
}

// TextAnalyzeResponse 文本分析响应
type TextAnalyzeResponse struct {
	Result map[string]interface{} `json:"result"`
}

// VoiceGenerateRequest 语音生成请求
type VoiceGenerateRequest struct {
	Text         string `json:"text"`
	VoiceProfile string `json:"voice_profile"`
	Language     string `json:"language"`
}

// VoiceGenerateResponse 语音生成响应
type VoiceGenerateResponse struct {
	AudioURL string `json:"audio_url"`
	Duration int    `json:"duration"`
}

// AIClient AI 服务客户端接口
type AIClient interface {
	GenerateImage(req *ImageGenerateRequest) (*ImageGenerateResponse, error)
	AnalyzeText(req *TextAnalyzeRequest) (*TextAnalyzeResponse, error)
	GenerateVoice(req *VoiceGenerateRequest) (*VoiceGenerateResponse, error)
}
